package main

import (
	"context"
	"log"
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

	healthHandler := handler.NewHealth()
	openAPIHandler := handler.NewOpenAPI()
	orderHandler := handler.NewOrderHandler(orderService)
	router := httpapp.NewRouter(healthHandler, openAPIHandler, orderHandler)
	srv := httpapp.NewServer(router)

	if err := srv.Run(cfg.ServerPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
