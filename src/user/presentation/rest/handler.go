package rest

import (
	"net/http"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/apm-dev/vending-machine/pkg/httputil"
	"github.com/apm-dev/vending-machine/user/presentation/rest/requests"
	"github.com/labstack/echo"
)

type UserHandler struct {
	us domain.UserService
}

func InitUserHandler(e *echo.Echo, auth *echo.Group, us domain.UserService) {
	h := &UserHandler{us}
	e.POST("/register", h.Register)
	e.POST("/login", h.Login)

	auth.POST("/logout/all", h.LogoutAll)
}

func (h *UserHandler) Register(c echo.Context) error {
	rr := new(requests.Register)
	if err := c.Bind(rr); err != nil {
		return c.JSON(http.StatusBadRequest, httputil.MakeResponse(
			http.StatusBadRequest, err.Error(), nil,
		))
	}
	if err := c.Validate(rr); err != nil {
		return c.JSON(http.StatusBadRequest, httputil.MakeResponse(
			http.StatusBadRequest, err.Error(), nil,
		))
	}

	token, err := h.us.Register(
		c.Request().Context(),
		rr.Username, rr.Password,
		domain.Role(rr.Role),
	)
	if err != nil {
		status := httputil.StatusCode(err)
		return c.JSON(status, httputil.MakeResponse(
			status, err.Error(), nil,
		))
	}

	return c.JSON(http.StatusOK, httputil.MakeResponse(
		http.StatusOK, "welcome "+rr.Username, echo.Map{"token": token},
	))
}

func (h *UserHandler) Login(c echo.Context) error {
	rr := new(requests.Login)
	if err := c.Bind(rr); err != nil {
		return c.JSON(http.StatusBadRequest, httputil.MakeResponse(
			http.StatusBadRequest, err.Error(), nil,
		))
	}
	if err := c.Validate(rr); err != nil {
		return c.JSON(http.StatusBadRequest, httputil.MakeResponse(
			http.StatusBadRequest, err.Error(), nil,
		))
	}

	token, activeSessions, err := h.us.Login(
		c.Request().Context(),
		rr.Username, rr.Password,
	)
	if err != nil {
		status := httputil.StatusCode(err)
		return c.JSON(status, httputil.MakeResponse(
			status, err.Error(), nil,
		))
	}
	// notify the user if there was another active session
	msg := "welcome " + rr.Username
	if activeSessions {
		msg += ", there is another active session."
	}
	return c.JSON(http.StatusOK, httputil.MakeResponse(
		http.StatusOK, msg, echo.Map{"token": token},
	))
}

func (h *UserHandler) LogoutAll(c echo.Context) error {
	err := h.us.TerminateActiveSessions(c.Request().Context())
	if err != nil {
		status := httputil.StatusCode(err)
		return c.JSON(status, httputil.MakeResponse(
			status, err.Error(), nil,
		))
	}

	return c.JSON(http.StatusOK, httputil.MakeResponse(
		http.StatusOK,
		"All other active sessions have been terminated.",
		nil,
	))
}
