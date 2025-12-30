package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const articleServiceName = "article-service"

func main() {
	exporter, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	tp, err := initTracer(exporter)
	if err != nil {
		log.Fatal(err)
	}
	defer tp.Shutdown(context.Background())

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(mainSpanMiddleware(articleServiceName))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Get("/articles", getArticle)

	log.Println("Server starting on :3000")
	http.ListenAndServe(":3000", r)
}

type Article struct {
	ID      string
	Title   string
	Content string
}

func getArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	mainSpan := getMainSpan(r)

	tracer := otel.Tracer(articleServiceName)
	_, span := tracer.Start(ctx, "getArticle")

	defer span.End()

	article := Article{
		ID:      "1",
		Title:   "はじめてのGo",
		Content: "Goは素晴らしいプログラミング言語です。",
	}

	span.SetAttributes(
		attribute.String("article.id", article.ID),
		attribute.String("article.title", article.Title),
	)

	// mainSpanには記事IDを属性として追加
	mainSpan.SetAttributes(
		attribute.String("article.id", article.ID),
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

func initTracer(exporter sdktrace.SpanExporter) (*sdktrace.TracerProvider, error) {
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("article-service"),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}
