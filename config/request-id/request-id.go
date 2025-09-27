package requestid

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetRequestID(app *echo.Echo) {
	app.Use(middleware.RequestID())
}
