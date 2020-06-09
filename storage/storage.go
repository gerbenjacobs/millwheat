package storage

import (
	"context"

	"github.com/google/uuid"

	app "github.com/gerbenjacobs/millwheat"
	"github.com/gerbenjacobs/millwheat/game"
)

type UserStorage interface {
	Create(ctx context.Context, user *app.User) error
	Read(ctx context.Context, userID uuid.UUID) (*app.User, error)
	Login(ctx context.Context, email, password string) (*app.User, error)
	Update(ctx context.Context, user *app.User) (*app.User, error)
}

type TownStorage interface {
	Get(ctx context.Context, id uuid.UUID) (*game.Town, error)
}
