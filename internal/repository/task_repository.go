package repository

import (
	"database/sql"
	"taskforge/internal/models"
	"time"
)

type TaskRepository interface {
	Create(task *models.Task) error
	GetByID(id int) (*models.Task, error)
	GetByUserID(userID int) ([]models.Task, error)
	Update(task *models.Task) error
	Delete(id int) error
}

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{
		db: db,
	}
}

func (r *taskRepository) Create(task *models.Task) error {
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	query := `
		INSERT INTO tasks (user_id, title, description, status, priority, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		task.UserID,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.DueDate,
		task.CreatedAt,
		task.UpdatedAt,
	).Scan(&task.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *taskRepository) GetByID(id int) (*models.Task, error) {
	query := `
		SELECT id, user_id, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks WHERE id = $1
	`

	var task models.Task
	err := r.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	return &task, nil
}

func (r *taskRepository) GetByUserID(userID int) ([]models.Task, error) {
	query := `SELECT * FROM tasks WHERE user_id = $1`

	var tasks []models.Task

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *taskRepository) Update(task *models.Task) error {
	task.UpdatedAt = time.Now()

	query := `
		UPDATE tasks
		SET title = $1, description  = $2, status = $3, priority = $4, 
			due_date = $5, updated_at = $6
		WHERE id = $7
	`

	result, err := r.db.Exec(
		query, 
		task.Title,
        task.Description,
        task.Status,
        task.Priority,
        task.DueDate,
        task.UpdatedAt,
        task.ID,
	)
	if err != nil {
		return err
	}

	rowsUpdated, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsUpdated == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (r *taskRepository) Delete(id int) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsDeleted, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsDeleted == 0 {
		return ErrTaskNotFound
	}

	return nil
}