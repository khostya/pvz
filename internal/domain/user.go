package domain

import "github.com/google/uuid"

const (
	UserRoleModerator Role = "moderator"
	UserRoleEmployee  Role = "employee"
)

type (
	Role string

	User struct {
		ID       uuid.UUID `json:"id"`
		Email    string    `json:"email"`
		Password string    `json:"password"`
		Role     Role      `json:"role"`
	}

	Token string
)
