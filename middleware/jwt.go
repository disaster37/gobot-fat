package middleware

import (
	"github.com/disaster37/gobot-fat/login/usecase"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (m *GoMiddleware) IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if u := c.Get("user"); u != nil {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(*usecase.JwtCustomClaims)
			isAdmin := claims.Admin
			if !isAdmin {
				return echo.ErrUnauthorized
			}
			return next(c)
		}

		return echo.ErrUnauthorized
	}
}
