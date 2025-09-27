package cors

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetCorsMiddlwares(app *echo.Echo) {
	app.Use(middleware.CORS())
}
