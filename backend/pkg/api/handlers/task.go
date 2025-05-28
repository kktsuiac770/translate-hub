package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"translatehub/pkg/api/models"
	"translatehub/pkg/api/services"
)

// TaskHandler handles task-related requests
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the URL path to extract project ID if present
	// Expected formats:
	// /tasks - List all tasks or create a new task
	// /projects/{projectId}/tasks - List or create tasks for a specific project
	// /tasks/{taskId} - Get, update, or delete a specific task

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if this is a project-specific tasks request
	if len(pathParts) >= 4 && pathParts[1] == "projects" {
		projectID, err := strconv.Atoi(pathParts[2])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		handleProjectTasks(w, r, projectID)
		return
	}

	// Check if this is a specific task request
	if len(pathParts) >= 3 && pathParts[1] == "tasks" {
		taskID, err := strconv.Atoi(pathParts[2])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.Method == http.MethodGet {
			handleGetTask(w, r, taskID)
			return
		}
	}

	// Handle general task requests
	switch r.Method {
	case http.MethodGet:
		handleListTasks(w, r)
	case http.MethodPost:
		handleCreateTask(w, r, 0) // 0 means no specific project
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleProjectTasks(w http.ResponseWriter, r *http.Request, projectID int) {
	switch r.Method {
	case http.MethodGet:
		handleListProjectTasks(w, r, projectID)
	case http.MethodPost:
		handleCreateTask(w, r, projectID)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleListProjectTasks(w http.ResponseWriter, r *http.Request, projectID int) {
	tasks, err := services.ListProjectTasks(projectID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to list tasks: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func handleListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := services.ListAllTasks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to list tasks: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func handleCreateTask(w http.ResponseWriter, r *http.Request, projectID int) {
	// Parse the multipart form
	err := r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Failed to parse form: %v", err)
		return
	}

	// Get form data
	name := r.FormValue("name")
	creator := r.FormValue("creator")
	if name == "" || creator == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Name and creator are required"))
		return
	}

	// Get the file
	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Failed to get file: %v", err)
		return
	}
	defer file.Close()

	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to read file: %v", err)
		return
	}

	lines := splitLines(string(content))

	// Get project details to get source and target languages
	project, err := services.GetProject(projectID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to get project details: %v", err)
		return
	}

	// Create dialogues with initial translations
	dialogues := make([]models.Dialogue, len(lines))
	for i, line := range lines {
		dialogues[i] = models.Dialogue{
			ID:    i + 1,
			Text:  line,
			Trans: services.GeminiTranslate(line, project.SourceLang, project.TargetLang),
		}
	}

	// Create task
	task := &models.Task{
		Name:      name,
		Creator:   creator,
		Status:    "new",
		ProjectID: projectID,
		Filename:  header.Filename,
		Dialogues: dialogues,
	}

	err = services.CreateTask(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to create task: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func splitLines(s string) []string {
	return strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
}

// SubmitChangeHandler allows a user to submit a translation change for a dialogue
func SubmitChangeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	taskID, _ := strconv.Atoi(r.FormValue("task_id"))
	dialogueID, _ := strconv.Atoi(r.FormValue("dialogue_id"))
	user := r.FormValue("user")
	newTrans := r.FormValue("new_trans")

	change := &models.Change{
		TaskID:     taskID,
		DialogueID: dialogueID,
		User:       user,
		NewTrans:   newTrans,
		Status:     "pending",
	}

	err := services.AddChange(change)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to add change: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(change)
}

// ReviewChangeHandler allows the task creator to approve/reject a change
func ReviewChangeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	taskID, _ := strconv.Atoi(r.FormValue("task_id"))
	changeID, _ := strconv.Atoi(r.FormValue("change_id"))
	status := r.FormValue("status") // approved or rejected
	user := r.FormValue("user")

	task, err := services.GetTask(taskID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Task not found: %v", err)
		return
	}

	if task.Creator != user {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Only the task creator can review changes"))
		return
	}

	err = services.UpdateChangeStatus(changeID, status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to update change: %v", err)
		return
	}

	if status == "approved" {
		change, err := services.GetChange(changeID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Failed to get change details: %v", err)
			return
		}

		err = services.UpdateDialogueTranslation(taskID, change.DialogueID, change.NewTrans)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Failed to update dialogue: %v", err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Change reviewed"))
}

func handleGetTask(w http.ResponseWriter, r *http.Request, taskID int) {
	task, err := services.GetTask(taskID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Failed to get task: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}
