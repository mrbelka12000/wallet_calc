package migrations

import (
	"embed"
	"log"

	"github.com/mrbelka12000/wallet_calc/pkg/gorm/postgres"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

func RunMigrations(database *postgres.Gorm) {
	db, err := database.DB.DB()
	if err != nil {
		log.Fatalf("migrations: failed get DB object: %s", err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("migrations: failed set dialect: %s", err)
	}

	if err := goose.Up(db, "."); err != nil {
		log.Fatalf("migrations: failed run Up command: %s", err)
	}
}
