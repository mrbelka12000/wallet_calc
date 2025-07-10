package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

const (
	apiURL          = "https://api.openai.com"
	pathCompletions = "/v1/chat/completions"

	gptModel = "gpt-4o-mini"

	parsePrompt = `
You are a financial data parser. Given the raw text of a bank statement, extract each transaction into a structured JSON array.

Each object in the JSON array should include the following fields:
- "date": transaction date in format DD.MM.YY
- "amount": numeric amount (positive for income, negative for expenses)
- "transaction": type of transaction (e.g., Purchases, Transfers, Replenishment, Others)
- "details": merchant or explanation
- "category": the best-matching category from the provided list
- "confidence": float between 0.0 and 1.0 representing confidence in the category assignment

### Rules:
- **Only use categories from the provided list. Do not invent new categories.**
- Use semantic clues from "details" and "transaction" to determine the category.
- Keep amounts as floats with two decimal places.
- Use lowercased, English category names (e.g., "fuel / gas", not "Fuel / Gas").
- Output must be a valid JSON array.
- Do not include any extra text—return JSON only.

### Example output (no newlines, no explanation, raw JSON only):
[
{
"date": "02.07.25",
"amount": -5000.00,
"transaction": "Purchases",
"details": "ТОО BEATRICE",
"category": "shopping",
"confidence": 1.0
},
{
"date": "30.06.25",
"amount": 20000.00,
"transaction": "Replenishment",
"details": "From card of other banks",
"category": "topup",
"confidence": 1.0
}
]


### Available categories:
%s


### Bank statement:
%s
`
)

type (
	Client struct {
		hc    *http.Client
		log   *slog.Logger
		token string
	}

	ParseRequest struct {
		TextToParse string
		Categories  []string
	}

	In struct {
		Model    string    `json:"model"`
		Messages []Message `json:"messages"`
	}

	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	Out struct {
		Choices []struct {
			Message struct {
				Role    string      `json:"role"`
				Content string      `json:"content"`
				Refusal interface{} `json:"refusal"`
			} `json:"message"`
		} `json:"choices"`
	}
	Record struct {
		Date        string  `json:"date"`
		Amount      float64 `json:"amount"`
		Transaction string  `json:"transaction"`
		Details     string  `json:"details"`
		Category    string  `json:"category"`
		Confidence  float64 `json:"confidence"`
	}
)

func NewClient(log *slog.Logger, token string) *Client {
	return &Client{
		hc:    &http.Client{},
		log:   log,
		token: token,
	}
}

func (c *Client) Parse(ctx context.Context, req ParseRequest) ([]Record, error) {
	var out Out

	err := c.do(ctx,
		In{
			Model: gptModel,
			Messages: []Message{
				{
					Role:    "user",
					Content: fmt.Sprintf(parsePrompt, strings.Join(req.Categories, ", "), req.TextToParse),
				},
			},
		},
		&out,
	)
	if err != nil {
		return nil, err
	}

	if len(out.Choices) == 0 {
		return nil, fmt.Errorf("no choices found")
	}

	var result []Record
	err = json.Unmarshal([]byte(out.Choices[0].Message.Content), &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling choice content: %w", err)
	}

	return result, nil
}
