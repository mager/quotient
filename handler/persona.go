package handler

import (
	"encoding/json"
	"net/http"
)

type Persona struct {
	User        string `json:"user"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *Handler) createPersona(w http.ResponseWriter, r *http.Request) {
	var persona Persona
	if err := json.NewDecoder(r.Body).Decode(&persona); err != nil {
		h.Log.Errorw("failed to decode request body", "error", err)
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	personas := h.Database.Collection("personas")
	docRef, _, err := personas.Add(r.Context(), map[string]interface{}{
		"user":        h.Database.Collection("users").Doc(persona.User),
		"name":        persona.Name,
		"description": persona.Description,
	})
	if err != nil {
		h.Log.Errorw("failed to create persona", "error", err)
		http.Error(w, "failed to create persona", http.StatusInternalServerError)
		return
	}
	h.Log.Infow("created persona", "id", docRef.ID)

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) getPersonas(w http.ResponseWriter, r *http.Request) {
	var persona Persona
	if err := json.NewDecoder(r.Body).Decode(&persona); err != nil {
		h.Log.Errorw("failed to decode request body", "error", err)
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	personas := h.Database.Collection("personas")
	// Fetch personas owned by user
	iter := personas.Where("user", "==", h.Database.Collection("users").Doc(persona.User)).Documents(r.Context())
	defer iter.Stop()

	var personasResp []Persona
	for {
		doc, err := iter.Next()
		if err != nil {
			h.Log.Errorw("failed to iterate over personas", "error", err)
			http.Error(w, "failed to iterate over personas", http.StatusInternalServerError)
			return
		}
		var personaResp Persona
		if err := doc.DataTo(&personaResp); err != nil {
			h.Log.Errorw("failed to decode persona", "error", err)
			http.Error(w, "failed to decode persona", http.StatusInternalServerError)
			return
		}
		personasResp = append(personasResp, personaResp)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(personasResp)
}
