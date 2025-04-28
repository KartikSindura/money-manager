package store

import (
	"database/sql"
	"time"
)

type Expense struct {
	ID        int64     `json:"id"`
	Amount    float64   `json:"amount"`
	Note      string    `json:"note"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Income struct {
	ID        int64     `json:"id"`
	Amount    float64   `json:"amount"`
	Source    string    `json:"source"`
	Note      string    `json:"note"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	CreateIncome(income *Income) (*Income, error)
	GetIncomeByID(id int64) (*Income, error)
	UpdateIncome(income *Income) error
	DeleteIncomeByID(id int64) error
}

func (pg *PostgresTransactionStore) CreateExpense(expense *Expense) (*Expense, error) {
	tx, err := pg.db.Begin()

	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO expenses (amount, note, date)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(query, expense.Amount, expense.Note, expense.Date).Scan(&expense.ID, &expense.CreatedAt, &expense.UpdatedAt)
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
	SELECT id, amount, note, date, created_at, updated_at
	FROM expenses
	WHERE id = $1
	`
	err := pg.db.QueryRow(query, id).Scan(&expense.ID, &expense.Amount, &expense.Note, &expense.Date, &expense.CreatedAt, &expense.UpdatedAt)
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
    SET amount = $1, note = $2, date = $3, updated_at = $4
    WHERE id = $5
    `
	result, err := pg.db.Exec(query, expense.Amount, expense.Note, expense.Date, time.Now(), expense.ID)
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
    INSERT INTO incomes (amount, source, note, date)
    VALUES ($1, $2, $3, $4)
    RETURNING id, created_at, updated_at
	`
	err = pg.db.QueryRow(query, income.Amount, income.Source, income.Note, income.Date).Scan(&income.ID, &income.CreatedAt, &income.UpdatedAt)
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
    SELECT id, amount, source, note, date, created_at, updated_at
    FROM incomes
    WHERE id = $1
    `
	err := pg.db.QueryRow(query, id).Scan(&income.ID, &income.Amount, &income.Source, &income.Note, &income.Date, &income.CreatedAt, &income.UpdatedAt)
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
    SET amount = $1, source = $2, note = $3, date = $4, updated_at = $5
    WHERE id = $6
    `
	result, err := pg.db.Exec(query, income.Amount, income.Source, income.Note, income.Date, time.Now(), income.ID)

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

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return err
	}
	return nil
}
