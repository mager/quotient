package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"google.golang.org/api/iterator"
)

type MeReq struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type User struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	Name  string `json:"name"`
}

type MeResp struct {
	User
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	resp := h.getUser(w, r)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) User {
	var req MeReq
	var user User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return user
	}

	users := h.Database.Collection("users")
	// Search users by email
	iter := users.Where("email", "==", req.Email).Documents(r.Context())
	defer iter.Stop()

	// If user exists, return early, otherwise, create user
	doc, err := iter.Next()
	if err == iterator.Done {
		// Create user
		docRef, _, err := users.Add(r.Context(), map[string]interface{}{
			"email": req.Email,
			"name":  req.Name,
		})
		if err != nil {
			h.Log.Errorw("Failed to create user", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return user
		}
		h.Log.Info("Created user", "id", docRef.ID)
		user.ID = docRef.ID
	} else if err != nil {
		h.Log.Errorw("Failed to fetch user", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return user
	} else {
		user.ID = doc.Ref.ID
		h.Log.Infow("Fetching user", "email", req.Email)

	}

	return user
}

func (h *Handler) getUserByEmail(email string) User {
	var user User
	users := h.Database.Collection("users")
	// Search users by email
	iter := users.Where("email", "==", email).Documents(context.TODO())
	defer iter.Stop()

	// If user exists, return early, otherwise, create user
	doc, err := iter.Next()
	if err == iterator.Done {
		h.Log.Errorw("User not found", "email", email)
	} else if err != nil {
		h.Log.Errorw("Failed to fetch user", "error", err.Error())
	} else {
		doc.DataTo(&user)
		h.Log.Infow("Fetching user", "email", email)
	}

	return user
}
