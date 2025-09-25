package repository

import (
	"context"
	"errors"
	"tasklist/db"

	"tasklist/internal/models"
)

var ErrTaskNotFound = errors.New("task not found")

type Task interface {
	Create(ctx context.Context, userId int, task models.TaskRequest) error
	Complete(ctx context.Context, taskId int) error
	List(ctx context.Context, userId int) ([]models.Task, error)
	GetByID(ctx context.Context, taskId, userId int) (*models.Task, error)
	Update(ctx context.Context, taskId, userId int, task models.TaskRequest) error
	Delete(ctx context.Context, id int) error

	ListAll(ctx context.Context) ([]models.Task, error)
	MarkOverdued(ctx context.Context, taskId int) error
}
type TaskRepo struct {
	db *db.Database
}

func NewTaskRepo(db *db.Database) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) Create(ctx context.Context, userId int, task models.TaskRequest) error {
	query := `INSERT INTO tasks (user_id,title, deadline) VALUES ($1, $2, $3) RETURNING id`
	_, err := r.db.Pool.Exec(ctx, query, userId, task.Title, task.Deadline)
	return err
}
func (r *TaskRepo) Complete(ctx context.Context, taskId int) error {
	query := `update tasks set completed=true where id=$1`
	_, err := r.db.Pool.Exec(ctx, query, taskId)
	return err
}

func (r *TaskRepo) List(ctx context.Context, userId int) ([]models.Task, error) {
	query := `SELECT id, title, completed, deadline, is_overdue
		FROM tasks
		where user_id = $1
		ORDER BY id DESC`

	rows, err := r.db.Pool.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Completed, &task.Deadline, &task.IsOverdue); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepo) ListAll(ctx context.Context) ([]models.Task, error) {
	query := `SELECT id, title, completed, deadline, is_overdue FROM tasks where is_overdue=false and completed=false`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Completed, &task.Deadline, &task.IsOverdue); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepo) GetByID(ctx context.Context, id, userId int) (*models.Task, error) {
	var task models.Task
	query := `SELECT id, title, completed, deadline, is_overdue
		FROM tasks
		WHERE id=$1 and user_id=$2`
	err := r.db.Pool.QueryRow(ctx, query, id, userId).Scan(&task.ID, &task.Title, &task.Completed, &task.Deadline, &task.IsOverdue)
	if err != nil {
		return nil, ErrTaskNotFound
	}
	return &task, nil
}

func (r *TaskRepo) Update(ctx context.Context, taskId, userId int, task models.TaskRequest) error {
	query := `UPDATE tasks SET title=$1, deadline=$2 WHERE id=$3 and user_id=$4`
	rows, err := r.db.Pool.Exec(ctx, query, task.Title, task.Deadline, taskId, userId)
	if err != nil {
		return err
	}
	if rows.RowsAffected() == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func (r *TaskRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tasks WHERE id=$1`
	rows, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if rows.RowsAffected() == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func (r *TaskRepo) MarkOverdued(ctx context.Context, taskId int) error {
	query := `update tasks set is_overdue=true where id=$1`
	rows, err := r.db.Pool.Exec(ctx, query, taskId)
	if err != nil {
		return err
	}
	if rows.RowsAffected() == 0 {
		return ErrTaskNotFound
	}
	return nil
}
