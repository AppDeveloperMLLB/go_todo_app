package store

import (
	"context"

	"example.com/sample/go_todo_app/entity"
)

func (r *Repository) ListTasks(
	ctx context.Context, db Queryer,
) (entity.Tasks, error) {
	tasks := entity.Tasks{}
	sql := `
		SELECT id, title, status, created_at
		FROM tasks;
	`
	if err := db.SelectContext(ctx, &tasks, sql); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *Repository) AddTask(
	ctx context.Context, db Queryer, t *entity.Task,
) error {
	t.CreatedAt = r.Clocker.Now()
	sql := `
		INSERT INTO tasks (title, status, created_at)
		VALUES ($1, $2, $3) returning id;
	`
	err := db.QueryRowxContext(
		ctx, sql, t.Title, t.Status, t.CreatedAt,
	).Scan(&t.ID)
	if err != nil {
		return err
	}
	return nil
}
