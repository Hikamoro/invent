package api

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Goods struct {
	ID          int // добавим ID, он нужен
	Name        string
	Unit        string
	Quantity    int    // с большой буквы для экспорта
	CategoryID  int    // переименовал Categoryes в CategoryID
	Description string // с большой буквы
}

// GetGoods возвращает список всех товаров
func GetGoods(DB *sql.DB) ([]Goods, error) {
	rows, err := DB.Query(`
        SELECT id, name, unit, quantity, category_id, description
        FROM goods
        ORDER BY id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goods []Goods
	for rows.Next() {
		var g Goods
		if err := rows.Scan(
			&g.ID,
			&g.Name,
			&g.Unit,
			&g.Quantity,
			&g.CategoryID,
			&g.Description,
		); err != nil {
			return nil, err
		}
		goods = append(goods, g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return goods, nil
}

func GetGoodsByID(id int, DB *sql.DB) (*Goods, error) {
	var g Goods
	err := DB.QueryRow(`
        SELECT id, name, unit, quantity, category_id, description
        FROM goods
        WHERE id = $1
    `, id).Scan(
		&g.ID,
		&g.Name,
		&g.Unit,
		&g.Quantity,
		&g.CategoryID,
		&g.Description,
	)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func GetGoodsByName(name string, DB *sql.DB) (*Goods, error) {
	var g Goods
	err := DB.QueryRow(`
        SELECT id, name, unit, quantity, category_id, description
        FROM goods
        WHERE name = $1
    `, name).Scan(
		&g.ID,
		&g.Name,
		&g.Unit,
		&g.Quantity,
		&g.CategoryID,
		&g.Description,
	)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func CreateGoods(g Goods, DB *sql.DB) (*Goods, error) {
	var newGoods Goods
	err := DB.QueryRow(`
        INSERT INTO goods (name, unit, quantity, category_id, description)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `, g.Name, g.Unit, g.Quantity, g.CategoryID, g.Description).Scan(&newGoods.ID)
	if err != nil {
		return nil, err
	}
	return &newGoods, nil
}

func UpdateGoods(g Goods, DB *sql.DB) (*Goods, error) {
	var updatedGoods Goods
	err := DB.QueryRow(`
        UPDATE goods
        SET name = $1, unit = $2, quantity = $3, category_id = $4, description = $5
        WHERE id = $6
        RETURNING id
    `, g.Name, g.Unit, g.Quantity, g.CategoryID, g.Description, g.ID).Scan(&updatedGoods.ID)
	if err != nil {
		return nil, err
	}
	return &updatedGoods, nil
}

func DeleteGoods(id int, DB *sql.DB) error {
	_, err := DB.Exec("DELETE FROM goods WHERE id = $1", id)
	return err
}
