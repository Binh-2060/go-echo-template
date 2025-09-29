package routes

import (
	"github.com/Binh-2060/go-echo-template/api/controllers"
	jwtpkg "github.com/Binh-2060/go-echo-template/pkg/jwt-pkg"
	"github.com/labstack/echo/v4"
)

func SetUserRoutes(router *echo.Group) {
	router.GET("/getData", controllers.GetUserController)
	router.POST("/create", controllers.CreateUserController)
	router.GET("/signToken", controllers.SingTokenController)
	//test with token middleware
	router.GET("/me", controllers.UserAuthController, jwtpkg.VerifyToken)
}
