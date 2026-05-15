package api

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetCategories(DB *sql.DB) ([]Category, error) {
	rows, err := DB.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	// ВАЖНО: проверяем ошибки, возникшие во время итерации
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func GetCategoryByID(id int, DB *sql.DB) (*Category, error) {
	var category Category
	err := DB.QueryRow("SELECT id, name FROM categories WHERE id = $1", id).Scan(&category.ID, &category.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Category not found
		}
		return nil, err
	}
	return &category, nil
}

func GetCategoryByName(name string, DB *sql.DB) (*Category, error) {
	var category Category
	err := DB.QueryRow("SELECT id, name FROM categories WHERE name = $1", name).Scan(&category.ID, &category.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Category not found
		}
		return nil, err
	}
	return &category, nil
}

func CreateCategory(name string, db *sql.DB) (*Category, error) {
	var cat Category
	err := db.QueryRow(
		"INSERT INTO categories (name) VALUES ($1) RETURNING id, name",
		name,
	).Scan(&cat.ID, &cat.Name)

	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func UpdateCategory(id int, name string, db *sql.DB) (*Category, error) {
	var cat Category
	err := db.QueryRow(
		"UPDATE categories SET name = $1 WHERE id = $2 RETURNING id, name",
		name, id,
	).Scan(&cat.ID, &cat.Name)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func DeleteCategory(id int, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM categories WHERE id = $1", id)
	return err
}
