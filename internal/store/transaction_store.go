package store

import (
	"database/sql"
	"fmt"
	"time"
)

type Expense struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	Amount     float64    `json:"amount"`
	CategoryID int64      `json:"category_id"`
	Category   *string    `json:"category,omitempty"`
	Note       string     `json:"note"`
	Date       *time.Time `json:"date"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type Income struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	Amount     float64    `json:"amount"`
	CategoryID int64      `json:"category_id"`
	Category   *string    `json:"category,omitempty"`
	Source     string     `json:"source"`
	Note       string     `json:"note"`
	Date       *time.Time `json:"date"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type Transaction struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Amount     float64   `json:"amount"`
	CategoryID int64     `json:"category_id"`
	Category   *string   `json:"category,omitempty"`
	Note       string    `json:"note"`
	Source     *string   `json:"source"` // for incomes
	Type       string    `json:"type"`   // income or expense
	Date       time.Time `json:"date"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type PostgresTransactionStore struct {
	db *sql.DB
}

func NewPostgresTransactionStore(db *sql.DB) *PostgresTransactionStore {
	return &PostgresTransactionStore{
		db: db,
	}
}

type TransactionStore interface {
	CreateExpense(expense *Expense) (*Expense, error)
	GetExpenseByID(id int64) (*Expense, error)
	UpdateExpense(expense *Expense) error
	DeleteExpenseByID(id int64) error
	GetExpenses(user_id int64, limit int, offset int) ([]Expense, error)
	GetTotalExpenses(user_id int64) (float64, error)

	CreateIncome(income *Income) (*Income, error)
	GetIncomeByID(id int64) (*Income, error)
	UpdateIncome(income *Income) error
	DeleteIncomeByID(id int64) error
	GetIncomes(user_id int64, limit int, offset int) ([]Income, error)
	GetTotalIncomes(user_id int64) (float64, error)

	GetTransactions(user_id int64, limit int, offset int, from *time.Time, to *time.Time, month *int, year *int, _type *string, category *int64) ([]Transaction, error)
}

func (pg *PostgresTransactionStore) CreateExpense(expense *Expense) (*Expense, error) {
	tx, err := pg.db.Begin()

	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO expenses (user_id, amount, category_id, note, date)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(query, expense.UserID, expense.Amount, expense.CategoryID, expense.Note, expense.Date).Scan(&expense.ID, &expense.CreatedAt, &expense.UpdatedAt)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return expense, nil
}

func (pg *PostgresTransactionStore) GetExpenseByID(id int64) (*Expense, error) {
	expense := &Expense{}

	query := `
	SELECT id, user_id, amount, category_id, note, date, created_at, updated_at
	FROM expenses
	WHERE id = $1
	`
	err := pg.db.QueryRow(query, id).Scan(&expense.ID, &expense.UserID, &expense.Amount, &expense.CategoryID, &expense.Note, &expense.Date, &expense.CreatedAt, &expense.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return expense, nil
}

func (pg *PostgresTransactionStore) UpdateExpense(expense *Expense) error {
	tx, err := pg.db.Begin()

	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `
    UPDATE expenses
    SET amount = $1, category_id = $2, note = $3, date = $4, updated_at = $5
    WHERE id = $6
    `
	result, err := tx.Exec(query, expense.Amount, expense.CategoryID, expense.Note, expense.Date, expense.UpdatedAt, expense.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return tx.Commit()
}

func (pg *PostgresTransactionStore) DeleteExpenseByID(id int64) error {
	query := `
    DELETE FROM expenses
    WHERE id = $1
    `
	result, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (pg *PostgresTransactionStore) CreateIncome(income *Income) (*Income, error) {
	tx, err := pg.db.Begin()

	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
    INSERT INTO incomes (user_id, amount, category_id, source, note, date)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id, created_at, updated_at
	`
	err = pg.db.QueryRow(query, income.UserID, income.Amount, income.CategoryID, income.Source, income.Note, income.Date).Scan(&income.ID, &income.CreatedAt, &income.UpdatedAt)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return income, nil
}

func (pg *PostgresTransactionStore) GetIncomeByID(id int64) (*Income, error) {
	income := &Income{}

	query := `
    SELECT id, user_id, amount, category_id, source, note, date, created_at, updated_at
    FROM incomes
    WHERE id = $1
    `
	err := pg.db.QueryRow(query, id).Scan(&income.ID, &income.UserID, &income.Amount, &income.CategoryID, &income.Source, &income.Note, &income.Date, &income.CreatedAt, &income.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return income, nil
}

func (pg *PostgresTransactionStore) UpdateIncome(income *Income) error {
	tx, err := pg.db.Begin()

	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `
    UPDATE incomes
    SET amount = $1, category_id = $2, source = $3, note = $4, date = $5, updated_at = $6
    WHERE id = $7
    `
	result, err := tx.Exec(query, income.Amount, income.CategoryID, income.Source, income.Note, income.Date, income.UpdatedAt, income.ID)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return tx.Commit()
}

func (pg *PostgresTransactionStore) DeleteIncomeByID(id int64) error {
	query := `
    DELETE FROM incomes
    WHERE id = $1
    `
	result, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return err
	}
	return nil
}

