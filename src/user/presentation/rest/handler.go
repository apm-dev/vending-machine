package rest

import (
	"net/http"
	"strings"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/apm-dev/vending-machine/pkg/httputil"
	"github.com/apm-dev/vending-machine/user/presentation/rest/requests"
	"github.com/labstack/echo"
)

type UserHandler struct {
	us domain.UserService
}

func InitUserHandler(e *echo.Echo, us domain.UserService) {
	h := &UserHandler{us}
	e.POST("/register", h.Register)
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
		domain.Role(strings.ToUpper(rr.Role)),
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
