package routes

import (
	"github.com/Binh-2060/go-echo-template/api/controllers"
	"github.com/labstack/echo/v4"
)

func SetUserRoutes(router *echo.Group) {
	router.GET("/getData", controllers.GetUserController)
	router.POST("/create", controllers.CreateUserController)
}
