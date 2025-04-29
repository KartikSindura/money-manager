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
	r.Get("/expenses", app.TransactionHandler.HandleGetExpenses)
	r.Get("/total-expenses", app.TransactionHandler.HandleGetTotalExpenses)

	r.Post("/incomes", app.TransactionHandler.HandleCreateIncome)
	r.Get("/incomes/{id}", app.TransactionHandler.HandleGetIncomeByID)
	r.Put("/incomes/{id}", app.TransactionHandler.HandleUpdateIncome)
	r.Delete("/incomes/{id}", app.TransactionHandler.HandleDeleteIncome)
	r.Get("/incomes", app.TransactionHandler.HandleGetIncomes)
	r.Get("/total-incomes", app.TransactionHandler.HandleGetTotalIncomes)

	r.Get("/transactions", app.TransactionHandler.HandleGetTransactions)
	// TODO: recurring expenses
	// TODO: listing expenses by categories
	// TODO: auth

	return r
}
