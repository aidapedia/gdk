package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	env "github.com/aidapedia/gdk/environment"
)

// Span is a struct to handle span
type Span struct {
	span trace.Span
}

func InitTracer(serviceName string, collectorURL string, isSecured bool) (*otlptrace.Exporter, error) {
	secureOption := otlptracegrpc.WithInsecure()

	isSecured = false
	if isSecured { // TODO add tls and fix it later
		secureOption = otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)

	if err != nil {
		return nil, err
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service", serviceName),
			attribute.String("library.language", "go"),
			attribute.String("library.language", "go"),
			attribute.String("environment", env.Development),
		),
	)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)
	return exporter, nil
}

// StartSpanFromContext create and start an opentracing and newrelic span from context
func StartSpanFromContext(ctx context.Context, operation string) (Span, context.Context) {
	ctx, sp := otel.Tracer("").Start(ctx, operation)
	return Span{
		span: sp,
	}, ctx
}

// Finish finish the span
func (s *Span) Finish(errors error) {
	var message string
	err := errors
	status := codes.Ok
	if err != nil {
		s.span.RecordError(err, trace.WithStackTrace(true))
		message = err.Error()
		status = codes.Error
	}
	s.span.SetStatus(status, message)
	s.span.End()
}
