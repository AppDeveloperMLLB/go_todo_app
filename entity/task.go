package entity

import "time"

type TaskID int64
type TaskStatus string

const (
	TaskStatusTodo TaskStatus = "todo"
	TaskStatusDoing
	TaskStatusDone TaskStatus = "done"
)

type Task struct {
	ID        TaskID     `json:"id" db:"id"`
	Title     string     `json:"title" db:"title"`
	Status    TaskStatus `json:"status" db:"status"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

type Tasks []*Task
