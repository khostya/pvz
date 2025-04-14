package domain

import "errors"

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserNotFound    = errors.New("user not found")
	ErrProductNotFound = errors.New("product not found")

	ErrEmployeeOnly  = errors.New("method is only available to employees")
	ErrModeratorOnly = errors.New("method is only available to moderators")

	ErrThereIsNoInProgressReception = errors.New("there is no in progress reception")

	ErrPreviousReceptionIsNotClosed = errors.New("previous reception is not closed")
	ErrReceptionAlreadyClosed       = errors.New("reception already closed")

	ErrReceptionNotFound = errors.New("reception not found")
)
