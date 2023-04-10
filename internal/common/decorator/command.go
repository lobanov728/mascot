package decorator

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

func ApplyCommandDecorators[H any, R any](handler CommandHandler[H, R], logger *logrus.Entry) CommandHandler[H, R] {
	return commandLoggingDecorator[H, R]{
		base:   handler,
		logger: logger,
	}
}

type CommandHandler[C any, R any] interface {
	Handle(ctx context.Context, cmd C) (R, error)
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}
