package xhttp

import (
	"github.com/labstack/echo"
	"github.com/uber-go/zap"
	"golang.org/x/net/context"
)

var reqRootLogger = zap.NewJSON()

const RequestLoggerKey = "_reqLogger_"

func RequestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := TakeRequestID(c)
		reqLogger := reqRootLogger.With(zap.String("requestID", requestID))
		c.SetContext(context.WithValue(c.Context(), RequestLoggerKey, reqLogger))
		return next(c)
	}
}

func CurrentLogger(c context.Context) zap.Logger {
	return c.Value(RequestLoggerKey).(zap.Logger)
}
