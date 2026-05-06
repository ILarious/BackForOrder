package http

import (
	"github.com/ILarious/BackForOrder/internal/app/http/handler"
	"github.com/go-chi/chi/v5"
)

func NewRouter(health *handler.Health, openAPI *handler.OpenAPI, orders *handler.OrderHandler) chi.Router {
	r := chi.NewRouter()

	r.Get("/health", health.Health)
	r.Get("/openapi.yaml", openAPI.Spec)
	r.Get("/docs", openAPI.Docs)
	r.Get("/orders", orders.List)
	r.Post("/orders", orders.Create)

	return r
}
