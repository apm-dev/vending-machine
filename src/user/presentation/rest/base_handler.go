package rest

import (
	"github.com/apm-dev/vending-machine/domain"
	"github.com/labstack/echo"
)

type UserHandler struct {
	us domain.UserService
}

func InitUserHandler(e *echo.Echo, auth *echo.Group, us domain.UserService) {
	h := &UserHandler{us}
	// public routes
	e.POST("/register", h.Register)
	e.POST("/login", h.Login)
	// auth required routes
	auth.POST("/logout/all", h.LogoutAll)
	auth.POST("/deposit", h.Deposit)
	auth.POST("/reset", h.ResetDeposit)
}
