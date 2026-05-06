package main

import (
	"log"

	"github.com/ILarious/BackForOrder/config"
	httpapp "github.com/ILarious/BackForOrder/internal/app/http"
	"github.com/ILarious/BackForOrder/internal/app/http/handler"
)

func main() {
	cfg := config.Load()

	healthHandler := handler.NewHealth()
	openAPIHandler := handler.NewOpenAPI()
	router := httpapp.NewRouter(healthHandler, openAPIHandler)
	srv := httpapp.NewServer(router)

	if err := srv.Run(cfg.ServerPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
