package parser

import (
	"context"
	"log/slog"
	"runtime"
	"strings"
	"sync"

	"github.com/mrbelka12000/wallet_calc/internal/client/ai"
	domainmodel "github.com/mrbelka12000/wallet_calc/internal/domain_model"
	"github.com/mrbelka12000/wallet_calc/pkg/config"
	pdftotext "github.com/mrbelka12000/wallet_calc/pkg/pdf_to_text"
)

const (
	kOffsetForParse = 20

	offsetForAI = 50
)

type KParser struct {
	aiClient     aiClient
	getter       getter
	log          *slog.Logger
	numOfWorkers int
}

func NewKParser(cfg config.Config, log *slog.Logger) *KParser {
	return &KParser{
		aiClient:     ai.NewClient(log.With("component", "ai"), cfg.AIKey),
		getter:       pdftotext.NewGetter(),
		log:          log,
		numOfWorkers: runtime.NumCPU(),
	}
}

func (p *KParser) ParseStatement(ctx context.Context, fileName string, categories []string) ([]domainmodel.BaseTransaction, error) {
	textFromPDF, err := p.getter.GetTextFromPDF(fileName, kOffsetForParse)
	if err != nil {
		return nil, err
	}

	var (
		dataToParse = strings.Split(textFromPDF, "\n")
		result      = make([]domainmodel.BaseTransaction, len(dataToParse))
		mx          = sync.Mutex{}
		wg          = sync.WaitGroup{}
		workerPool  = make(chan struct{}, p.numOfWorkers)
	)

	for i := 0; i < len(dataToParse); i += offsetForAI {
		wg.Add(1)
		workerPool <- struct{}{}

		go func(l int) {
			defer wg.Done()

			r := l + offsetForAI
			if r > len(dataToParse) {
				r = len(dataToParse)
			}

			var tmpResult []ai.Record
			tmpResult, err = p.aiClient.Parse(ctx, ai.ParseRequest{
				TextToParse: strings.Join(dataToParse[l:r], "\n"),
				Categories:  categories,
			})
			if err != nil {
				p.log.Warn("parse error", "err", err)
				return
			}

			mx.Lock()
			defer mx.Unlock()

			ind := l
			for _, tx := range tmpResult {
				result[ind] = domainmodel.BaseTransaction{
					Date:        tx.Date,
					Description: tx.Details,
					Amount:      tx.Amount,
					Transaction: tx.Transaction,
					Category:    tx.Category,
					Confidence:  tx.Confidence,
				}
				ind++
			}

			<-workerPool
		}(i)
	}

	wg.Wait()
	close(workerPool)

	filteredResult := make([]domainmodel.BaseTransaction, 0, len(result))
	var empty domainmodel.BaseTransaction
	for _, tx := range result {
		if tx == empty {
			continue
		}
		filteredResult = append(filteredResult, tx)
	}

	return filteredResult, err
}
