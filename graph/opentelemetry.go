package graph

import (
	"context"

	"go.opentelemetry.io/contrib/detectors/aws/ecs"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
)

var Tracer = otel.Tracer("github.com/ryohei1216/graphql-trace")

func New(ctx context.Context) (func(context.Context) error, error) {
	// OTLP Exporterを作成する
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	// ECSのデータを検出する
	ecsResourceDetector := ecs.NewResourceDetector()

	resource, err := resource.New(
		ctx,
		resource.WithDetectors(ecsResourceDetector),
	)
	if err != nil {
		return nil, err
	}

	tp := sdkTrace.NewTracerProvider(
		sdkTrace.WithSampler(sdkTrace.AlwaysSample()),
		sdkTrace.WithBatcher(traceExporter),
		sdkTrace.WithIDGenerator(xray.NewIDGenerator()),
		sdkTrace.WithResource(resource),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})

	return func(context.Context) error {
		err = tp.Shutdown(ctx)
		if err != nil {
			return err
		}
		return nil
	}, nil
}
