package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/kr/pretty"
)

type GetQuoteReq struct {
	Emotion string `json:"emotion"`
}

type Quote struct {
	Text string `json:"text"`
	Attr string `json:"attr"`
}

type GetQuoteResp struct {
	Req   GetQuoteReq `json:"req"`
	Quote Quote       `json:"quote"`
}

const (
	prompt = "Can you provide a famous quote about %s with attribution to the author?"
)

func (h *Handler) getQuote(w http.ResponseWriter, r *http.Request) {
	var req GetQuoteReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.Log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	completion, err := h.OpenAI.Completion(context.TODO(), gpt3.CompletionRequest{
		Prompt:           []string{fmt.Sprintf(prompt, req.Emotion)},
		MaxTokens:        gpt3.IntPtr(256),
		TopP:             gpt3.Float32Ptr(1),
		Temperature:      gpt3.Float32Ptr(0.7),
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	})
	if err != nil {
		h.Log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(completion.Choices) == 0 {
		h.Log.Error("No choices returned from OpenAI")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := GetQuoteResp{
		Req:   req,
		Quote: adaptQuote(completion.Choices[0].Text),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func adaptQuote(q string) Quote {
	// Parse the string and remove the leading 3 newlines

	pretty.Print(q)
	return Quote{
		Text: q,
		Attr: "Unknown",
	}
}
