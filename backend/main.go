package main

import (
	"log"
	"net/http"

	"translatehub/api/handlers"
	"translatehub/api/services"
)

func main() {
	if err := services.InitDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer services.CloseDB()

	http.HandleFunc("/tasks", handlers.TaskHandler)
	http.HandleFunc("/changes", handlers.SubmitChangeHandler)
	http.HandleFunc("/review", handlers.ReviewChangeHandler)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
