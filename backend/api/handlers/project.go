package handlers

import (
	"encoding/json"
	"net/http"
	"translatehub/api/models"
	"translatehub/api/services"
)

// ProjectHandler handles project creation and listing
func ProjectHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleCreateProject(w, r)
	case http.MethodGet:
		handleListProjects(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleCreateProject(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name       string `json:"name"`
		SourceLang string `json:"source_lang"`
		TargetLang string `json:"target_lang"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid project name"))
		return
	}
	project := models.Project{Name: req.Name, SourceLang: req.SourceLang, TargetLang: req.TargetLang}
	if err := services.CreateProject(&project); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create project in DB"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

func handleListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := services.ListProjects()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to list projects"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}
