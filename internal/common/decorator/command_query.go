package decorator

import (
	"context"
)

func ApplyCommandQueryDecorators[H any, R any](handler CommandQueryHandler[H, R]) CommandQueryHandler[H, R] {
	return handler
}

type CommandQueryHandler[CQ any, R any] interface {
	Handle(ctx context.Context, cmdQr CQ) (R, error)
}
