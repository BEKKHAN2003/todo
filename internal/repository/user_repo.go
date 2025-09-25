package repository

import (
	"context"
	"errors"
	"tasklist/db"
	"tasklist/internal/models"
)

type User interface {
	Create(ctx context.Context, user *models.User) (int, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
}

type UserRepo struct {
	db *db.Database
}

func NewUserRepo(db *db.Database) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *models.User) (int, error) {
	var userId int
	query := `INSERT INTO users (username, password) VALUES ($1,$2) RETURNING id`
	if err := r.db.Pool.QueryRow(ctx, query, user.Username, user.Password).Scan(&userId); err != nil {
		return userId, err
	}

	return userId, nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	query := `SELECT id, username, password FROM users WHERE username=$1`

	row := r.db.Pool.QueryRow(ctx, query, username)
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
