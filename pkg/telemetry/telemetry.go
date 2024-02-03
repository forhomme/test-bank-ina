package telemetry

import (
	"context"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	tracer "go.opentelemetry.io/otel/trace"
	"test_ina_bank/config"
)

type OtelSdk struct {
	Tracer   tracer.Tracer
	Metric   metric.Meter
	Shutdown func(ctx context.Context) error
}

func NewOtelSdk(cfg *config.Cfg) *OtelSdk {
	// Configure OpenTelemetry with sensible defaults.
	uptrace.ConfigureOpentelemetry(
		// copy your project DSN here or use UPTRACE_DSN env var
		uptrace.WithDSN(cfg.General.UptraceDSN),

		uptrace.WithServiceName(cfg.General.AppName),
		uptrace.WithServiceVersion(cfg.General.AppVersion),
	)

	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errs from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = fn(ctx)
		}
		shutdownFuncs = nil
		return err
	}

	shutdownFuncs = append(shutdownFuncs, uptrace.Shutdown)

	return &OtelSdk{
		Tracer:   otel.Tracer(cfg.General.AppName),
		Metric:   otel.Meter(cfg.General.AppName),
		Shutdown: shutdown,
	}
}
