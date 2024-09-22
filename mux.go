package main

import (
	"context"
	"net/http"

	"example.com/sample/go_todo_app/clock"
	"example.com/sample/go_todo_app/config"
	"example.com/sample/go_todo_app/handler"
	"example.com/sample/go_todo_app/service"
	"example.com/sample/go_todo_app/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	mux := chi.NewRouter()
	mux.HandleFunc(
		"/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		},
	)
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	r := store.Repository{Clocker: clock.RealClocker{}}
	v := validator.New()
	at := &handler.AddTask{
		Service: &service.AddTask{
			DB:   db,
			Repo: &r,
		},
		Validator: v,
	}
	mux.Post("/tasks", at.ServeHTTP)
	lt := &handler.ListTask{
		Service: &service.ListTask{
			DB:   db,
			Repo: &r,
		},
	}
	mux.Get("/tasks", lt.ServeHTTP)
	ru := &handler.RegisterUser{
		Service:   &service.RegisterUser{DB: db, Repo: &r},
		Validator: v,
	}
	mux.Post("/register", ru.ServeHTTP)
	return mux, cleanup, nil
}
