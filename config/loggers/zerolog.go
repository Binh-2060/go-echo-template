package loggers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

type ErrorResponse struct {
	Code    int
	Message string
}

func InitLogger() {
	log.Println("#### Begin Load Loggers Config ####")
	Logger = zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()
}

func catchError(err error) ErrorResponse {
	//set by default
	var res = ErrorResponse{
		Code:    http.StatusOK,
		Message: "-",
	}

	if err == nil {
		return res
	}

	if httpError, ok := err.(*echo.HTTPError); ok {
		if m, ok := httpError.Message.(string); ok {
			res.Message = m
		} else {
			res.Message = err.Error()
		}
		res.Code = httpError.Code
	} else {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
	}

	return res
}

// Middleware for structured request logging using zerolog
func SetEchoZeroLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Skip logging OPTIONS preflight requests
		if c.Request().Method == http.MethodOptions {
			return next(c)
		}

		start := time.Now()
		requestJSON := formatJsonRequest(c)

		// Run handler
		err := next(c)
		latency := time.Since(start)
		var errorRes = catchError(err)

		// Prepare base event
		event := Logger.Info().
			Str("method", c.Request().Method).
			Str("path", c.Request().URL.Path).
			Str("query_params", c.Request().URL.RawQuery).
			Str("remote_ip", c.RealIP()).
			Str("request_id", c.Response().Header().Get(echo.HeaderXRequestID)).
			Float64("latency_sec", latency.Seconds()).
			Int("status", errorRes.Code).
			Str("error", errorRes.Message).
			Interface("user", c.Get("user")).
			Interface("request_body", requestJSON)

		event.Msg("Request completed")
		return err
	}
}
