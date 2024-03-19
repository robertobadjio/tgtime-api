package decorator

import (
	"context"
)

func ApplyCommandQueryDecorators[H any](handler CommandQueryHandler[H]) CommandQueryHandler[H] {
	return handler
}

type CommandQueryHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}
