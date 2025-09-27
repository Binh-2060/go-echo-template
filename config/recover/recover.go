package recover

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetRecoverMiddleware(app *echo.Echo) {
	app.Use(middleware.Recover())
}
