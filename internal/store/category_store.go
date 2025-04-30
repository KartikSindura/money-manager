package store

import (
	"database/sql"
	"time"
)

type Category struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
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
	GetCategoryIDByName(name *string, user_id int64) (*int64, error)
	GetCategories(user_id int64) ([]Category, error)
}

func (p *PostgresCategoryStore) FindOrCreateCategoryByName(category *Category) (*Category, error) {
	query := `
    INSERT INTO categories (user_id, name)
    VALUES ($1, $2)
    ON CONFLICT (user_id, name) DO NOTHING
    RETURNING id, created_at
    `
	err := p.db.QueryRow(query, category.UserID, category.Name).Scan(&category.ID, &category.CreatedAt)
	if err == sql.ErrNoRows {
		// category already exists, fetch it
		query := `SELECT id, created_at FROM categories WHERE user_id = $1 AND name = $2`
		err = p.db.QueryRow(query, category.UserID, category.Name).Scan(&category.ID, &category.CreatedAt)
		if err != nil {
			return nil, err
		}
		return category, nil
	} else if err != nil {
		return nil, err
	}
	return category, nil
}

func (p *PostgresCategoryStore) GetCategoryIDByName(name *string, user_id int64) (*int64, error) {
	query := `SELECT id 
	FROM categories
	WHERE name = $1 AND user_id = $2
	`
	var id int64
	err := p.db.QueryRow(query, name, user_id).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (p *PostgresCategoryStore) GetCategories(user_id int64) ([]Category, error) {
	query := `SELECT id, user_id, name, created_at
	FROM categories
	WHERE user_id = $1`
	rows, err := p.db.Query(query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.UserID, &category.Name, &category.CreatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
