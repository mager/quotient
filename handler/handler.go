package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Handler struct {
	fx.In

	Router *mux.Router
	Log    *zap.SugaredLogger
}

// New creates a Handler struct
func New(h Handler) *Handler {
	h.registerRoutes()
	return &h
}

// RegisterRoutes registers all the routes for the route handler
func (h *Handler) registerRoutes() {
	h.Router.HandleFunc("/health", h.health).Methods("GET")
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
