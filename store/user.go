package store

import (
	"context"
	"fmt"

	"example.com/sample/go_todo_app/entity"
)

func (r *Repository) RegisterUser(ctx context.Context, db Queryer, u *entity.User) error {
	u.CreatedAt = r.Clocker.Now()
	u.UpdatedAt = r.Clocker.Now()
	const query = `
		INSERT INTO users (name, password, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var id entity.UserID
	err := db.QueryRowxContext(ctx, query, u.Name, u.Password, u.Role).Scan(&id)
	if err != nil {
		// var pgErr pq.Error
		// if errors.As(err, &pgErr) {
		// 	if pgErr.Code == "23505" {
		// 		return fmt.Errorf("user already exists: %w", err)
		// 	}
		// }
		return fmt.Errorf("failed to insert user: %w", err)
	}
	u.ID = id
	return nil
}
