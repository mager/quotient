package handler

import (
	"testing"
)

func TestAdaptQuote(t *testing.T) {
	testCases := []struct {
		input         string
		expectedQuote Quote
	}{
		{
			input: "\"Happiness is when what you think, what you say, and what you do are in harmony.\" - Mahatma Gandhi",
			expectedQuote: Quote{
				Text: "Happiness is when what you think, what you say, and what you do are in harmony.",
				Attr: "Mahatma Gandhi",
			},
		},
		{
			input: "\"Love is like pizza - even when it's bad, it's still pretty good.\" - Unknown",
			expectedQuote: Quote{
				Text: "Love is like pizza - even when it's bad, it's still pretty good.",
				Attr: "Unknown",
			},
		},
		{
			input: "\"A quote without attribution is like a book without an author.\"",
			expectedQuote: Quote{
				Text: "A quote without attribution is like a book without an author.",
				Attr: "Unknown",
			},
		},
	}

	for _, tc := range testCases {
		actualQuote := adaptQuote(tc.input)

		if actualQuote.Text != tc.expectedQuote.Text {
			t.Errorf("For input: %s\nExpected quote: %s\nActual quote: %s\n", tc.input, tc.expectedQuote.Text, actualQuote.Text)
		}

		if actualQuote.Attr != tc.expectedQuote.Attr {
			t.Errorf("For input: %s\nExpected attribution: %s\nActual attribution: %s\n", tc.input, tc.expectedQuote.Attr, actualQuote.Attr)
		}
	}
}
