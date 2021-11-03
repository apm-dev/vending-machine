package user

import (
	"context"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/pkg/errors"
)

func (s *Service) fetchContextUser(ctx context.Context) (*domain.User, error) {
	const op string = "user.helper.fetchContextUser"

	uid, err := domain.UserIdFromContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	user, err := s.ur.FindById(ctx, uid)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	return user, nil
}
