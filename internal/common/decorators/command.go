package decorators

import (
	"context"
	"fmt"
	"strings"
	"test_ina_bank/pkg/baselogger"
	"test_ina_bank/pkg/telemetry"
)

func ApplyCommandDecorators[H any](handler CommandHandler[H], logger *baselogger.Logger, tracer *telemetry.OtelSdk) CommandHandler[H] {
	return commandLoggingDecorator[H]{
		base: commandMetricsDecorator[H]{
			base:   handler,
			client: tracer,
		},
		logger: logger,
	}
}

type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}
