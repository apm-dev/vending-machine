package domain

import (
	"context"
	"time"
)

type (
	Role string
)

const (
	Seller Role = "seller"
	Buyer  Role = "buyer"
)

type User struct {
	Id       int64
	Username string
	Password string
	Role     Role
	Deposit  int64
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
	// Register create new user and return id, jwt token or error
	Register(ctx context.Context, uname, pass string, role Role) (int64, string, error)
	// Login check credentials, generate and return jwt token and a boolean
	// which says is there another active session using this account or not
	Login(ctx context.Context, uname, pass string) (string, bool, error)
	// Authorize parse jwt token and return related user
	Authorize(ctx context.Context, token string) (*User, error)
	// TerminateActiveSessions terminate all other active sessions
	TerminateActiveSessions(token string) error
}

type UserRepository interface {
	FindById(ctx context.Context, id int64) *User
	Insert(ctx context.Context, u User) int64
	Update(ctx context.Context, u *User)
	Delete(ctx context.Context, id int64)
}

type JwtRepository interface {
	Exists(ctx context.Context, token string) (bool, error)
	Insert(ctx context.Context, userId int64, token string, ttl time.Duration) error
	DeleteTokensOfUserExcept(ctx context.Context, userId int64, exceptionToken string) error
}
