package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func (m *GoMiddleware) IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if u := c.Get("user"); u != nil {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			isAdmin := claims["admin"].(bool)
			if isAdmin == false {
				return echo.ErrUnauthorized
			}
			return next(c)
		}

		return echo.ErrUnauthorized
	}
}
