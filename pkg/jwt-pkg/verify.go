package jwtpkg

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var publicKey []byte

func VerifyToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token format")
		}

		tokenStr := parts[1]
		if publicKey == nil {
			publicKey, err = os.ReadFile("app.rsa.pub")
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
		}

		key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("failed to authenticate with provided token")
			}
			return key, nil
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "failed to authenticate with provided token")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user", claims)
			return next(c)
		} else {
			return echo.NewHTTPError(http.StatusUnauthorized, "failed to authenticate with provided token")
		}
	}

}
