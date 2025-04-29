package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/KartikSindura/money/internal/api"
	"github.com/KartikSindura/money/internal/store"
	"github.com/KartikSindura/money/migrations"
)

type Application struct {
	Logger             *log.Logger
	DB                 *sql.DB
	TransactionHandler *api.TransactionHandler
}

func NewApplication() (*Application, error) {
	pgDb, err := store.Open()
	if err != nil {
		return nil, err
	}
	err = store.MigrateFS(pgDb, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// stores
	postgresTransactionStore := store.NewPostgresTransactionStore(pgDb)
	postgresCategoryStore := store.NewPostgresCategoryStore(pgDb)

	// handlers
	transactionHandler := api.NewTransactionHandler(postgresTransactionStore, postgresCategoryStore, logger)

	app := &Application{
		Logger:             logger,
		DB:                 pgDb,
		TransactionHandler: transactionHandler,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Health gud")
}
