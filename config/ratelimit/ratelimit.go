package ratelimit

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func SetEchoRateLimit(app *echo.Echo) {

	app.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(1000), Burst: 30, ExpiresIn: 1 * time.Minute},
		),
		IdentifierExtractor: func(c echo.Context) (string, error) {
			return c.RealIP(), nil
		},
		ErrorHandler: func(c echo.Context, err error) error {
			res := map[string]interface{}{
				"timestamp": time.Now().Format("2006-01-02-15-04-05"),
				"status":    0,
				"items":     nil,
				"error":     "Too many requests",
			}
			return c.JSON(http.StatusTooManyRequests, res)
		},
	}))

}
