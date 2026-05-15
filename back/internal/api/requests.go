package api

import (
	//"database/sql"
	"encoding/json"
	"invent/back/internal/db"
	"net/http"

	"github.com/gorilla/mux"
)

func IssuesCategoriesRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		categories, err := GetCategories(db.DB)
		if err != nil {
			http.Error(w, "unable to get categories: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(categories)
	} else if r.Method == http.MethodPost {
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
		CreateCategory(newCategory.Name, db.DB)
	}
	return
}

func IssuesCategoriesRequestById(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		vars := mux.Vars(r)
		id := vars["id"]

		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		DeleteCategory(id, db.DB)
	}
	return
}

func IssuesGoodsRequest(w http.ResponseWriter, r *http.Request) {

	return
}

func IssuesGoodsRequestById(w http.ResponseWriter, r *http.Request) {

	return
}
