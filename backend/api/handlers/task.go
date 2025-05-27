package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"translatehub/api/models"
	"translatehub/api/services"
)

var (
	tasks    = make([]models.Task, 0)
	taskMux  sync.Mutex
	taskID   = 1
	changeID = 1
)

// TaskHandler handles task creation and listing
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleCreateTask(w, r)
	case http.MethodGet:
		handleListTasks(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleCreateTask(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("File upload error"))
		return
	}
	defer file.Close()
	content, _ := io.ReadAll(file)
	lines := splitLines(string(content))
	sourceLang := r.FormValue("source_lang")
	targetLang := r.FormValue("target_lang")
	dialogues := make([]models.Dialogue, len(lines))
	for i, line := range lines {
		dialogues[i] = models.Dialogue{
			ID:    i + 1,
			Text:  line,
			Trans: services.GeminiTranslate(line, sourceLang, targetLang),
		}
	}
	taskMux.Lock()
	task := models.Task{
		ID:        taskID,
		Creator:   r.FormValue("user"),
		Filename:  handler.Filename,
		Dialogues: dialogues,
		Status:    "open",
	}
	err = services.CreateTask(&task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create task in DB"))
		return
	}
	tasks = append(tasks, task)
	taskID++
	taskMux.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func handleListTasks(w http.ResponseWriter, r *http.Request) {
	taskMux.Lock()
	defer taskMux.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
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

	taskMux.Lock()
	defer taskMux.Unlock()
	for i, t := range tasks {
		if t.ID == taskID {
			change := models.Change{
				ID:         changeID,
				TaskID:     taskID,
				DialogueID: dialogueID,
				User:       user,
				NewTrans:   newTrans,
				Status:     "pending",
			}
			err := services.AddChange(&change)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Failed to add change in DB"))
				return
			}
			tasks[i].Changes = append(tasks[i].Changes, change)
			changeID++
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(change)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Task not found"))
}

// ReviewChangeHandler allows the task creator to approve/reject a change
func ReviewChangeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	taskID, _ := strconv.Atoi(r.FormValue("task_id"))
	changeIDVal, _ := strconv.Atoi(r.FormValue("change_id"))
	status := r.FormValue("status") // approved or rejected
	user := r.FormValue("user")

	taskMux.Lock()
	defer taskMux.Unlock()
	for i, t := range tasks {
		if t.ID == taskID && t.Creator == user {
			for j, c := range t.Changes {
				if c.ID == changeIDVal {
					tasks[i].Changes[j].Status = status
					if status == "approved" {
						for k, d := range t.Dialogues {
							if d.ID == c.DialogueID {
								tasks[i].Dialogues[k].Trans = c.NewTrans
								err := services.UpdateTask(&tasks[i])
								if err != nil {
									w.WriteHeader(http.StatusInternalServerError)
									w.Write([]byte("Failed to update task in DB"))
									return
								}
							}
						}
					}
					err := services.UpdateChange(&tasks[i].Changes[j])
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte("Failed to update change in DB"))
						return
					}
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("Change reviewed"))
					return
				}
			}
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Change or task not found, or not authorized"))
}
