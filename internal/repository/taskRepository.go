package repository

import (
	"context"
	"database/sql"
	"fmt"
	"task-manager/internal/contextkeys"
	"task-manager/internal/db"
	"task-manager/internal/helpers"
	"task-manager/internal/models"
)

type TaskRepository interface {
	Store(ctx context.Context, p models.TaskPayload) error
	Update(ctx context.Context, p models.UpdateTask) (models.Task, error)
	Show(id int) (models.Task, error)
	Index(uID int64) (models.TasksList, error)
	Delete(id int) error
	IsTaskOwner(uID int64, id int) (bool, error)
}

type taskRepository struct {
	d db.DB
}

func NewTaskRepository(d db.DB) TaskRepository {
	return &taskRepository{
		d: d,
	}
}

func (r taskRepository) Store(ctx context.Context, p models.TaskPayload) error {
	q, err := db.GetQuery(helpers.GetQueryPath("task/InsertTask.sql"))
	if err != nil {
		return fmt.Errorf("store: failed to read query: %v", err)
	}
	uID, ok := ctx.Value(contextkeys.UserID).(int64)
	if !ok || uID == 0 {
		return fmt.Errorf("store: failed to get user id")
	}

	tx, err := r.d.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store: failed to begin tx: %v", err)
	}
	_, err = tx.ExecContext(
		ctx,
		q,
		p.Name,
		p.Priority,
		p.Description,
		p.DueDate,
		p.CreatedAt,
		uID,
	)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("store: failed to insert a new task: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: failed to commit tx: %v", err)
	}

	return nil
}
func (r taskRepository) Update(ctx context.Context, p models.UpdateTask) (models.Task, error) {
	q, err := db.GetQuery(helpers.GetQueryPath("task/GetTask.sql"))
	if err != nil {
		return models.Task{}, fmt.Errorf("update: failed to read query: %v", err)
	}
	var t models.Task
	var desc sql.NullString

	q = fmt.Sprintf("%s for update", q)

	if err := r.d.QueryRow(q, p.ID).Scan(&t.Name, &t.Priority, &desc, &t.DueDate, &t.CreatedAt, &t.CreatedBy); err != nil {
		return models.Task{}, fmt.Errorf("update: failed to get task from db: %v", err)
	}
	if desc.Valid {
		t.Description = desc.String
	}

	t.Name = p.Name
	t.Priority = p.Priority
	t.Description = p.Description
	t.DueDate = p.DueDate
	t.ID = int(p.ID)

	uID := ctx.Value("uID").(int64)
	if uID != t.CreatedBy {
		return models.Task{}, fmt.Errorf("user not authorized for this action")
	}

	uq, err := db.GetQuery(helpers.GetQueryPath("task/UpdateTask.sql"))
	if err != nil {
		return models.Task{}, fmt.Errorf("update: failed to read update query, %v", err)
	}

	tx, err := r.d.BeginTx(ctx, nil)
	if err != nil {
		return models.Task{}, fmt.Errorf("update: failed to begin tx:%v", err)
	}

	_, err = tx.ExecContext(
		ctx,
		uq,
		t.Name,
		t.Priority,
		t.Description,
		t.DueDate,
		t.ID,
	)
	if err != nil {
		_ = tx.Rollback()
		return models.Task{}, fmt.Errorf("update: failed to execute update query: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return models.Task{}, fmt.Errorf("update: failed to commit tx: %v", err)
	}
	return t, nil
}
func (r taskRepository) Show(id int) (models.Task, error) {
	q, err := db.GetQuery(helpers.GetQueryPath("task/GetTask.sql"))
	if err != nil {
		return models.Task{}, fmt.Errorf("show: failed to read query:%v", err)
	}
	var t models.Task
	var desc sql.NullString

	err = r.d.QueryRow(q, id).Scan(&t.Name, &t.Priority, &desc, &t.DueDate, &t.CreatedAt, &t.CreatedBy)
	if err != nil {
		return models.Task{}, fmt.Errorf("show: failed to execute query:%v ", err)
	}
	if desc.Valid {
		t.Description = desc.String
	}

	return t, nil
}
func (r taskRepository) Index(uID int64) (models.TasksList, error) {
	q, err := db.GetQuery(helpers.GetQueryPath("task/GetTasksList.sql"))
	if err != nil {
		return models.TasksList{}, fmt.Errorf("index: failed to read query: %v", err)
	}
	var l models.TasksList

	rows, err := r.d.Query(q, uID)
	if err != nil {
		return models.TasksList{}, fmt.Errorf("index: failed to execute query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Task
		var d sql.NullString
		if err := rows.Scan(&t.ID, &t.Name, &t.Priority, &d, &t.DueDate, &t.CreatedAt, &t.CreatedBy); err != nil {
			return models.TasksList{}, fmt.Errorf("index: failed to read results: %v", err)
		}
		if d.Valid {
			t.Description = d.String
		}

		l.Tasks = append(l.Tasks, t)
	}

	if err := rows.Err(); err != nil {
		return models.TasksList{}, fmt.Errorf("index: query failed: %v", err)
	}

	return l, nil
}
func (r taskRepository) Delete(id int) error {
	q, err := db.GetQuery(helpers.GetQueryPath("task/DeleteTask.sql"))
	if err != nil {
		return fmt.Errorf("delete: failed to read query:%v", err)
	}

	_, err = r.d.Exec(q, id)
	if err != nil {
		return fmt.Errorf("delete: failed to execute query:%v", err)
	}

	return nil
}
func (r taskRepository) IsTaskOwner(uID int64, id int) (bool, error) {
	q, err := db.GetQuery(helpers.GetQueryPath("task/GetTaskOwner.sql"))
	if err != nil {
		return false, fmt.Errorf("isTaskOwner: failed to read query: %v", err)
	}

	var isOwner bool

	if err := r.d.QueryRow(q, uID, id).Scan(&isOwner); err != nil {
		return false, fmt.Errorf("isTaskOwner: failed to execute query: %v", err)
	}

	return isOwner, nil
}
