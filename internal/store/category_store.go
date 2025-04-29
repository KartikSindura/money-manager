package store

import (
	"database/sql"
	"time"
)

type Category struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type PostgresCategoryStore struct {
	db *sql.DB
}

func NewPostgresCategoryStore(db *sql.DB) *PostgresCategoryStore {
	return &PostgresCategoryStore{
		db: db,
	}
}

type CategoryStore interface {
	FindOrCreateCategoryByName(category *Category) (*Category, error)
	GetCategoryIDByName(name *string) (*int64, error)
	// GetCategoryByID(id int64) (*Category, error)
	// UpdateCategory(category *Category) error
	DeleteCategoryByID(id int64) error
}

func (p *PostgresCategoryStore) FindOrCreateCategoryByName(category *Category) (*Category, error) {
	query := `
    INSERT INTO categories (name)
    VALUES ($1)
    ON CONFLICT (name) DO NOTHING
    RETURNING id, created_at
    `
	err := p.db.QueryRow(query, category.Name).Scan(&category.ID, &category.CreatedAt)
	if err == sql.ErrNoRows {
		// category already exists, fetch it
		query := `SELECT id, created_at FROM categories WHERE name = $1`
		err = p.db.QueryRow(query, category.Name).Scan(&category.ID, &category.CreatedAt)
		if err != nil {
			return nil, err
		}
		return category, nil
	} else if err != nil {
		return nil, err
	}
	return category, nil
}

func (p *PostgresCategoryStore) GetCategoryIDByName(name *string) (*int64, error) {
	query := `SELECT id FROM categories WHERE name = $1`
	var id int64
	err := p.db.QueryRow(query, name).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (p *PostgresCategoryStore) DeleteCategoryByID(id int64) error {
	query := `
    DELETE FROM categories
    WHERE id = $1
    `
	_, err := p.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
