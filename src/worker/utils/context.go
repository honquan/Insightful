package utils

import (
	"context"
)

const ArgsCtxKey = "args"

// get set args
func SetArgs(ctx context.Context, args []string) context.Context {
	return context.WithValue(ctx, ArgsCtxKey, args)
}
func GetArgs(ctx context.Context) []string {
	uid, ok := ctx.Value(ArgsCtxKey).([]string)
	if !ok {
		return nil
	}
	return uid
}
