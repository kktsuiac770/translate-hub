package services

import (
	"context"
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
	err = tx.QueryRow(ctx, `INSERT INTO tasks (creator, filename, status) VALUES ($1, $2, $3) RETURNING id`, task.Creator, task.Filename, task.Status).Scan(&taskID)
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
