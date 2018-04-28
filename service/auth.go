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
	Register(ctx context.Context, u domain.User) (domain.User, domain.Token, error)
	Login(ctx context.Context, username, password string) (domain.User, domain.Token, error)
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

func (s repositoryAuthService) Register(_ context.Context, u domain.User) (domain.User, domain.Token, error) {
	valErr := u.Validate()
	if valErr != nil {
		return domain.User{}, domain.Token{}, valErr
	}
	u.EncryptPassword()
	user, err := s.userRepo.Insert(u)
	if err != nil {
		return domain.User{}, domain.Token{}, err
	}

	token, tokenErr := s.tokenRepo.Insert(domain.Token{
		UserID:    user.ID,
		Token:     randToken(),
		ExpiresAt: time.Now().AddDate(0, 0, 30),
	})
	if tokenErr != nil {
		return domain.User{}, domain.Token{}, tokenErr
	}
	return user, token, nil
}

func (s repositoryAuthService) Login(_ context.Context, username, password string) (domain.User, domain.Token, error) {

	user, userErr := s.userRepo.FindByUsername(username)
	if userErr != nil {
		return domain.User{}, domain.Token{}, userErr
	}
	user.ClearPassword()

	passErr := user.ComparePassword(password)
	if passErr != nil {
		return domain.User{}, domain.Token{}, domain.ErrInvalidPassword
	}

	token, tokenErr := s.tokenRepo.Insert(domain.Token{
		UserID:    user.ID,
		Token:     randToken(),
		ExpiresAt: time.Now().AddDate(0, 0, 30),
	})
	if tokenErr != nil {
		return domain.User{}, domain.Token{}, tokenErr
	}

	return user, token, nil
}

func randToken() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
