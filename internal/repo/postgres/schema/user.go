package schema

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
	"time"
)

type (
	schemaModel struct {
		CreatedAt time.Time           `db:"created_at"`
		UpdatedAt sql.Null[time.Time] `db:"updated_at"`
		DeletedAt sql.Null[time.Time] `db:"deleted_at"`
	}

	User struct {
		schemaModel
		ID       uuid.UUID `db:"id"`
		Email    string    `db:"email"`
		Password string    `db:"password"`
		Role     string    `db:"role"`
	}
)

func NewUser(d *domain.User) *User {
	return &User{
		schemaModel: schemaModel{
			CreatedAt: time.Now(),
		},
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

func (u User) Columns() []string {
	return []string{"id", "email", "password", "role", "created_at", "updated_at", "deleted_at"}
}

func (u User) Values() []any {
	return []any{u.ID, u.Email, u.Password, u.Role, u.CreatedAt, u.UpdatedAt, u.DeletedAt}
}

func (u User) UpdateColumns() []string {
	return []string{"id", "email", "password", "role", "updated_at"}
}

func (u User) UpdateValues() []any {
	return []any{u.ID, u.Email, u.Password, u.Role, time.Now()}
}
