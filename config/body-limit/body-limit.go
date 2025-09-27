package bodylimit

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetBodyLimit(app *echo.Echo) {
	app.Use(middleware.BodyLimit("50M"))
}
