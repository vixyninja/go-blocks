package logx

import "context"

type Logx interface {
	Debug(c context.Context, msg string, args ...any)
	Info(c context.Context, msg string, args ...any)
	Warn(c context.Context, msg string, args ...any)
	Error(c context.Context, msg string, args ...any)
	Fatal(c context.Context, msg string, args ...any)

	// REMAINING: This is for cloning the logger and adding fields to it
	With(c context.Context, fields map[string]any) Logx
}
