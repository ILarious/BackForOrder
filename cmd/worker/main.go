package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/ILarious/BackForOrder/config"
	"github.com/ILarious/BackForOrder/internal/domain/service"
	orderkafka "github.com/ILarious/BackForOrder/internal/infrastructure/kafka"
	orderpostgres "github.com/ILarious/BackForOrder/internal/infrastructure/postgres"
	"github.com/ILarious/BackForOrder/internal/infrastructure/postgres/migration"
	"github.com/ILarious/BackForOrder/pkg/postgres"
	"github.com/ILarious/BackForOrder/pkg/worker_pool"
)

func main() {
	cfg := config.Load()

	db, err := postgres.New(cfg.Postgres)
	if err != nil {
		log.Fatalf("failed to connect postgres: %v", err)
	}
	defer func() {
		if err := postgres.Close(db); err != nil {
			log.Printf("failed to close postgres: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := migration.Up(ctx, db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	orderRepository, err := orderpostgres.NewOrderRepository(db)
	if err != nil {
		log.Fatalf("failed to create order repository: %v", err)
	}

	workers, err := worker_pool.NewWorkerPool(cfg.WorkerPool.Size)
	if err != nil {
		log.Fatalf("failed to create worker pool: %v", err)
	}

	orderProducer := orderkafka.NewOrderProducer(cfg.Kafka.Brokers, cfg.Kafka.OrderRequestTopic)
	orderService := service.NewOrderService(orderRepository)
	outboxDispatcher := service.NewOutboxDispatcher(orderRepository, orderProducer, workers)
	orderConsumer := orderkafka.NewOrderConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.OrderResponseTopic,
		cfg.Kafka.OrderResponseGroupID,
		orderService,
		workers,
	)

	workerCtx, stopWorker := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopWorker()

	go outboxDispatcher.Run(workerCtx)
	go orderConsumer.Run(workerCtx)

	<-workerCtx.Done()
	stopWorker()

	if err := orderConsumer.Close(); err != nil {
		log.Printf("failed to close order consumer: %v", err)
	}
	if err := orderProducer.Close(); err != nil {
		log.Printf("failed to close order producer: %v", err)
	}

	workers.Close()
	workers.Wait()
}
