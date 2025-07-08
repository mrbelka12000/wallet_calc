package v1

import (
	"log/slog"

	"github.com/gorilla/mux"

	"github.com/mrbelka12000/wallet_calc/pkg/config"
)

type (
	Controller struct {
		log                *slog.Logger
		userUsecase        userUsecase
		transactionUsecase transactionUsecase
		categoryUsecase    categoryUsecase
		jwtSecretKey       []byte
	}
)

func InitControllers(cfg config.Config, mx *mux.Router, u userUsecase, t transactionUsecase, cu categoryUsecase, log *slog.Logger) {
	c := &Controller{
		log:                log.With("component", "controller"),
		userUsecase:        u,
		transactionUsecase: t,
		categoryUsecase:    cu,
		jwtSecretKey:       []byte(cfg.SecretKey),
	}

	setUpRoutes(mx, c)
}
