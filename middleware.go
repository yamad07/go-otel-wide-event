package main

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

var mainSpanContextKey = "mainSpan"

func mainSpanMiddleware(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			url := r.URL.String()
			tracer := otel.Tracer(serviceName)
			ctx, span := tracer.Start(ctx, url)

			ctx = context.WithValue(ctx, mainSpanContextKey, span)
			span.SetAttributes(
				semconv.URLPath(url),
				semconv.URLScheme(getScheme(r)),
				semconv.HTTPRequestMethodKey.String(r.Method),
			)

			defer span.End()

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)

		})
	}
}

func getMainSpan(
	r *http.Request,
) trace.Span {
	ctx := r.Context()
	span, ok := ctx.Value(mainSpanContextKey).(trace.Span)
	if !ok {
		return nil
	}
	return span
}

func getScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	return "http"
}
