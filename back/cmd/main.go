package main

import (
	"html/template"
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
	mux.Use(corsMiddleware)

	// Статические файлы (frontend)
	//mux.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("../front"))))
	//mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./front/"))))
	// API endpoints for categories

	mux.HandleFunc("/api/categories", api.IssuesCategoriesRequest).Methods("GET", "POST", "OPTIONS")
	mux.HandleFunc("/api/categories/{id}", api.IssuesCategoriesRequestById).Methods("GET", "PUT", "DELETE", "OPTIONS")

	// API endpoints for goods
	mux.HandleFunc("/api/goods", api.IssuesGoodsRequest).Methods("GET", "POST", "OPTIONS")
	mux.HandleFunc("/api/goods/{id}", api.IssuesGoodsRequestById).Methods("GET", "PUT", "DELETE", "OPTIONS")

	// mux.PathPrefix("/").Handler(http.FileServer(http.Dir("./front")))
	mux.HandleFunc("/", home)
	mux.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./front/"))))
	log.Println("Server starting on http://localhost:8080")
	log.Println("Frontend available at: http://localhost:8080/")
	log.Println("API available at: http://localhost:8080/api/")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	Tmpl, _ := template.ParseFiles("./front/index.html")
	Tmpl.Execute(w, nil)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
