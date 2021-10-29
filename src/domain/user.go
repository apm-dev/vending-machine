package domain

import (
	"context"
	"time"
)

const (
	Seller Role = "seller"
	Buyer  Role = "buyer"
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

func NewUser(uname, passwd string, role Role) *User {
	return &User{
		Username: uname,
		Password: passwd,
		Role:     role,
		Deposit:  0,
	}
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
	Insert(ctx context.Context, u User) uint
	FindById(ctx context.Context, id uint) *User
	Update(ctx context.Context, u *User)
	Delete(ctx context.Context, id uint)
}

type JwtRepository interface {
	Exists(ctx context.Context, token string) (bool, error)
	Insert(ctx context.Context, userId uint, token string, ttl time.Duration) error
	DeleteTokensOfUserExcept(ctx context.Context, userId uint, exceptionToken string) error
}
