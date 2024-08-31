package handler

import (
	"net/http"

	"example.com/sample/go_todo_app/entity"
	"example.com/sample/go_todo_app/store"
)

type ListTask struct {
	Store *store.TaskStore
}

type task struct {
	ID     entity.TaskID `json:"id"`
	Title  string        `json:"title"`
	Status string        `json:"status"`
}

func (lt *ListTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tasks := lt.Store.All()
	rsp := []task{}
	for _, t := range tasks {
		rsp = append(rsp, task{
			ID:     t.ID,
			Title:  t.Title,
			Status: string(t.Status),
		})
	}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
