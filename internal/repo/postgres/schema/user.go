package schema

import (
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
)

type (
	User struct {
		ID       uuid.UUID `db:"users.id"`
		Email    string    `db:"users.email"`
		Password string    `db:"users.password"`
		Role     string    `db:"users.role"`
	}
)

func NewUser(d *domain.User) *User {
	return &User{
		ID:       d.ID,
		Email:    d.Email,
		Password: d.Password,
		Role:     string(d.Role),
	}
}

func NewDomainUser(d *User) *domain.User {
	return &domain.User{
		ID:       d.ID,
		Email:    d.Email,
		Password: d.Password,
		Role:     domain.Role(d.Role),
	}
}

func (User) TableName() string {
	return "users"
}

func (u User) InsertColumns() []string {
	return []string{"id", "email", "password", "role"}
}

func (u User) Columns() []string {
	return []string{"users.id as \"users.id\"", "users.email as \"users.email\"",
		"users.password as \"users.password\"", "users.role as \"users.role\""}
}

func (u User) Values() []any {
	return []any{u.ID, u.Email, u.Password, u.Role}
}
