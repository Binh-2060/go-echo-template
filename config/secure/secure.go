package secure

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetSecureMiddilware(app *echo.Echo) {
	app.Use(middleware.Secure())
}
