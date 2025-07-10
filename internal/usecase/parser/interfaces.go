package parser

import (
	"context"

	"github.com/mrbelka12000/wallet_calc/internal/client/ai"
)

type (
	aiClient interface {
		Parse(ctx context.Context, req ai.ParseRequest) ([]ai.Record, error)
	}

	getter interface {
		GetTextFromPDF(fileName string, customOffset int) (string, error)
	}
)
