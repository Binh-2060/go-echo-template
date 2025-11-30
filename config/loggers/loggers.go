package loggers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

type ErrorResponse struct {
	Code    int
	Message string
}

type CustomResponseWriter struct {
	echo.Response
	body *bytes.Buffer
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

// Write captures the response body
func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.Response.Write(b)
}

// format requestjson with content type
func formatJsonRequest(c echo.Context) interface{} {
	var requestJSON interface{}
	contentType := c.Request().Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		// Parse multipart form (set a reasonable size limit, e.g., 10MB)
		if err := c.Request().ParseMultipartForm(10 << 20); err != nil {
			requestJSON = map[string]interface{}{"error": "failed to parse multipart form"}
		} else {
			// Initialize request body map
			requestJSON = make(map[string]interface{})
			// Add form fields
			for key, values := range c.Request().MultipartForm.Value {
				if len(values) == 1 {
					requestJSON.(map[string]interface{})[key] = values[0]
				} else {
					requestJSON.(map[string]interface{})[key] = values
				}
			}
			requestJSON.(map[string]interface{})["is_file_uploads"] = len(c.Request().MultipartForm.File) > 0
		}
	} else {
		// Handle JSON or other content types
		var requestBody []byte
		if c.Request().Body != nil {
			requestBody, _ = io.ReadAll(c.Request().Body)
			// Restore the body for downstream handlers
			c.Request().Body = io.NopCloser(bytes.NewBuffer(requestBody))
			// Parse request body as JSON
			if len(requestBody) > 0 {
				if err := json.Unmarshal(requestBody, &requestJSON); err != nil {
					// Fallback to empty map if JSON is invalid
					requestJSON = map[string]interface{}{}
				}
			} else {
				requestJSON = map[string]interface{}{}
			}
		} else {
			requestJSON = map[string]interface{}{}
		}
	}

	return requestJSON
}

func formatErrorMessage(err error) string {
	var msg = "-"
	if err == nil {
		return msg
	}

	if he, ok := err.(*echo.HTTPError); ok {
		if m, ok := he.Message.(string); ok {
			msg = m
		} else {
			msg = err.Error()
		}
	} else {
		msg = err.Error()
	}

	return msg
}

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

		// Calculate latency
		latency := time.Since(start)

		var errorRes = catchError(err)
		var event *zerolog.Event
		if errorRes.Code == http.StatusOK {
			event = Logger.Info()
		} else {
			event = Logger.Error()
		}

		// Prepare base event
		//you can custom as you want to loggings
		event = Logger.Info().
			Str("method", c.Request().Method).
			Str("path", c.Request().URL.Path).
			Str("query_params", c.Request().URL.RawQuery).
			Str("remote_ip", c.RealIP()).
			Str("request_id", c.Response().Header().Get(echo.HeaderXRequestID)).
			Str("latency_sec", fmt.Sprint(latency.Seconds(), "s")).
			Int("status", errorRes.Code).
			Str("error", errorRes.Message).
			// Int64("bytes_out", c.Response().Size).
			Interface("request_body", requestJSON)

		if user := c.Get("user"); user != nil {
			event.Interface("user", user)
		}

		event.Msg("Request completed")
		return err
	}
}
