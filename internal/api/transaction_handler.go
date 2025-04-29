package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/KartikSindura/money/internal/store"
	"github.com/KartikSindura/money/utils"
)

type TransactionHandler struct {
	transactionStore store.TransactionStore
	categoryStore    store.CategoryStore
	logger           *log.Logger
}

func NewTransactionHandler(transactionStore store.TransactionStore, categoryStore store.CategoryStore, logger *log.Logger) *TransactionHandler {
	return &TransactionHandler{
		transactionStore: transactionStore,
		categoryStore:    categoryStore,
		logger:           logger,
	}
}

func (h *TransactionHandler) HandleCreateExpense(w http.ResponseWriter, r *http.Request) {
	expense := &store.Expense{}
	err := json.NewDecoder(r.Body).Decode(&expense)
	if err != nil {
		h.logger.Printf("Error: decodingHandleCreateExpense: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "error decoding expense"})
		return
	}

	// if category is not provided, set it to "Uncategorized"
	categoryName := "uncategorized"
	if expense.Category != nil {
		categoryName = strings.ToLower(*expense.Category)
	}

	category := &store.Category{Name: categoryName}

	// get category id
	category, err = h.categoryStore.FindOrCreateCategoryByName(category)
	if err != nil {
		h.logger.Printf("Error: HandleCreateExpense: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "error getting category"})
		return
	}

	expense.CategoryID = category.ID

	// if date is not provided, set it to now
	if expense.Date == nil {
		now := time.Now()
		expense.Date = &now
	}

	expense, err = h.transactionStore.CreateExpense(expense)
	if err != nil {
		h.logger.Printf("Error: HandleCreateExpense: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "error creating expense"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"expense": expense})
}

func (h *TransactionHandler) HandleGetExpenseByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("Error: HandleGetExpenseByID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid id parameter"})
		return
	}
	expense, err := h.transactionStore.GetExpenseByID(id)
	if err != nil {
		h.logger.Printf("Error: HandleGetExpenseByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "error getting expense"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"expense": expense})
}

func (h *TransactionHandler) HandleUpdateExpense(w http.ResponseWriter, r *http.Request) {
	expenseID, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("Error: HandleUpdateExpense: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid id parameter"})
		return
	}

	existingExpense, err := h.transactionStore.GetExpenseByID(expenseID)
	if err != nil {
		h.logger.Printf("Error: HandleUpdateExpense: %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "cannot find expense"})
		return
	}

	if existingExpense == nil {
		h.logger.Printf("Error: HandleUpdateExpense: %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "expense not found"})
		return
	}

	var updatedExpenseRequest struct {
		Amount     *float64   `json:"amount"`
		Category   *string    `json:"category"`
		CategoryID int64      `json:"category_id"`
		Note       *string    `json:"note"`
		Date       *time.Time `json:"date"`
	}

	err = json.NewDecoder(r.Body).Decode(&updatedExpenseRequest)
	if err != nil {
		h.logger.Printf("Error: decodingHandleUpdateExpense: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "failed to decode request body"})
		return
	}

	// amount must not be nil
	if updatedExpenseRequest.Amount == nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "amount is required"})
		return
	}

	if updatedExpenseRequest.Amount != nil {
		existingExpense.Amount = *updatedExpenseRequest.Amount
	}

	var categoryName string
	if updatedExpenseRequest.Category != nil {
		categoryName = strings.ToLower(*updatedExpenseRequest.Category)
	} else {
		categoryName = "uncategorized"
	}

	category := &store.Category{
		Name: categoryName,
	}
	category, err = h.categoryStore.FindOrCreateCategoryByName(category)
	if err != nil {
		h.logger.Printf("Error: HandleUpdateExpense: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "error getting category"})
		return
	}
	existingExpense.CategoryID = category.ID

	if updatedExpenseRequest.Note != nil {
		existingExpense.Note = *updatedExpenseRequest.Note
	}
	// if date is not provided, set it to now
	if updatedExpenseRequest.Date == nil {
		now := time.Now()
		existingExpense.Date = &now
	}
	if updatedExpenseRequest.Date != nil {
		existingExpense.Date = updatedExpenseRequest.Date
	}
	existingExpense.UpdatedAt = time.Now()

	err = h.transactionStore.UpdateExpense(existingExpense)
	if err != nil {
		h.logger.Printf("Error: HandleUpdateExpense: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to update expense"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"expense": existingExpense})
}

func (h *TransactionHandler) HandleDeleteExpense(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("Error: decodingHandleDeleteExpense: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid id parameter"})
		return
	}
	err = h.transactionStore.DeleteExpenseByID(id)
	if err != nil {
		h.logger.Printf("Error: HandleDeleteExpense: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "error deleting expense"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"status": "expense deleted"})
}

func (h *TransactionHandler) HandleGetExpenses(w http.ResponseWriter, r *http.Request) {
	limit, offset := utils.GetLimitOffset(r)
	expenses, err := h.transactionStore.GetExpenses(limit, offset)
	if err != nil {
		h.logger.Printf("Error: HandleGetExpenses: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "error getting expenses"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"expenses": expenses})
}

func (h *TransactionHandler) HandleGetTotalExpenses(w http.ResponseWriter, r *http.Request) {
	totalExpenses, err := h.transactionStore.GetTotalExpenses()
	if err != nil {
		h.logger.Printf("Error: HandleGetTotalExpenses: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "error getting total expenses"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"totalExpenses": totalExpenses})
}

func (h *TransactionHandler) HandleCreateIncome(w http.ResponseWriter, r *http.Request) {
	income := &store.Income{}
	err := json.NewDecoder(r.Body).Decode(&income)
	if err != nil {
		h.logger.Printf("Error: decodingHandleCreateIncome: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "error decoding income"})
		return
	}

	categoryName := "uncategorized"
	if income.Category != nil {
		categoryName = strings.ToLower(*income.Category)
	}
	category := &store.Category{
		Name: categoryName,
	}
	category, err = h.categoryStore.FindOrCreateCategoryByName(category)
	if err != nil {
		h.logger.Printf("Error: HandleCreateIncome: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "error getting category"})
		return
	}

	income.CategoryID = category.ID

	// if date is not provided, set it to now
	if income.Date == nil {
		now := time.Now()
		income.Date = &now
	}

	income, err = h.transactionStore.CreateIncome(income)
	if err != nil {
		h.logger.Printf("Error: HandleCreateIncome: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "error creating income"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"income": income})
}

func (h *TransactionHandler) HandleGetIncomeByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("Error: decodingHandleGetIncomeByID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid id parameter"})
		return
	}
	income, err := h.transactionStore.GetIncomeByID(id)
	if err != nil {
		h.logger.Printf("Error: HandleGetIncomeByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "error getting income"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"income": income})
}

func (h *TransactionHandler) HandleUpdateIncome(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("Error: decodingHandleUpdateIncome: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid id parameter"})
		return
	}

	existingIncome, err := h.transactionStore.GetIncomeByID(id)
	if err != nil {
		h.logger.Printf("Error: HandleUpdateIncome: %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "cannot find income"})
		return
	}

	var updatedIncomeRequest struct {
		Amount     *float64   `json:"amount"`
		Category   *string    `json:"category"`
		CategoryID int64      `json:"category_id"`
		Note       *string    `json:"note"`
		Source     *string    `json:"source"`
		Date       *time.Time `json:"date"`
	}

	err = json.NewDecoder(r.Body).Decode(&updatedIncomeRequest)
	if err != nil {
		h.logger.Printf("Error: decodingHandleUpdateIncome: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "failed to decode request body"})
		return
	}

	// amount must not be nil
	if updatedIncomeRequest.Amount == nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "amount is required"})
		return
	}

	if updatedIncomeRequest.Amount != nil {
		existingIncome.Amount = *updatedIncomeRequest.Amount
	}

	var categoryName string
	if updatedIncomeRequest.Category != nil {
		categoryName = strings.ToLower(*updatedIncomeRequest.Category)
	} else {
		categoryName = "uncategorized"
	}

	category := &store.Category{
		Name: categoryName,
	}
	category, err = h.categoryStore.FindOrCreateCategoryByName(category)
	if err != nil {
		h.logger.Printf("Error: HandleUpdateIncome: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "error getting category"})
		return
	}
	existingIncome.CategoryID = category.ID

	if updatedIncomeRequest.Note != nil {
		existingIncome.Note = *updatedIncomeRequest.Note
	}

	if updatedIncomeRequest.Source != nil {
		existingIncome.Source = *updatedIncomeRequest.Source
	}
	// if date is not provided, set it to now
	if updatedIncomeRequest.Date == nil {
		now := time.Now()
		existingIncome.Date = &now
	}
	if updatedIncomeRequest.Date != nil {
		existingIncome.Date = updatedIncomeRequest.Date
	}
	existingIncome.UpdatedAt = time.Now()

	err = h.transactionStore.UpdateIncome(existingIncome)
	if err != nil {
		h.logger.Printf("Error: HandleUpdateIncome: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to update income"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"income": existingIncome})
}

func (h *TransactionHandler) HandleDeleteIncome(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("Error: decodingHandleDeleteIncome: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid id parameter"})
		return
	}
	err = h.transactionStore.DeleteIncomeByID(id)
	if err != nil {
		h.logger.Printf("Error: HandleDeleteIncome: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "error deleting income"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"status": "income deleted"})
}

func (h *TransactionHandler) HandleGetIncomes(w http.ResponseWriter, r *http.Request) {
	limit, offset := utils.GetLimitOffset(r)
	incomes, err := h.transactionStore.GetIncomes(limit, offset)
	if err != nil {
		h.logger.Printf("Error: HandleGetIncomes: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "error getting incomes"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"incomes": incomes})
}

func (h *TransactionHandler) HandleGetTotalIncomes(w http.ResponseWriter, r *http.Request) {
	totalIncomes, err := h.transactionStore.GetTotalIncomes()
	if err != nil {
		h.logger.Printf("Error: HandleGetTotalIncomes: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "error getting total incomes"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"totalIncomes": totalIncomes})
}

func (h *TransactionHandler) HandleGetTransactions(w http.ResponseWriter, r *http.Request) {
	limit, offset := utils.GetLimitOffset(r)
	from, to, month, year, _type, categoryName, err := utils.GetTransactionQueryParams(r)
	if err != nil {
		h.logger.Printf("Error: HandleGetTransactions: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "error parsing transaction query params"})
		return
	}

	var categoryID *int64
	if categoryName != nil {
		var err error
		categoryID, err = h.categoryStore.GetCategoryIDByName(categoryName)
		if err == sql.ErrNoRows {
			h.logger.Printf("Error: HandleGetTransactions: %v", err)
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "category not found"})
			return
		}
		if err != nil {
			h.logger.Printf("Error: HandleGetTransactions: %v", err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "fetching category_id by name failed"})
			return
		}
	}
	transactions, err := h.transactionStore.GetTransactions(limit, offset, from, to, month, year, _type, categoryID)
	if err != nil {
		h.logger.Printf("Error: HandleGetTransactions: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "error getting transactions"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"transactions": transactions})
}
