package loggers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type CustomResponseWriter struct {
	echo.Response
	body *bytes.Buffer
}

// Write captures the response body
func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.Response.Write(b)
}

// LoggingMiddleware logs request and response bodies in structured JSON format
func SetEchoLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Start time for logging
		start := time.Now()

		// Initialize request body variable
		var requestJSON interface{}
		contentType := c.Request().Header.Get("Content-Type")

		// Handle multipart/form-data
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
			// Reset the body to allow downstream handlers to read it
			c.Request().Body = io.NopCloser(bytes.NewReader([]byte{}))
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

		// Wrap response writer to capture response body
		// buf := new(bytes.Buffer)
		// customWriter := &CustomResponseWriter{
		// 	Response: *c.Response(),
		// 	body:     buf,
		// }
		// c.SetResponse(echo.NewResponse(customWriter, c.Echo()))

		// Call the next handler
		err := next(c)

		// Parse response body as JSON
		// var responseJSON interface{}
		// if buf.Len() > 0 {
		// 	if err := json.Unmarshal(buf.Bytes(), &responseJSON); err != nil {
		// 		// Fallback to empty map if JSON is invalid
		// 		responseJSON = map[string]interface{}{}
		// 	}
		// } else {
		// 	responseJSON = map[string]interface{}{}
		// }

		var error string
		if err != nil {
			error = err.Error()
		} else {
			error = "-"
		}
		// Prepare log data
		logData := map[string]interface{}{
			"timestamp":     time.Now().Format(time.RFC3339Nano),
			"method":        c.Request().Method,
			"request_body":  requestJSON,
			"path":          c.Request().URL.Path,
			"query_params":  c.Request().URL.RawQuery,
			"status":        c.Response().Status,
			"remote_ip":     c.RealIP(),
			"response_body": "-",
			"latency":       fmt.Sprint(time.Since(start).Seconds(), " s"),
			"error":         error,
		}

		// Encode and log as JSON
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false) // Keeps '&' instead of \u0026
		_ = enc.Encode(logData)

		return err
	}
}
