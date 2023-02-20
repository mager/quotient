package handler

import (
	"encoding/json"
	"net/http"

	"google.golang.org/api/iterator"
)

type MeReq struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type MeResp struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	Name  string `json:"name"`
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	var req MeReq
	var resp MeResp
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.Log.Infow("Fetching user", "email", req.Email)

	// Call Firestore
	// If user doesn't exist, create user
	// If user exists, return user

	users := h.Database.Collection("users")
	// Search users by email
	iter := users.Where("email", "==", req.Email).Documents(r.Context())
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		// If user doesn't exist, create user
		if err == iterator.Done {
			// Create user
			docRef, _, err := users.Add(r.Context(), map[string]interface{}{
				"email": req.Email,
				"name":  req.Name,
			})
			if err != nil {
				h.Log.Errorw("Failed to create user", "error", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp.ID = docRef.ID
			break
		}

		if err != nil {
			h.Log.Errorw("Failed to fetch user", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// If user exists, return user
		resp.ID = doc.Ref.ID
	}

	resp.Email = req.Email
	resp.Name = req.Name

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
