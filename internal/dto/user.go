package dto

import "github.com/khostya/pvz/internal/domain"

type (
	LoginUserParam struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	RegisterUserParam struct {
		Email    string      `json:"email" validate:"required"`
		Password string      `json:"password" validate:"required"`
		Role     domain.Role `json:"role" validate:"required"`
	}

	DummyLoginUserParam struct {
		Role domain.Role `json:"role" validate:"required"`
	}
)
