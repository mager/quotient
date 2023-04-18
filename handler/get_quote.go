package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type GetQuoteReq struct {
	Topic  string `json:"topic"`
	Source string `json:"source"`
	ID     string `json:"id"`
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
	promptTemplate = "Provide me 1 quote from %s with attribution about %s?"
)

func (h *Handler) getQuote(w http.ResponseWriter, r *http.Request) {
	var req GetQuoteReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	prompt := fmt.Sprintf(promptTemplate, req.Source, req.Topic)

	completion, err := h.OpenAI.CreateChatCompletion(
		context.TODO(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:        64,
			TopP:             1,
			Temperature:      0.7,
			FrequencyPenalty: 0,
			PresencePenalty:  0,
		},
	)

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
		Quote:  adaptQuote(completion.Choices[0].Message.Content),
	}

	h.Log.Infow("Got quote",
		"quote", resp.Quote.Text,
		"attr", resp.Quote.Attr,
		"topic", req.Topic,
		"user", req.ID,
	)

	// Add the quote to the database
	quotes := h.Database.Collection("quotes")
	_, _, err = quotes.Add(r.Context(), map[string]interface{}{
		"text":  resp.Quote.Text,
		"attr":  resp.Quote.Attr,
		"user":  h.Database.Doc("users/" + req.ID),
		"topic": req.Topic,
	})
	if err != nil {
		h.Log.Errorw("Failed to add quote to database", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
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
