package services

import (
	"context"

	"github.com/google/uuid"

	app "github.com/gerbenjacobs/millwheat"
	"github.com/gerbenjacobs/millwheat/game"
)

type UserService interface {
	Add(ctx context.Context, user *app.User) error
	User(ctx context.Context, userID uuid.UUID) (*app.User, error)
	Login(ctx context.Context, email, password string) (*app.User, error)
	Update(ctx context.Context, user *app.User) (*app.User, error)
}

type TownService interface {
	Town(ctx context.Context, id uuid.UUID) (*game.Town, error)
}
