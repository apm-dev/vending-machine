package user

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/apm-dev/vending-machine/pkg/logger"
	"github.com/pkg/errors"
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
		}
	}
	return UserService
}

// Register creates new user and return jwt token or error
func (s *Service) Register(ctx context.Context, uname, pass string, role domain.Role) (string, error) {
	const op string = "user.service.Register"
	// create domain user object
	user, err := domain.NewUser(uname, pass, role)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return "", domain.ErrInternalServer
	}
	// persist user
	_, err = s.ur.Insert(ctx, *user)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return "", domain.ErrUserAlreadyExists
		}
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return "", domain.ErrInternalServer
	}

	// generate jwt token with user claims
	token, err := s.jwt.Generate(user)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return "", domain.ErrInternalServer
	}

	logger.Log(logger.INFO, fmt.Sprintf("%s registered", uname))

	return token, nil
}

// Login checks credentials, generate and return jwt token and a boolean
// which says is there another active session using this account or not
func (s *Service) Login(ctx context.Context, uname, pass string) (string, bool, error) {
	const op string = "user.service.Login"
	// fetch user from db
	user, err := s.ur.FindByUsername(ctx, uname)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", false, domain.ErrUserNotFound
		}
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return "", false, domain.ErrInternalServer
	}
	// check password
	if !user.CheckPassword(pass) {
		logger.Log(logger.INFO, errors.Wrap(domain.ErrWrongCredentials, user.Username).Error())
		return "", false, domain.ErrWrongCredentials
	}
	// generate jwt token with user claims
	token, err := s.jwt.Generate(user)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return "", false, domain.ErrInternalServer
	}

	// persist jwt token
	err = s.jr.Insert(ctx, user.Id, token, s.jwt.tokenDuration)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return "", false, domain.ErrInternalServer
	}

	// check is there any other active session or not
	count, err := s.jr.UserTokensCount(ctx, user.Id)
	if err != nil {
		logger.Log(logger.WARN, errors.Wrap(err, op).Error())
	}

	logger.Log(logger.INFO, fmt.Sprintf(
		"%s logged-in", user.Username,
	))

	return token, count > 1, nil
}

// Authorize parses jwt token and return related user
func (s *Service) Authorize(ctx context.Context, token string) (uint, error) {
	const op string = "user.service.Authorize"

	// verify and get claims of token
	claims, err := s.jwt.Verify(token)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidToken) {
			logger.Log(logger.INFO, errors.Wrap(err, op).Error())
			return 0, domain.ErrInvalidToken
		}
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return 0, domain.ErrInternalServer
	}

	// check token existans in db to see is it still active or not
	exists, err := s.jr.Exists(ctx, token)
	if err != nil {
		return 0, domain.ErrInternalServer
	}
	if !exists {
		return 0, domain.ErrInvalidToken
	}

	return claims.Id, nil
}

// TerminateActiveSessions terminates all other active sessions
func (s *Service) TerminateActiveSessions(ctx context.Context, token string) error {
	const op string = "user.service.TerminateActiveSessions"

	uid, err := domain.UserIdOfContext(ctx)
	if err != nil {
		logger.Log(logger.WARN, errors.Wrap(err, op).Error())
		return domain.ErrInternalServer
	}

	err = s.jr.DeleteTokensOfUserExcept(ctx, uid, token)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return domain.ErrInternalServer
	}

	return nil
}

// Deposit increases buyer(user) deposit
func (s *Service) Deposit(ctx context.Context, coin domain.Coin) (uint, error) {
	const op string = "user.service.Deposit"
	
	ctx, cancel := context.WithTimeout(ctx, s.depositTimeout)
	defer cancel()

	if !coin.IsValid() {
		return 0, domain.ErrInvalidCoin
	}

	uid, err := domain.UserIdOfContext(ctx)
	if err != nil {
		return 0, domain.ErrInternalServer
	}

	user, err := s.ur.FindById(ctx, uid)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return 0, domain.ErrInternalServer
	}

	if user.Role != domain.BUYER {
		return 0, domain.ErrUnauthorized
	}

	// use locks because of concurrent requests
	// we make sure to add user deposits without conflict
	// there are better ways to handle this kind of issue
	// but it's just for POC
	s.depositLock.Lock()
	defer s.depositLock.Unlock()

	user.AddDeposit(coin)

	err = s.ur.Update(ctx, user)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return 0, domain.ErrInternalServer
	}

	return user.Deposit, nil
}

// ResetDeposit reset buyer(user) deposits back to zero
func (s *Service) ResetDeposit(ctx context.Context) error {
	panic("not implemented") // TODO: Implement
}
