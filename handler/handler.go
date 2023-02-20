package handler

import (
	"net/http"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Handler struct {
	fx.In

	Router *mux.Router
	Log    *zap.SugaredLogger
	OpenAI gpt3.Client
}

// New creates a Handler struct
func New(h Handler) *Handler {
	h.registerRoutes()
	return &h
}

// RegisterRoutes registers all the routes for the route handler
func (h *Handler) registerRoutes() {
	h.Router.HandleFunc("/health", h.health).Methods("GET")

	h.Router.HandleFunc("/q", h.getQuote).Methods("POST")
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
