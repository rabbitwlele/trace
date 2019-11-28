package trace

import "context"

type Logger interface {
	Print(ctx context.Context, v ...interface{})
}
