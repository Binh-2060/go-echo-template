package controllers

import (
	"net/http"

	"github.com/Binh-2060/go-echo-template/api/presenters"
	"github.com/Binh-2060/go-echo-template/api/schema/requestbody"
	"github.com/Binh-2060/go-echo-template/api/services"
	"github.com/Binh-2060/go-echo-template/api/validators"
	jwtpkg "github.com/Binh-2060/go-echo-template/pkg/jwt-pkg"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GetUserController(c echo.Context) error {
	users, err := services.GetUserService(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, presenters.ResponseSuccess(users))
}

func CreateUserController(c echo.Context) error {
	var req requestbody.UserRequestBody
	var err error
	if err := validators.ParseAndValidateBody(c, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = services.CreateUserService(ctx, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, presenters.ResponseSuccess("SUCCESS"))
}

func SingTokenController(c echo.Context) error {
	token, err := jwtpkg.SignToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, presenters.ResponseSuccess(token))
}

func UserAuthController(c echo.Context) error {
	var user = c.Get("user").(jwt.MapClaims)
	return c.JSON(http.StatusOK, presenters.ResponseSuccess(user))
}
