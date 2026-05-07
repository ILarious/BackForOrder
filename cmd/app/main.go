package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ILarious/BackForOrder/config"
	httpapp "github.com/ILarious/BackForOrder/internal/app/http"
	"github.com/ILarious/BackForOrder/internal/app/http/handler"
	"github.com/ILarious/BackForOrder/internal/domain/service"
	orderpostgres "github.com/ILarious/BackForOrder/internal/infrastructure/postgres"
	"github.com/ILarious/BackForOrder/internal/infrastructure/postgres/migration"
	"github.com/ILarious/BackForOrder/pkg/postgres"
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

	orderService := service.NewOrderService(orderRepository)

	appCtx, stopApp := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopApp()

	healthHandler := handler.NewHealth()
	openAPIHandler := handler.NewOpenAPI()
	orderHandler := handler.NewOrderHandler(orderService)
	router := httpapp.NewRouter(healthHandler, openAPIHandler, orderHandler)
	srv := httpapp.NewServer(router)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- srv.Run(cfg.ServerPort)
	}()

	select {
	case <-appCtx.Done():
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start server: %v", err)
		}
	}

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelShutdown()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}

	stopApp()
}
