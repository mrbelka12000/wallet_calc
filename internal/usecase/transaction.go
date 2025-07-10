package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/mrbelka12000/wallet_calc/internal/converter"
	domainmodel "github.com/mrbelka12000/wallet_calc/internal/domain_model"
	"github.com/mrbelka12000/wallet_calc/internal/dto"
	"github.com/mrbelka12000/wallet_calc/internal/usecase/parser"
	"github.com/mrbelka12000/wallet_calc/pkg/config"
)

var (
	defaultGUID = uuid.MustParse("a6fb4757-62dd-48e5-acbd-a91e0fe9d3ff")
)

type (
	TransactionUsecase struct {
		parsers map[uuid.UUID]baseParser
		cu      categoryUsecase
	}
)

func NewTransactionUsecase(cfg config.Config, log *slog.Logger, cu categoryUsecase) *TransactionUsecase {
	t := &TransactionUsecase{
		parsers: make(map[uuid.UUID]baseParser),
		cu:      cu,
	}

	t.init(cfg, log)

	return t
}

func (t *TransactionUsecase) init(cfg config.Config, log *slog.Logger) {
	t.parsers[defaultGUID] = parser.NewKParser(cfg, log)
}

func (t *TransactionUsecase) ParseStatement(ctx context.Context, fileName string, parserID uuid.UUID) ([]dto.BaseTransaction, error) {
	parser, ok := t.parsers[parserID]
	if !ok {
		return nil, errors.New("parser not found")
	}

	categories, err := t.cu.List(ctx)
	if err != nil {
		return nil, err
	}

	categoriesStr := make([]string, len(categories))
	for i, category := range categories {
		categoriesStr[i] = category.Name
	}

	baseTransactions, err := parser.ParseStatement(ctx, fileName, categoriesStr)
	if err != nil {
		return nil, err
	}

	return converter.ConvertFromDomainModelBtToDTO(baseTransactions), nil
}

func (t *TransactionUsecase) SaveTransactions(ctx context.Context, bt []domainmodel.BaseTransaction) error {
	return nil
}

func groupTransactions(bts []domainmodel.BaseTransaction) map[string][]domainmodel.BaseTransaction {
	group := make(map[string][]domainmodel.BaseTransaction)

	for _, bt := range bts {
		key := getGroupKey(bt)
		group[key] = append(group[key], bt)
	}

	return group
}

func getGroupKey(bt domainmodel.BaseTransaction) string {
	return fmt.Sprintf("%v:%v:%v:%v", bt.Date, bt.Amount, bt.Category, bt.Description)
}
