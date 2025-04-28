package routes

import (
	"github.com/KartikSindura/money/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)

	r.Post("/expenses", app.TransactionHandler.HandleCreateExpense)
	r.Get("/expenses/{id}", app.TransactionHandler.HandleGetExpenseByID)
	r.Put("/expenses/{id}", app.TransactionHandler.HandleUpdateExpense)
	r.Delete("/expenses/{id}", app.TransactionHandler.HandleDeleteExpense)

	r.Post("/incomes", app.TransactionHandler.HandleCreateIncome)
	r.Get("/incomes/{id}", app.TransactionHandler.HandleGetIncomeByID)
	r.Put("/incomes/{id}", app.TransactionHandler.HandleUpdateIncome)
	r.Delete("/incomes/{id}", app.TransactionHandler.HandleDeleteIncome)
	// TODO: list all transactions
	// TODO: filter transactions by date
	// TODO: total income
	// TODO: total expense
	// TODO: summarize transactions for a month/year
	// TODO: recurring expenses
	// TODO: add a category field in expenses
	// TODO: listing expenses by categories
	// TODO: auth
	// TODO: pagination
	// TODO:

	return r
}
