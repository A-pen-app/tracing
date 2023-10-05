package tracing

import (
	"context"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

var tp *sdktrace.TracerProvider
var projectID string
var serviceName string
var deploymentEnvironment string
var tracerName string

type Config struct {
	ProjectID             string `json:"project_id" yaml:"project_id"`
	TracerName            string `json:"tracer_name" yaml:"tracer_name"`
	ServiceName           string `json:"service_name" yaml:"service_name"`
	DeploymentEnvironment string `json:"deployment_environment" yaml:"deployment_environment"`
}

func Start(ctx context.Context, name string) trace.Span {
	_, span := tp.Tracer(tracerName).Start(ctx, name)
	return span
}

func Initialize(ctx context.Context, c *Config) error {
	if c != nil {
		projectID = c.ProjectID
		serviceName = c.ServiceName
		deploymentEnvironment = c.DeploymentEnvironment
		tracerName = c.TracerName
	}

	exporter, err := texporter.New(texporter.WithProjectID(projectID))
	if err != nil {
		return err
	}
	res, err := resource.New(ctx,
		resource.WithDetectors(gcp.NewDetector()),
		//resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.DeploymentEnvironmentKey.String(deploymentEnvironment),
		),
	)
	if err != nil {
		return err
	}

	tp = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return nil
}

func Finalize(ctx context.Context) {
	if tp != nil {
		tp.Shutdown(ctx)
	}
}
