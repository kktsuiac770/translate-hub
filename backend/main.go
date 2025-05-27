package main

import (
	"log"
	"net/http"

	"translatehub/api/handlers"
	"translatehub/api/services"
)

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		h(w, r)
	}
}

func main() {
	if err := services.InitDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer services.CloseDB()

	// Project routes
	http.HandleFunc("/projects", withCORS(handlers.ProjectHandler))
	http.HandleFunc("/projects/", withCORS(handlers.TaskHandler)) // For project-specific tasks

	// Task routes
	http.HandleFunc("/tasks", withCORS(handlers.TaskHandler))
	http.HandleFunc("/tasks/", withCORS(handlers.TaskHandler))

	// Change routes
	http.HandleFunc("/changes", withCORS(handlers.SubmitChangeHandler))
	http.HandleFunc("/review", withCORS(handlers.ReviewChangeHandler))
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
