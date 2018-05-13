package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/project-tilas/svc-auth/domain"
	"github.com/project-tilas/svc-auth/repository"
)

type AuthService interface {
	Register(ctx context.Context, u domain.User) (*domain.User, *domain.Token, error)
	Login(ctx context.Context, username, password string) (*domain.User, *domain.Token, error)
	LoginWithToken(ctx context.Context, userID, token string) (*domain.User, *domain.Token, error)
}

type repositoryAuthService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
}

func NewRepositoryAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
) AuthService {
	return &repositoryAuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (s repositoryAuthService) Register(_ context.Context, u domain.User) (*domain.User, *domain.Token, error) {
	valErr := u.Validate()
	if valErr != nil {
		return nil, nil, valErr
	}
	u.EncryptPassword()
	user, err := s.userRepo.Insert(u)
	if err != nil {
		return nil, nil, err
	}

	token, tokenErr := s.tokenRepo.Insert(domain.Token{
		UserID:    user.ID,
		Token:     randToken(),
		ExpiresAt: time.Now().AddDate(0, 0, 30),
	})
	if tokenErr != nil {
		return nil, nil, tokenErr
	}
	return &user, &token, nil
}

func (s repositoryAuthService) Login(_ context.Context, username, password string) (*domain.User, *domain.Token, error) {

	user, userErr := s.userRepo.FindByUsername(username)
	fmt.Println(userErr)
	if userErr != nil {
		return nil, nil, userErr
	}

	passErr := user.ComparePassword(password)
	if passErr != nil {
		return nil, nil, domain.ErrInvalidPassword
	}
	user.ClearPassword()

	token, tokenErr := s.tokenRepo.Insert(domain.Token{
		UserID:    user.ID,
		Token:     randToken(),
		ExpiresAt: time.Now().AddDate(0, 0, 30),
	})
	if tokenErr != nil {
		return nil, nil, tokenErr
	}

	return &user, &token, nil
}

func (s repositoryAuthService) LoginWithToken(_ context.Context, userID, token string) (*domain.User, *domain.Token, error) {

	_, findErr := s.tokenRepo.FindByUserIDAndToken(userID, token)
	if findErr != nil {
		return nil, nil, findErr
	}

	user, userErr := s.userRepo.FindByID(userID)
	if userErr != nil {
		return nil, nil, userErr
	}
	user.ClearPassword()

	insert, insertErr := s.tokenRepo.Insert(domain.Token{
		UserID:    user.ID,
		Token:     randToken(),
		ExpiresAt: time.Now().AddDate(0, 0, 30),
	})
	if insertErr != nil {
		return nil, nil, insertErr
	}

	deleteErr := s.tokenRepo.Remove(token)
	return &user, &insert, deleteErr
}

func randToken() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
