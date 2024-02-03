package decorators

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"strings"
	"test_ina_bank/pkg/telemetry"
	"time"
)

type commandMetricsDecorator[C any] struct {
	base   CommandHandler[C]
	client *telemetry.OtelSdk
}

func (d commandMetricsDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	start := time.Now()
	actionName := strings.ToLower(generateActionName(cmd))
	counter, _ := d.client.Metric.Int64Counter(fmt.Sprintf("commands.%s.duration", actionName),
		metric.WithUnit("ms"),
		metric.WithDescription("Function latency in ms"))

	successCounter, _ := d.client.Metric.Int64Counter(fmt.Sprintf("commands.%s.success", actionName),
		metric.WithUnit("1"),
		metric.WithDescription("Total command success"))

	failedCounter, _ := d.client.Metric.Int64Counter(fmt.Sprintf("commands.%s.failure", actionName),
		metric.WithUnit("1"),
		metric.WithDescription("Total command failure"))

	ctx, span := d.client.Tracer.Start(ctx, fmt.Sprintf("usecase.%s", actionName))
	defer span.End()

	defer func() {
		end := time.Since(start)

		counter.Add(ctx, end.Microseconds())
		if err == nil {
			successCounter.Add(ctx, 1)
		} else {
			failedCounter.Add(ctx, 1)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	return d.base.Handle(ctx, cmd)
}

type queryMetricsDecorator[C any, R any] struct {
	base   QueryHandler[C, R]
	client *telemetry.OtelSdk
}

func (d queryMetricsDecorator[C, R]) Handle(ctx context.Context, query C) (result R, err error) {
	start := time.Now()
	actionName := strings.ToLower(generateActionName(query))
	counter, _ := d.client.Metric.Int64Counter(fmt.Sprintf("queries.%s.duration", actionName),
		metric.WithUnit("ms"),
		metric.WithDescription("Function latency in ms"))

	successCounter, _ := d.client.Metric.Int64Counter(fmt.Sprintf("queries.%s.success", actionName),
		metric.WithUnit("1"),
		metric.WithDescription("Total query success"))

	failedCounter, _ := d.client.Metric.Int64Counter(fmt.Sprintf("queries.%s.failure", actionName),
		metric.WithUnit("1"),
		metric.WithDescription("Total query failure"))

	ctx, span := d.client.Tracer.Start(ctx, fmt.Sprintf("usecase.%s", actionName))
	defer span.End()

	defer func() {
		end := time.Since(start)

		counter.Add(ctx, end.Microseconds())
		if err == nil {
			successCounter.Add(ctx, 1)
		} else {
			failedCounter.Add(ctx, 1)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	return d.base.Handle(ctx, query)
}
