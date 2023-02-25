package handler

import (
	"encoding/json"
	"net/http"

	"google.golang.org/api/iterator"
)

type HistoryReq struct {
	Email string `json:"email"`
}

type HistoryResp struct {
	Quotes []Quote `json:"quotes"`
}

func (h *Handler) history(w http.ResponseWriter, r *http.Request) {
	var resp HistoryResp

	// Get user
	user := h.getUser(w, r)
	if user.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get quotes
	quotes := h.getHistory(w, r, user.ID)
	resp.Quotes = quotes

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) getHistory(w http.ResponseWriter, r *http.Request, userID string) []Quote {
	var resp []Quote

	// Get history from Firestore
	quotes := h.Database.Collection("quotes")
	iter := quotes.Where("user", "==", h.Database.Doc("users/"+userID)).Documents(r.Context())
	defer iter.Stop()

	// Iterate through history
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			h.Log.Errorw("Failed to get history", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return resp
		}

		var quote Quote
		err = doc.DataTo(&quote)
		if err != nil {
			h.Log.Errorw("Failed to get history", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return resp
		}

		resp = append(resp, quote)
	}

	return resp
}
