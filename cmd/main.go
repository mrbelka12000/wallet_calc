package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	"github.com/mrbelka12000/wallet_calc/internal/client/ai"
	v1 "github.com/mrbelka12000/wallet_calc/internal/controller/http/v1"
	"github.com/mrbelka12000/wallet_calc/internal/repo"
	"github.com/mrbelka12000/wallet_calc/internal/usecase"
	"github.com/mrbelka12000/wallet_calc/migrations"
	"github.com/mrbelka12000/wallet_calc/pkg/config"
	"github.com/mrbelka12000/wallet_calc/pkg/gorm/postgres"
	"github.com/mrbelka12000/wallet_calc/pkg/server"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		panic(err)
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service_name", cfg.ServiceName)

	db, err := postgres.New(cfg.PGURL)
	if err != nil {
		log.With("error", err).Error("failed to connect to database")
		return
	}

	migrations.RunMigrations(db)

	userRepo := repo.NewUser()
	userUsecase := usecase.NewUserUsecase(userRepo, db)

	categoryRepo := repo.NewCategory()
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo, db)

	transactionUsecase := usecase.NewTransactionUsecase(cfg, log)

	aiClient := ai.NewClient(log.With("component", "ai"), cfg.APIKey)
	_ = aiClient

	mx := mux.NewRouter()
	v1.InitControllers(cfg, mx, userUsecase, transactionUsecase, categoryUsecase, log)

	srv := server.New(mx, cfg.HTTPPort)
	srv.Start()

	log.Info("service started", "port", cfg.HTTPPort)

	gs := make(chan os.Signal, 1)
	signal.Notify(gs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-gs:
		log.Info(fmt.Sprintf("Received signal: %d", sig))
		log.Info("Server stopped properly")
		srv.Stop()
		close(gs)
	case err := <-srv.Ch():
		log.With("error", err).Error("Server stopped")
	}
}
