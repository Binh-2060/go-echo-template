package routes

import "github.com/labstack/echo/v4"

func SetRoutes(echo *echo.Group) {
	userRoutes := echo.Group("/users")
	SetUserRoutes(userRoutes)
}
