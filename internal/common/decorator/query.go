package decorator

import "context"

func ApplyQueryDecorators[H any, R any](handler QueryHandler[H, R]) QueryHandler[H, R] {
	return handler
}

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, qr Q) (R, error)
}
