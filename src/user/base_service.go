package user

import (
	"context"
	"sync"
	"time"

	"github.com/apm-dev/vending-machine/domain"
)

type Service struct {
	ur             domain.UserRepository
	jr             domain.JwtRepository
	jwt            *JWTManager
	depositTimeout time.Duration
	depositLock    *sync.RWMutex
}

var UserService *Service

func InitService(
	ur domain.UserRepository,
	jr domain.JwtRepository,
	jwt *JWTManager,
	depositTimeout time.Duration,
) domain.UserService {
	if UserService == nil {
		UserService = &Service{
			ur: ur, jr: jr, jwt: jwt,
			depositTimeout: depositTimeout,
			depositLock:    &sync.RWMutex{},
		}
	}
	return UserService
}

func (s *Service) Update(ctx context.Context, passwd string) error {
	panic("not implemented") // TODO: Implement
}

func (s *Service) Delete(ctx context.Context) ([]uint, error) {
	panic("not implemented") // TODO: Implement
}

func (s *Service) Get(ctx context.Context, id uint) (*domain.User, error) {
	panic("not implemented") // TODO: Implement
}

func (s *Service) List(ctx context.Context) ([]domain.User, error) {
	panic("not implemented") // TODO: Implement
}
