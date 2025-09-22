package dd

import (
	"context"

	"github.com/gin-gonic/gin"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func Load(dd_service string, dd_env string, dd_version string) {
	// Implementação fictícia para iniciar o tracer do Datadog
	tracer.Start(
		tracer.WithServiceName(dd_service),
		tracer.WithEnv(dd_env),
		tracer.WithServiceVersion(dd_version),
	)
}

func Stop() {
	tracer.Stop()
}

func StartSpan(ctx context.Context, name string, opts ...tracer.StartSpanOption) (tracer.Span, context.Context) {
	return tracer.StartSpanFromContext(ctx, name, opts...)
}

func FinishSpan(span tracer.Span) {
	if span != nil {
		span.Finish()
	}
}

func SetSpanError(span tracer.Span, err error) {
	if span != nil && err != nil {
		span.SetTag("error", err)
	}
}

func SetSpanTag(span tracer.Span, key string, value interface{}) {
	if span != nil {
		span.SetTag(key, value)
	}
}

func GinMiddleware(service string) gin.HandlerFunc {
	return gintrace.Middleware(service)
}
