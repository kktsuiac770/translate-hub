package main

import (
	"log"
	"net/http"

	"translatehub/api/handlers"
)

func main() {
	http.HandleFunc("/tasks", handlers.TaskHandler)
	http.HandleFunc("/changes", handlers.SubmitChangeHandler)
	http.HandleFunc("/review", handlers.ReviewChangeHandler)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
