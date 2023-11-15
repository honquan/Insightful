package log

import (
	"context"
	"go.uber.org/zap"
)

const RequestIDKey = "_request_id"
const MessageIDKey = "_message_id"

var logger *zap.Logger

func Logger() *zap.Logger {
	return logger
}

func init() {
	env := "DEV"
	ops := []zap.Option{
		zap.AddCallerSkip(2),
	}
	switch env {
	case "DEV":
		logger, _ = zap.NewDevelopment(ops...)
	default:
		logger, _ = zap.NewProduction(ops...)
	}
	zap.ReplaceGlobals(logger)
}

func getRequestID(ctx context.Context) string {
	id, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		return ""
	}
	return id
}

func getMessageID(ctx context.Context) string {
	id, ok := ctx.Value(MessageIDKey).(string)
	if !ok {
		return ""
	}
	return id
}

func expand(ctx context.Context, data ...interface{}) []interface{} {
	if ctx == nil {
		return data
	}
	// append requestID
	requestID := getRequestID(ctx)
	if requestID != "" {
		data = append(data, "r_id", requestID)
	}
	// append messageID
	msgID := getMessageID(ctx)
	if msgID != "" {
		data = append(data, "m_id", msgID)
	}
	return data
}

func callLogFn(logFn func(msg string, kv ...interface{}), msg string, data ...interface{}) {
	if len(data) == 0 {
		logFn(msg)
		return
	}
	logFn(msg, data...)
}

func Infow(ctx context.Context, msg string, data ...interface{}) {
	callLogFn(zap.S().Infow, msg, expand(ctx, data...)...)
}

func Warnw(ctx context.Context, msg string, data ...interface{}) {
	callLogFn(zap.S().Warnw, msg, expand(ctx, data...)...)
}

func Errorw(ctx context.Context, msg string, data ...interface{}) {
	callLogFn(zap.S().Errorw, msg, expand(ctx, data...)...)
}

func Debugw(ctx context.Context, msg string, data ...interface{}) {
	callLogFn(zap.S().Debugw, msg, expand(ctx, data...)...)
}
