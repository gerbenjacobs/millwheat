package services

import (
	"context"

	"github.com/google/uuid"
)

func UserFromContext(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(CtxKeyUserID).(uuid.UUID)
	if ok {
		return v
	}

	return uuid.UUID{}
}

func TownFromContext(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(CtxKeyTownID).(uuid.UUID)
	if ok {
		return v
	}

	return uuid.UUID{}
}
