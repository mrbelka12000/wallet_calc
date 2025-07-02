package v1

import (
	"log/slog"

	"github.com/gorilla/mux"
)

type (
	Controller struct {
		log                *slog.Logger
		userUsecase        userUsecase
		transactionUsecase transactionUsecase
		categoryUsecase    categoryUsecase
	}
)

func InitControllers(mx *mux.Router, u userUsecase, t transactionUsecase, cu categoryUsecase, log *slog.Logger) {
	c := &Controller{
		log:                log,
		userUsecase:        u,
		transactionUsecase: t,
		categoryUsecase:    cu,
	}

	setUpRoutes(mx, c)
}
