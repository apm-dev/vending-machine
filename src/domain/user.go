package domain

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	USER_ID_CONTEXT_KEY = "userId"

	SELLER Role = "seller"
	BUYER  Role = "buyer"
)

type (
	Role string
)

type User struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Role     Role   `json:"role"`
	Deposit  uint   `json:"deposit"`
}

func NewUser(uname, passwd string, role Role) (*User, error) {
	const op string = "domain.user.NewUser"
	// storing hash of password for security reasons
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return &User{
		Username: uname,
		Password: string(hash),
		Role:     role,
		Deposit:  0,
	}, nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) AddDeposit(coin Coin) {
	u.Deposit += uint(coin)
}

func UserIdOfContext(ctx context.Context) (uint, error) {
	uid, ok := ctx.Value(USER_ID_CONTEXT_KEY).(uint)
	if !ok || uid == 0 {
		return 0, ErrUserNotFound
	}
	return uid, nil
}

type UserService interface {
	// Register creates new user and return jwt token or error
	Register(ctx context.Context, uname, pass string, role Role) (string, error)
	// Login checks credentials, generate and return jwt token and a boolean
	// which says is there another active session using this account or not
	Login(ctx context.Context, uname, pass string) (string, bool, error)
	// Authorize parses jwt token and return related user
	Authorize(ctx context.Context, token string) (uint, error)
	// TerminateActiveSessions terminates all other active sessions
	TerminateActiveSessions(ctx context.Context, token string) error
	// Deposit increases buyer(user) deposit
	Deposit(ctx context.Context, coin Coin) (uint, error)
	// ResetDeposit reset buyer(user) deposits back to zero
	ResetDeposit(ctx context.Context) error
}

type UserRepository interface {
	Insert(ctx context.Context, u User) (uint, error)
	FindById(ctx context.Context, id uint) (*User, error)
	FindByUsername(ctx context.Context, un string) (*User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id uint) error
}

type JwtRepository interface {
	Insert(ctx context.Context, userId uint, token string, ttl time.Duration) error
	Exists(ctx context.Context, token string) (bool, error)
	UserTokensCount(ctx context.Context, uid uint) (uint, error)
	DeleteTokensOfUserExcept(ctx context.Context, userId uint, exceptionToken string) error
}
