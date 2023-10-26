package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Dataset struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Source string `json:"source"`
}

type GetDatasetReq struct{}

type GetDatasetResp struct {
	Dataset
}

func (h *Handler) getDataset(w http.ResponseWriter, r *http.Request) {
	// Get path parameters
	var resp GetDatasetResp
	vars := mux.Vars(r)
	resp.UserID = vars["userID"]
	resp.Slug = vars["id"]

	// Fetch dataset from database

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
