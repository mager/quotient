package handler

import (
	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	openai "github.com/sashabaranov/go-openai"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Handler struct {
	fx.In

	Database *firestore.Client
	Log      *zap.SugaredLogger
	Router   *mux.Router

	OpenAI *openai.Client
}

// New creates a Handler struct
func New(h Handler) *Handler {
	h.registerRoutes()
	return &h
}

// RegisterRoutes registers all the routes for the route handler
func (h *Handler) registerRoutes() {
	h.Router.HandleFunc("/health", h.health).Methods("GET")

	// Emoquote
	h.Router.HandleFunc("/me", h.me).Methods("POST")
	h.Router.HandleFunc("/q", h.getQuote).Methods("POST")
	h.Router.HandleFunc("/history", h.history).Methods("POST")

	// BotWhisper
	h.Router.HandleFunc("/persona", h.createPersona).Methods("POST")

	// Geotory
	h.Router.HandleFunc("/g/datasets/{userID}/{id}", h.getDataset).Methods("GET")
}
