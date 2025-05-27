package services

import (
	"context"
	"fmt"
	"translatehub/api/models"
)

// CreateTask inserts a new task and its dialogues into the database
func CreateTask(task *models.Task) error {
	ctx := context.Background()
	tx, err := DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var taskID int
	fmt.Println("Creating task:", task.Name, "for project ID:", task.ProjectID)
	err = tx.QueryRow(ctx, `INSERT INTO tasks (creator, name, filename, status, project_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		task.Creator, task.Name, task.Filename, task.Status, task.ProjectID).Scan(&taskID)
	if err != nil {
		return err
	}
	for _, d := range task.Dialogues {
		var dialogueID int
		err = tx.QueryRow(ctx, `INSERT INTO dialogues (task_id, text, trans) VALUES ($1, $2, $3) RETURNING id`, taskID, d.Text, d.Trans).Scan(&dialogueID)
		if err != nil {
			return err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	task.ID = taskID
	return nil
}

// UpdateTask updates a task record in the database
func UpdateTask(task *models.Task) error {
	ctx := context.Background()
	_, err := DB.Exec(ctx, `UPDATE tasks SET status=$1 WHERE id=$2`, task.Status, task.ID)
	if err != nil {
		return err
	}
	for _, d := range task.Dialogues {
		_, err := DB.Exec(ctx, `UPDATE dialogues SET trans=$1 WHERE id=$2 AND task_id=$3`, d.Trans, d.ID, task.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddChange inserts a new change for a task
func AddChange(change *models.Change) error {
	ctx := context.Background()
	var changeID int
	err := DB.QueryRow(ctx, `INSERT INTO changes (task_id, dialogue_id, user, new_trans, status) VALUES ($1, $2, $3, $4, $5) RETURNING id`, change.TaskID, change.DialogueID, change.User, change.NewTrans, change.Status).Scan(&changeID)
	if err != nil {
		return err
	}
	change.ID = changeID
	return nil
}

// UpdateChange updates a change record in the database
func UpdateChange(change *models.Change) error {
	ctx := context.Background()
	_, err := DB.Exec(ctx, `UPDATE changes SET status=$1, new_trans=$2 WHERE id=$3`, change.Status, change.NewTrans, change.ID)
	return err
}

// ListProjectTasks returns all tasks for a specific project
func ListProjectTasks(projectID int) ([]models.Task, error) {
	ctx := context.Background()
	rows, err := DB.Query(ctx, `
		SELECT t.id, t.name, t.creator, t.filename, t.status, t.project_id 
		FROM tasks t 
		WHERE t.project_id = $1
		ORDER BY t.id DESC`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Name, &t.Creator, &t.Filename, &t.Status, &t.ProjectID); err != nil {
			return nil, err
		}
		// Get dialogues for this task
		dRows, err := DB.Query(ctx, `
			SELECT id, text, trans 
			FROM dialogues 
			WHERE task_id = $1 
			ORDER BY id`, t.ID)
		if err != nil {
			return nil, err
		}
		defer dRows.Close()

		for dRows.Next() {
			var d models.Dialogue
			if err := dRows.Scan(&d.ID, &d.Text, &d.Trans); err != nil {
				return nil, err
			}
			t.Dialogues = append(t.Dialogues, d)
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

// ListAllTasks returns all tasks in the system
func ListAllTasks() ([]models.Task, error) {
	ctx := context.Background()
	rows, err := DB.Query(ctx, `
		SELECT id, creator, filename, status, project_id 
		FROM tasks 
		ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Creator, &t.Filename, &t.Status, &t.ProjectID); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

// GetTask returns a specific task by ID with all its dialogues and changes
func GetTask(taskID int) (*models.Task, error) {
	ctx := context.Background()
	var task models.Task
	err := DB.QueryRow(ctx, `
		SELECT id, name, creator, filename, status, project_id 
		FROM tasks 
		WHERE id = $1`, taskID).Scan(
		&task.ID, &task.Name, &task.Creator, &task.Filename, &task.Status, &task.ProjectID)
	if err != nil {
		return nil, err
	}

	// Get dialogues
	rows, err := DB.Query(ctx, `
		SELECT id, text, trans 
		FROM dialogues 
		WHERE task_id = $1 
		ORDER BY id`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d models.Dialogue
		if err := rows.Scan(&d.ID, &d.Text, &d.Trans); err != nil {
			return nil, err
		}
		task.Dialogues = append(task.Dialogues, d)
	}

	// Get changes
	cRows, err := DB.Query(ctx, `
		SELECT id, dialogue_id, user, new_trans, status 
		FROM changes 
		WHERE task_id = $1 
		ORDER BY id DESC`, taskID)
	if err != nil {
		return nil, err
	}
	defer cRows.Close()

	for cRows.Next() {
		var c models.Change
		if err := cRows.Scan(&c.ID, &c.DialogueID, &c.User, &c.NewTrans, &c.Status); err != nil {
			return nil, err
		}
		c.TaskID = taskID
		task.Changes = append(task.Changes, c)
	}

	return &task, nil
}

// GetChange returns a specific change by ID
func GetChange(changeID int) (*models.Change, error) {
	ctx := context.Background()
	var change models.Change
	err := DB.QueryRow(ctx, `
		SELECT id, task_id, dialogue_id, user, new_trans, status 
		FROM changes 
		WHERE id = $1`, changeID).Scan(
		&change.ID, &change.TaskID, &change.DialogueID, &change.User, &change.NewTrans, &change.Status)
	if err != nil {
		return nil, err
	}
	return &change, nil
}

// UpdateChangeStatus updates the status of a change
func UpdateChangeStatus(changeID int, status string) error {
	ctx := context.Background()
	_, err := DB.Exec(ctx, `UPDATE changes SET status = $1 WHERE id = $2`, status, changeID)
	return err
}

// UpdateDialogueTranslation updates the translation of a dialogue
func UpdateDialogueTranslation(taskID int, dialogueID int, translation string) error {
	ctx := context.Background()
	_, err := DB.Exec(ctx, `
		UPDATE dialogues 
		SET trans = $1 
		WHERE task_id = $2 AND id = $3`,
		translation, taskID, dialogueID)
	return err
}