func (pg *PostgresTransactionStore) GetExpenses(user_id int64, limit int, offset int) ([]Expense, error) {

	query := `
	SELECT * FROM expenses
	WHERE user_id = $1
	ORDER BY date DESC LIMIT $2 OFFSET $3
	`
	rows, err := pg.db.Query(query, user_id, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("unable to query expenses: %v", err)
	}
	defer rows.Close()

	expenses := []Expense{}
	for rows.Next() {
		expense := Expense{}
		err := rows.Scan(&expense.ID, &expense.UserID, &expense.Amount, &expense.CategoryID, &expense.Note, &expense.Date, &expense.CreatedAt, &expense.UpdatedAt)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	return expenses, nil
}

func (pg *PostgresTransactionStore) GetIncomes(user_id int64, limit int, offset int) ([]Income, error) {

	query := `
    SELECT * FROM incomes
	WHERE user_id = $1
    ORDER BY date DESC LIMIT $2 OFFSET $3
    `
	rows, err := pg.db.Query(query, user_id, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("unable to query incomes: %v", err)
	}
	defer rows.Close()

	incomes := []Income{}
	for rows.Next() {
		income := Income{}
		err := rows.Scan(&income.ID, &income.UserID, &income.Amount, &income.CategoryID, &income.Source, &income.Note, &income.Date, &income.CreatedAt, &income.UpdatedAt)
		if err != nil {
			return nil, err
		}
		incomes = append(incomes, income)
	}
	return incomes, nil
}

func (pg *PostgresTransactionStore) GetTransactions(user_id int64, limit int, offset int, from *time.Time, to *time.Time, month *int, year *int, _type *string, categoryID *int64) ([]Transaction, error) {
	query := `
	SELECT id, user_id, amount, category_id, note, NULL AS source, 'expense' AS type, date, created_at, updated_at
	FROM expenses
	WHERE user_id = $9
	AND	($3::timestamp IS NULL OR date >= $3)
	AND ($4::timestamp IS NULL OR date <= $4)
	AND ($5::int IS NULL OR EXTRACT(MONTH FROM date) = $5)
	AND ($6::int IS NULL OR EXTRACT(YEAR FROM date) = $6)
	AND ($7::text IS NULL OR $7 = 'expense')
	AND ($8::int IS NULL OR $8 = category_id)

	UNION ALL

	SELECT id, user_id, amount, category_id, note, source, 'income' AS type, date, created_at, updated_at
	FROM incomes
	WHERE user_id = $9
	AND ($3::timestamp IS NULL OR date >= $3)
	AND ($4::timestamp IS NULL OR date <= $4)
	AND ($5::int IS NULL OR EXTRACT(MONTH FROM date) = $5)
	AND ($6::int IS NULL OR EXTRACT(YEAR FROM date) = $6)
	AND ($7::text IS NULL OR $7 = 'income')
    AND ($8::int IS NULL OR $8 = category_id)

	ORDER BY date DESC
	LIMIT $1 OFFSET $2
	`
	rows, err := pg.db.Query(query, limit, offset, from, to, month, year, _type, categoryID, user_id)
	if err != nil {
		return nil, fmt.Errorf("unable to query transactions: %v", err)
	}
	defer rows.Close()

	transactions := []Transaction{}
	for rows.Next() {
		transaction := Transaction{}
		err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.CategoryID, &transaction.Note, &transaction.Source, &transaction.Type, &transaction.Date, &transaction.CreatedAt, &transaction.UpdatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

// FIX: sum on 0 entries
func (pg *PostgresTransactionStore) GetTotalExpenses(user_id int64) (float64, error) {
	query := `
    SELECT SUM(amount) FROM expenses
    WHERE user_id = $1
    `
	var total float64
	err := pg.db.QueryRow(query, user_id).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// FIX: sum on 0 entries
func (pg *PostgresTransactionStore) GetTotalIncomes(user_id int64) (float64, error) {
	query := `
    SELECT SUM(amount) FROM incomes
	WHERE user_id = $1
    `
	var total float64
	err := pg.db.QueryRow(query, user_id).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
