package handler

import (
	"encoding/json"
	"net/http"
)

type MeReq struct {
	Email string `json:"email"`
}

type MeResp struct {
	Email string `json:"email"`
	ID    string `json:"id"`
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	var req MeReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.Log.Infow("Fetching user", "email", req.Email)

	// Call Firestore
	// If user doesn't exist, create user
	// If user exists, return user

	resp := MeResp{
		Email: req.Email,
		ID:    "UNF5qof1W02cOKr1JpYx",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}
