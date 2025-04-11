package appctx

import (
	"context"
	"github.com/google/uuid"
	"github.com/khostya/pvz/internal/domain"
	"github.com/labstack/echo/v4"
)

const (
	userID  = "user_id"
	isDummy = "is_dummy"
	role    = "role"
)

func SetUserID(ctx context.Context, uuid uuid.UUID) context.Context {
	return set(ctx, userID, uuid)
}

func SetIsDummy(ctx context.Context, is bool) context.Context {
	return set(ctx, isDummy, is)
}

func SetRole(ctx context.Context, r domain.Role) context.Context {
	return set(ctx, role, r)
}

func GetIsDummy(ctx context.Context) (bool, bool) {
	return get[bool](ctx, isDummy)
}

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	return get[uuid.UUID](ctx, userID)
}

func GetRole(ctx context.Context) (domain.Role, bool) {
	return get[domain.Role](ctx, role)
}

func set(ctx context.Context, key string, v any) context.Context {
	return context.WithValue(ctx, key, v)
}

func get[T any](ctx context.Context, key string) (T, bool) {
	v, ok := ctx.Value(key).(T)
	return v, ok
}

func EchoGetRole(ctx echo.Context) (domain.Role, bool) {
	return eGet[domain.Role](ctx, role)
}

func eGet[T any](ctx echo.Context, key string) (T, bool) {
	v, ok := ctx.Get(key).(T)
	return v, ok
}

func SetEcho(ctx context.Context, eCtx echo.Context) {
	r, ok := GetRole(ctx)
	if ok {
		eCtx.Set(role, r)
	}

	uid, ok := GetUserID(ctx)
	if ok {
		eCtx.Set(userID, uid)
	}

	isD, ok := GetIsDummy(ctx)
	if ok {
		eCtx.Set(isDummy, isD)
	}
}

func EchoSetRole(eCtx echo.Context, r domain.Role) {
	eCtx.Set(role, r)
}
