package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	app "github.com/gerbenjacobs/millwheat"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *app.User) error {
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO users (id, email, password, token, createdAt, updatedAt) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	uid, _ := user.ID.MarshalBinary()
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(password) // replace the actual password with the hashed version

	_, err = stmt.ExecContext(ctx, uid, user.Email, password, user.Token, user.CreatedAt, user.UpdatedAt)
	merr, ok := err.(*mysql.MySQLError)
	if ok && merr.Number == 1062 {
		return app.ErrUserEmailUniqueness
	}

	return err
}

func (r *UserRepository) Read(ctx context.Context, userID uuid.UUID) (*app.User, error) {
	uid, _ := userID.MarshalBinary()
	row := r.db.QueryRowContext(ctx, "SELECT id, email, password, token, createdAt, updatedAt FROM users WHERE id = ?", uid)

	user := new(app.User)
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Token, &user.CreatedAt, &user.UpdatedAt)
	switch {
	case err == sql.ErrNoRows:
		return nil, fmt.Errorf("user with ID %q not found: %w", userID, app.ErrUserNotFound)
	case err != nil:
		return nil, fmt.Errorf("unknown error while scanning user: %v", err)
	}

	return user, nil
}

func (r *UserRepository) Login(ctx context.Context, email, password string) (*app.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, password FROM users WHERE email = ?", email)

	var id uuid.UUID
	var hashedPassword []byte
	err := row.Scan(&id, &hashedPassword)
	switch {
	case err == sql.ErrNoRows:
		return nil, fmt.Errorf("user with email %q not found: %w", email, app.ErrEmailNotFound)
	case err != nil:
		return nil, fmt.Errorf("unknown error while scanning user: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		return nil, fmt.Errorf("%w: %v", app.ErrWrongPassword, err)
	}

	return r.Read(ctx, id)
}

func (r *UserRepository) Update(ctx context.Context, user *app.User) (*app.User, error) {
	uid, _ := user.ID.MarshalBinary()

	query := "UPDATE users SET email = ?, token = ?, createdAt = ?, updatedAt = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, user.Email, user.Token, user.CreatedAt, user.UpdatedAt, uid)

	return user, err
}
