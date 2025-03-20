package opentelemetry

import (
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer trace.Tracer
	once   sync.Once
)

func GetTracer() trace.Tracer {
	once.Do(func() {
		tracer = otel.Tracer("repos_db")
	})
	return tracer
}
