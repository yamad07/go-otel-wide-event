# go-otel-wide-event

A Go sample project demonstrating the **Wide Event pattern** with OpenTelemetry.
Aggregates child span attributes into the parent span for more efficient trace analysis.

## Quick Start

### Prerequisites

- Go 1.21+
- Docker (for Jaeger)

### 1. Start Jaeger

```bash
docker run --rm \
  -p 16686:16686 \
  -p 4318:4318 \
  jaegertracing/jaeger:2.13.0
```

### 2. Start the server

```bash
go run .
```

### 3. Test it

```bash
curl http://localhost:3000/articles
```

→ View traces at [http://localhost:16686](http://localhost:16686)

## API

| Method | Endpoint    | Description       |
|--------|-------------|-------------------|
| GET    | `/`         | Welcome message   |
| GET    | `/articles` | Get article data  |

## Wide Event Pattern

### Trace Structure

```
Main Span: "/articles"
├── url.path: "/articles"
├── url.scheme: "http"
├── http.method: "GET"
└── article.id: "1"        ← Propagated from child span
    │
    └── Child Span: "getArticle"
        ├── article.id: "1"
        └── article.title: "はじめてのGo"
```

### Implementation

**1. Store parent span in Context via middleware**

```go
// middleware.go
ctx, span := tracer.Start(ctx, url)
ctx = context.WithValue(ctx, mainSpanContextKey, span)
```

**2. Retrieve parent span in handler and add attributes**

```go
// main.go
mainSpan := getMainSpan(r)
mainSpan.SetAttributes(
    attribute.String("article.id", article.ID),
)
```

## Configuration

Change the Jaeger endpoint in `main.go`:

```go
otlptracehttp.WithEndpoint("localhost:4318"),
```

## References

- [OpenTelemetry Go SDK](https://opentelemetry.io/docs/languages/go/)
- [Wide Events - Honeycomb Blog](https://www.honeycomb.io/blog/wide-events-and-the-future-of-observability)
