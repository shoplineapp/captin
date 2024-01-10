package helpers

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func Tracer() trace.Tracer {
	return otel.GetTracerProvider().Tracer("captin")
}
