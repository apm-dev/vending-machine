package middlewares

import (
	"net/http"
	"strings"

	"github.com/apm-dev/vending-machine/pkg/httputil"
	"github.com/labstack/echo"
)

func (m *UserMiddleware) JwtAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.JSON(http.StatusUnauthorized, httputil.MakeResponse(
				http.StatusUnauthorized,
				"Authorization header is required",
				nil,
			))
		}

		token = strings.TrimSpace(strings.Replace(token, "Bearer", "", 1))
		uid, err := m.us.Authorize(c.Request().Context(), token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, httputil.MakeResponse(
				http.StatusUnauthorized,
				err.Error(),
				nil,
			))
		}
		c.Set("userId", uid)

		return next(c)
	}
}
