package services

import (
	"context"
	"translatehub/api/models"
)

func CreateProject(project *models.Project) error {
	ctx := context.Background()
	var id int
	err := DB.QueryRow(ctx, `INSERT INTO projects (name, source_lang, target_lang) VALUES ($1, $2, $3) RETURNING id`, project.Name, project.SourceLang, project.TargetLang).Scan(&id)
	if err != nil {
		return err
	}
	project.ID = id
	return nil
}

func ListProjects() ([]models.Project, error) {
	ctx := context.Background()
	rows, err := DB.Query(ctx, `SELECT id, name, source_lang, target_lang FROM projects`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var projects []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.SourceLang, &p.TargetLang); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

// GetProject retrieves a project by ID
func GetProject(id int) (*models.Project, error) {
	ctx := context.Background()
	var project models.Project
	err := DB.QueryRow(ctx, `SELECT id, name, source_lang, target_lang FROM projects WHERE id = $1`, id).Scan(
		&project.ID, &project.Name, &project.SourceLang, &project.TargetLang)
	if err != nil {
		return nil, err
	}
	return &project, nil
}
