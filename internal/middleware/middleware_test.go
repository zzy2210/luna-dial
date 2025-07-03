package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"okr-web/internal/types"
)

func TestErrorHandler(t *testing.T) {
	e := echo.New()
	e.HTTPErrorHandler = ErrorHandler()

	tests := []struct {
		name         string
		err          error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "AppError",
			err:          types.ErrUserNotFound,
			expectedCode: http.StatusNotFound,
			expectedBody: `{"success":false,"error":"USER_NOT_FOUND","message":"用户不存在","code":404}`,
		},
		{
			name:         "HTTPError",
			err:          echo.NewHTTPError(http.StatusBadRequest, "Bad Request"),
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"success":false,"error":"HTTP_ERROR","message":"Bad Request","code":400}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			e.HTTPErrorHandler(tt.err, c)

			assert.Equal(t, tt.expectedCode, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}

func TestCORSConfig(t *testing.T) {
	e := echo.New()
	e.Use(CORSConfig())

	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
}

func TestLoggerConfig(t *testing.T) {
	e := echo.New()
	e.Use(LoggerConfig())

	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "test", rec.Body.String())
}

func TestRecoverConfig(t *testing.T) {
	e := echo.New()
	e.Use(RecoverConfig())

	e.GET("/panic", func(c echo.Context) error {
		panic("test panic")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()

	// 这不应该导致程序崩溃
	assert.NotPanics(t, func() {
		e.ServeHTTP(rec, req)
	})

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
