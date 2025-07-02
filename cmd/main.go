package main

import (
	"log/slog"
	"os"

	"github.com/mrbelka12000/wallet_calc/pkg/config"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		panic(err)
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service_name", cfg.ServiceName)

	log.Info("starting service")
}
