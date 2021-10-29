package user

import (
	"context"

	"github.com/apm-dev/vending-machine/domain"
)

type Service struct {
	ur domain.UserRepository
	jr domain.JwtRepository
}

func InitService(ur domain.UserRepository, jr domain.JwtRepository) domain.UserService {
	return &Service{
		ur: ur,
		jr: jr,
	}
}

// Register creates new user and return jwt token or error
func (s *Service) Register(ctx context.Context, uname string, pass string, role domain.Role) (string, error) {
	panic("not implemented") // TODO: Implement
}

// Login checks credentials, generate and return jwt token and a boolean
// which says is there another active session using this account or not
func (s *Service) Login(ctx context.Context, uname string, pass string) (string, bool, error) {
	panic("not implemented") // TODO: Implement
}

// Authorize parses jwt token and return related user
func (s *Service) Authorize(ctx context.Context, token string) (uint, error) {
	panic("not implemented") // TODO: Implement
}

// TerminateActiveSessions terminates all other active sessions
func (s *Service) TerminateActiveSessions(ctx context.Context, token string) error {
	panic("not implemented") // TODO: Implement
}

// Deposit increases buyer(user) deposit
func (s *Service) Deposit(ctx context.Context, coin domain.Coin) (uint, error) {
	panic("not implemented") // TODO: Implement
}

// ResetDeposit reset buyer(user) deposits back to zero
func (s *Service) ResetDeposit(ctx context.Context) error {
	panic("not implemented") // TODO: Implement
}
