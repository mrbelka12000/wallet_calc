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
	}
)

func InitControllers(mx *mux.Router, u userUsecase, t transactionUsecase, log *slog.Logger) {
	c := &Controller{
		log:                log,
		userUsecase:        u,
		transactionUsecase: t,
	}

	setUpRoutes(mx, c)
}
