package service

import (
	"context"
	"fmt"

	"example.com/sample/go_todo_app/entity"
	"example.com/sample/go_todo_app/store"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUser struct {
	DB   store.Queryer
	Repo UserRegister
}

func (r *RegisterUser) RegisterUser(
	ctx context.Context,
	name, password, role string) (*entity.User, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	u := &entity.User{
		Name:     name,
		Password: string(pw),
		Role:     role,
	}

	if err := r.Repo.RegisterUser(ctx, r.DB, u); err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}
	return u, nil
}
