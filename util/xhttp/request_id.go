package xhttp

import (
	"github.com/labstack/echo"
	"golang.org/x/net/context"
	"strconv"
	"time"
)

const (
	RequestIDHeader = "X-RequestID"
	RequestIDKey    = "_RequestID_"
)

func RequestID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := c.Request().Header().Get(RequestIDHeader)
		if requestID == "" {
			requestID = strconv.FormatInt(time.Now().UnixNano(), 10)
		}
		ctx := c.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		ctx = context.WithValue(ctx, RequestIDKey, requestID)
		c.SetContext(ctx)
		c.Response().Header().Set(RequestIDHeader, requestID)
		return next(c)
	}
}

func TakeRequestID(ctx context.Context) string {
	return ctx.Value(RequestIDKey).(string)
}
