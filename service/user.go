package service

import (
	"context"

	"github.com/project-tilas/svc-auth/domain"
	"github.com/project-tilas/svc-auth/repository"
)

// UserService is the core business logic for authentication
type UserService interface {
	RegisterUser(ctx context.Context, u domain.User) (domain.User, error)
	DeleteUser(ctx context.Context, id string) error
	Login(ctx context.Context, username, password string) (domain.User, error)
}

type repositoryUserService struct {
	repo repository.UserRepository
}

func NewRepositoryUserService(repo repository.UserRepository) UserService {
	return &repositoryUserService{
		repo: repo,
	}
}

func (s repositoryUserService) RegisterUser(_ context.Context, u domain.User) (domain.User, error) {
	user, err := s.repo.Insert(u)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s repositoryUserService) DeleteUser(_ context.Context, id string) error {
	return nil //TODO
}

func (s repositoryUserService) Login(_ context.Context, username, password string) (domain.User, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
