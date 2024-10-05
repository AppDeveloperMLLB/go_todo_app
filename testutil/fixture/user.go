package fixture

import (
	"math/rand"
	"strconv"
	"time"

	"example.com/sample/go_todo_app/entity"
)

func User(u *entity.User) *entity.User {
	result := &entity.User{
		ID:        entity.UserID(rand.Int()),
		Name:      "User" + strconv.Itoa(rand.Int())[:5],
		Password:  "password",
		Role:      "admin",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if u == nil {
		return result
	}
	if u.ID != 0 {
		result.ID = u.ID
	}
	if u.Name != "" {
		result.Name = u.Name
	}
	if u.Password != "" {
		result.Password = u.Password
	}
	if u.Role != "" {
		result.Role = u.Role
	}
	if !u.CreatedAt.IsZero() {
		result.CreatedAt = u.CreatedAt
	}
	if !u.UpdatedAt.IsZero() {
		result.UpdatedAt = u.UpdatedAt
	}
	return result
}
