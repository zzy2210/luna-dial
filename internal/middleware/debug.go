package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

var DebugMode = false

// DebugLoggerMiddleware 打印每个请求和响应（仅debug模式）
func DebugLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !DebugMode {
			return next(c)
		}

		request := c.Request()
		method := request.Method
		path := request.URL.Path
		var bodyCopy []byte
		if request.Body != nil {
			bodyCopy, _ = ioutil.ReadAll(request.Body)
			request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyCopy))
		}
		fmt.Printf("[DEBUG][REQ] %s %s\nHeaders: %v\nBody: %s\n", method, path, request.Header, string(bodyCopy))

		// 捕获响应
		rec := &responseRecorder{ResponseWriter: c.Response().Writer, buf: new(bytes.Buffer)}
		c.Response().Writer = rec
		start := time.Now()
		err := next(c)
		duration := time.Since(start)

		status := c.Response().Status
		respBody := rec.buf.String()
		fmt.Printf("[DEBUG][RESP] %s %s %d (%v)\nBody: %s\n", method, path, status, duration, respBody)

		return err
	}
}

type responseRecorder struct {
	ResponseWriter http.ResponseWriter
	buf            *bytes.Buffer
}

func (r *responseRecorder) Header() http.Header {
	return r.ResponseWriter.Header()
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.buf.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
}
