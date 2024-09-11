package service

import (
	"context"
	"fmt"

	"example.com/sample/go_todo_app/entity"
	"example.com/sample/go_todo_app/store"
)

type AddTask struct {
	DB   store.Queryer
	Repo TaskAdder
}

func (a *AddTask) AddTask(ctx context.Context, title string) (*entity.Task, error) {
	t := &entity.Task{
		Title:  title,
		Status: entity.TaskStatusTodo,
	}
	err := a.Repo.AddTask(ctx, a.DB, t)
	if err != nil {
		return nil, fmt.Errorf("failed to add task: %w", err)
	}
	return t, nil
}
