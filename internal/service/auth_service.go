package service

import (
	"context"
	"errors"
	"tasklist/pkg/config"

	"golang.org/x/crypto/bcrypt"

	"tasklist/internal/models"
	"tasklist/internal/repository"
	"tasklist/pkg/auth"
)

type Auth interface {
	Register(ctx context.Context, username, password string) (string, error)
	Login(ctx context.Context, username, password string) (string, error)
}
type AuthService struct {
	repo repository.User
	cfg  *config.Config
}

func NewAuthService(r repository.User, cfg *config.Config) *AuthService {
	return &AuthService{
		repo: r,
		cfg:  cfg,
	}
}

func (s *AuthService) Register(ctx context.Context, username, password string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("username or password is empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := models.User{
		Username: username,
		Password: string(hash),
	}
	userId, err := s.repo.Create(ctx, &user)
	if err != nil {
		return "", err
	}

	token, err := auth.GenerateToken(userId, s.cfg)
	if err != nil {
		return "", err
	}

	return token, nil
}
func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("username or password is empty")
	}
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("password incorrect")
	}
	return auth.GenerateToken(user.ID, s.cfg)
}
