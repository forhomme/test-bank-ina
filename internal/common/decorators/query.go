package decorators

import (
	"context"
	"test_ina_bank/pkg/baselogger"
	"test_ina_bank/pkg/telemetry"
)

func ApplyQueryDecorators[H any, R any](handler QueryHandler[H, R], logger *baselogger.Logger, tracer *telemetry.OtelSdk) QueryHandler[H, R] {
	return queryLoggingDecorator[H, R]{
		base: queryMetricsDecorator[H, R]{
			base:   handler,
			client: tracer,
		},
		logger: logger,
	}
}

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}
