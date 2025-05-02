package routes

import (
	"github.com/KartikSindura/money/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(app.Middleware.Authenticate)

		r.Post("/expenses", app.Middleware.RequireUser(app.TransactionHandler.HandleCreateExpense))
		r.Get("/expenses/{id}", app.Middleware.RequireUser(app.TransactionHandler.HandleGetExpenseByID))
		r.Put("/expenses/{id}", app.Middleware.RequireUser(app.TransactionHandler.HandleUpdateExpense))
		r.Delete("/expenses/{id}", app.Middleware.RequireUser(app.TransactionHandler.HandleDeleteExpense))
		r.Get("/expenses", app.Middleware.RequireUser(app.TransactionHandler.HandleGetExpenses))
		r.Get("/total-expenses", app.Middleware.RequireUser(app.TransactionHandler.HandleGetTotalExpenses))
		r.Post("/incomes", app.Middleware.RequireUser(app.TransactionHandler.HandleCreateIncome))
		r.Get("/incomes/{id}", app.Middleware.RequireUser(app.TransactionHandler.HandleGetIncomeByID))
		r.Put("/incomes/{id}", app.Middleware.RequireUser(app.TransactionHandler.HandleUpdateIncome))
		r.Delete("/incomes/{id}", app.Middleware.RequireUser(app.TransactionHandler.HandleDeleteIncome))
		r.Get("/incomes", app.Middleware.RequireUser(app.TransactionHandler.HandleGetIncomes))
		r.Get("/total-incomes", app.Middleware.RequireUser(app.TransactionHandler.HandleGetTotalIncomes))
		r.Get("/transactions", app.Middleware.RequireUser(app.TransactionHandler.HandleGetTransactions))
		r.Get("/categories", app.Middleware.RequireUser(app.TransactionHandler.HandleGetCategories))
	})

	r.Get("/health", app.HealthCheck)

	r.Post("/register", app.UserHandler.HandleRegisterUser)
	r.Post("/login", app.UserHandler.HandleLoginUser)
	// TODO: total expenses filtered by month
	// TODO: total incomes filtered by month
	// TODO: recurring transactions
	// TODO: budgeting

	return r
}
