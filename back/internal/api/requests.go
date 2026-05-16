package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"invent/back/internal/db"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// IssuesCategoriesRequest - обработчик для /api/categories (GET, POST)
func IssuesCategoriesRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		categories, err := GetCategories(db.DB)
		if err != nil {
			http.Error(w, "unable to get categories: "+err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(categories)

	case http.MethodPost:
		var newCategory Category
		if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Валидация: имя не должно быть пустым
		if newCategory.Name == "" {
			http.Error(w, "category name is required", http.StatusBadRequest)
			return
		}

		createdCategory, err := CreateCategory(newCategory.Name, db.DB)
		if err != nil {
			http.Error(w, "unable to create category: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdCategory)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// IssuesCategoriesRequestById - обработчик для /api/categories/{id} (GET, PUT, DELETE)
func IssuesCategoriesRequestById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	idStr := vars["id"]

	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		category, err := GetCategoryByID(id, db.DB)
		if err != nil {
			http.Error(w, "unable to get category: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if category == nil {
			http.Error(w, "category not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(category)

	case http.MethodPut:
		var updatedCategory Category
		if err := json.NewDecoder(r.Body).Decode(&updatedCategory); err != nil {
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if updatedCategory.Name == "" {
			http.Error(w, "category name is required", http.StatusBadRequest)
			return
		}

		updatedCategory.ID = id
		result, err := UpdateCategory(id, updatedCategory.Name, db.DB)
		if err != nil {
			http.Error(w, "unable to update category: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if result == nil {
			http.Error(w, "category not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(result)

	case http.MethodDelete:
		err := DeleteCategory(id, db.DB)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "category not found", http.StatusNotFound)
				return
			}
			http.Error(w, "unable to delete category: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// IssuesGoodsRequest - обработчик для /api/goods (GET, POST)
func IssuesGoodsRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// Получаем все товары
		goods, err := GetGoods(db.DB)
		if err != nil {
			http.Error(w, "unable to get goods: "+err.Error(), http.StatusInternalServerError)
			return
		}
		
		// Фильтруем по category_id если передан параметр
		categoryIdStr := r.URL.Query().Get("category_id")
		if categoryIdStr != "" {
			categoryId, err := strconv.Atoi(categoryIdStr)
			if err != nil {
				http.Error(w, "invalid category_id format", http.StatusBadRequest)
				return
			}
			
			// Фильтруем товары по категории
			var filteredGoods []Goods
			for _, good := range goods {
				if good.CategoryID == categoryId {
					filteredGoods = append(filteredGoods, good)
				}
			}
			goods = filteredGoods
		}
		
		json.NewEncoder(w).Encode(goods)

	case http.MethodPost:
		var newGoods Goods
		if err := json.NewDecoder(r.Body).Decode(&newGoods); err != nil {
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Валидация
		if newGoods.Name == "" {
			http.Error(w, "goods name is required", http.StatusBadRequest)
			return
		}
		if newGoods.Unit == "" {
			http.Error(w, "goods unit is required", http.StatusBadRequest)
			return
		}
		if newGoods.Quantity < 0 {
			http.Error(w, "quantity cannot be negative", http.StatusBadRequest)
			return
		}
		if newGoods.CategoryID <= 0 {
			http.Error(w, "valid category_id is required", http.StatusBadRequest)
			return
		}

		// Проверяем существует ли категория
		_, err := GetCategoryByID(newGoods.CategoryID, db.DB)
		if err != nil {
			http.Error(w, "invalid category_id", http.StatusBadRequest)
			return
		}

		createdGoods, err := CreateGoods(newGoods, db.DB)
		if err != nil {
			http.Error(w, "unable to create goods: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdGoods)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// IssuesGoodsRequestById - обработчик для /api/goods/{id} (GET, PUT, DELETE)
func IssuesGoodsRequestById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	idStr := vars["id"]

	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		goods, err := GetGoodsByID(id, db.DB)
		if err != nil {
			http.Error(w, "unable to get goods: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if goods == nil {
			http.Error(w, "goods not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(goods)

	case http.MethodPut:
		var updatedGoods Goods
		if err := json.NewDecoder(r.Body).Decode(&updatedGoods); err != nil {
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Валидация
		if updatedGoods.Name == "" {
			http.Error(w, "goods name is required", http.StatusBadRequest)
			return
		}
		if updatedGoods.Unit == "" {
			http.Error(w, "goods unit is required", http.StatusBadRequest)
			return
		}
		if updatedGoods.Quantity < 0 {
			http.Error(w, "quantity cannot be negative", http.StatusBadRequest)
			return
		}
		if updatedGoods.CategoryID <= 0 {
			http.Error(w, "valid category_id is required", http.StatusBadRequest)
			return
		}

		updatedGoods.ID = id
		result, err := UpdateGoods(updatedGoods, db.DB)
		if err != nil {
			http.Error(w, "unable to update goods: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if result == nil {
			http.Error(w, "goods not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(result)

	case http.MethodDelete:
		err := DeleteGoods(id, db.DB)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "goods not found", http.StatusNotFound)
				return
			}
			http.Error(w, "unable to delete goods: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
