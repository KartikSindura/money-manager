package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/KartikSindura/money/internal/store"
	"github.com/KartikSindura/money/utils"
)

type TransactionHandler struct {
	transactionStore store.TransactionStore
	logger           *log.Logger
}

func NewTransactionHandler(transactionStore store.TransactionStore, logger *log.Logger) *TransactionHandler {
	return &TransactionHandler{
		transactionStore: transactionStore,
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
		Amount *float64   `json:"amount"`
		Note   *string    `json:"note"`
		Date   *time.Time `json:"date"`
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

	if updatedExpenseRequest.Note != nil {
		existingExpense.Note = *updatedExpenseRequest.Note
	}
	// date must not be nil
	if updatedExpenseRequest.Date == nil {
		existingExpense.Date = time.Now()
	}
	if updatedExpenseRequest.Date != nil {
		existingExpense.Date = *updatedExpenseRequest.Date
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

func (h *TransactionHandler) HandleCreateIncome(w http.ResponseWriter, r *http.Request) {
	income := &store.Income{}
	err := json.NewDecoder(r.Body).Decode(&income)
	if err != nil {
		h.logger.Printf("Error: decodingHandleCreateIncome: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "error decoding income"})
		return
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
		Amount *float64   `json:"amount"`
		Note   *string    `json:"note"`
		Source *string    `json:"source"`
		Date   *time.Time `json:"date"`
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

	if updatedIncomeRequest.Note != nil {
		existingIncome.Note = *updatedIncomeRequest.Note
	}

	if updatedIncomeRequest.Source != nil {
		existingIncome.Source = *updatedIncomeRequest.Source
	}
	// date must not be nil
	if updatedIncomeRequest.Date == nil {
		existingIncome.Date = time.Now()
	}
	if updatedIncomeRequest.Date != nil {
		existingIncome.Date = *updatedIncomeRequest.Date
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
