package parser

import (
	"context"
	"log/slog"
	"strings"
	"sync"

	"github.com/mrbelka12000/wallet_calc/internal/client/ai"
	domainmodel "github.com/mrbelka12000/wallet_calc/internal/domain_model"
	"github.com/mrbelka12000/wallet_calc/pkg/config"
	pdftotext "github.com/mrbelka12000/wallet_calc/pkg/pdf_to_text"
)

const (
	kOffsetForParse = 20

	offsetForAI = 1000
)

type KParser struct {
	aiClient aiClient
	getter   getter
	log      *slog.Logger
}

func NewKParser(cfg config.Config, log *slog.Logger) *KParser {
	return &KParser{
		aiClient: ai.NewClient(log.With("component", "ai"), cfg.APIKey),
		getter:   pdftotext.NewGetter(),
		log:      log,
	}
}

func (p *KParser) ParseStatement(ctx context.Context, fileName string) ([]domainmodel.BaseTransaction, error) {
	textFromPDF, err := p.getter.GetTextFromPDF(fileName, kOffsetForParse)
	if err != nil {
		return nil, err
	}

	var (
		dataToParse = strings.Split(textFromPDF, "\n")
		result      = make([]domainmodel.BaseTransaction, 0, len(dataToParse))
		mx          = sync.Mutex{}
		wg          = sync.WaitGroup{}
	)

	for i := 0; i < len(dataToParse); i += offsetForAI {
		wg.Add(1)

		go func(l int) {
			defer wg.Done()

			r := l + offsetForAI
			if r > len(dataToParse) {
				r = len(dataToParse)
			}

			var tmpResult []ai.Record
			tmpResult, err = p.aiClient.Parse(ctx, strings.Join(dataToParse[l:r], "\n"))
			if err != nil {
				p.log.Warn("parse error", "err", err)
				return
			}

			mx.Lock()
			defer mx.Unlock()

			for _, tx := range tmpResult {
				result = append(result, domainmodel.BaseTransaction{
					Date:        tx.Date,
					Description: tx.Details,
					Amount:      tx.Amount,
					Transaction: tx.Transaction,
					Category:    tx.Category,
					Confidence:  tx.Confidence,
				})
			}
		}(i)
	}

	wg.Wait()

	return result, err
}
