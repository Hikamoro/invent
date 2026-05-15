package main

import (
	"invent/back/internal/api"
	"invent/back/internal/db"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

func main() {
	db.InitDB()
	defer db.DB.Close()

	mux := mux.NewRouter()
	//r := mux.NewRouter()
	mux.Use(corsMiddleware)

	mux.HandleFunc("/api/categories", api.IssuesCategoriesRequest).Methods("GET", "OPTIONS")
	mux.HandleFunc("/api/categories/{id}", api.IssuesCategoriesRequestById).Methods("GET", "OPTIONS")
	mux.HandleFunc("/api/goods", api.IssuesGoodsRequest).Methods("GET", "OPTIONS")
	mux.HandleFunc("/api/goods/{id}", api.IssuesGoodsRequestById).Methods("GET", "OPTIONS")

	err := http.ListenAndServe(":8080", mux) // поднятие сервера
	if err != nil {
		log.Fatal(err)
	}
}

func ins(w http.ResponseWriter, r *http.Request) {}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Для production лучше указывать конкретный origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

//func IssuesHandler(w http.ResponseWriter, r *http.Request) {
// 	origin := r.Header.Get("Origin")
// 	if origin == "https://trusted.com" {
// 		w.Header().Set("Access-Control-Allow-Origin", origin)
// 	}
// 	w.Header().Set("Access-Control-Allow-Credentials", "true") // если нужны куки
// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 	if r.Method == http.MethodOptions {
// 		w.WriteHeader(http.StatusNoContent)
// 		return
// 	}
// }
