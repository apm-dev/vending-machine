package user

import (
	"context"

	"github.com/apm-dev/vending-machine/domain"
)

func (s *Service) fetchContextUser(ctx context.Context) (*domain.User, error) {
	uid, err := domain.UserIdFromContext(ctx)
	if err != nil {
		return nil, domain.ErrInternalServer
	}

	user, err := s.ur.FindById(ctx, uid)
	if err != nil {
		return nil, domain.ErrInternalServer
	}
	return user, nil
}
