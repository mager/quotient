package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
)

type GetQuoteReq struct {
	Emotion string `json:"emotion"`
}

type Quote struct {
	Text string `json:"text"`
	Attr string `json:"attr"`
}

type GetQuoteResp struct {
	Req    GetQuoteReq `json:"req"`
	Quote  Quote       `json:"quote"`
	Prompt string      `json:"prompt"`
}

const (
	promptTemplate = "Can you provide one famous quote with attribution about %s?"
)

func (h *Handler) getQuote(w http.ResponseWriter, r *http.Request) {
	var req GetQuoteReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	prompt := fmt.Sprintf(promptTemplate, req.Emotion)

	completion, err := h.OpenAI.Completion(context.TODO(), gpt3.CompletionRequest{
		Prompt:           []string{prompt},
		MaxTokens:        gpt3.IntPtr(64),
		TopP:             gpt3.Float32Ptr(1),
		Temperature:      gpt3.Float32Ptr(0.7),
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		LogProbs:         gpt3.IntPtr(1),
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
		Req:    req,
		Prompt: prompt,
		Quote:  adaptQuote(completion.Choices[0].Text),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func adaptQuote(q string) Quote {
	// Strip out new lines
	q = strings.ReplaceAll(q, "\n", "")

	// Parse out the quote and attribution
	re := regexp.MustCompile(`^"(.+)"\s+-\s+(.+)$`)
	match := re.FindStringSubmatch(q)

	if len(match) != 3 {
		return Quote{
			Text: strings.ReplaceAll(q, "\"", ""),
			Attr: "Unknown",
		}
	}

	return Quote{
		Text: match[1],
		Attr: match[2],
	}
}
