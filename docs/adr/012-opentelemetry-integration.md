# ADR-012: OpenTelemetry Integration

## Status
Proposto / Planejado

## Contexto
Aplicações modernas production-ready necessitam observabilidade:
- **Tracing:** Requests através de múltiplos serviços (handlers → services → repositories)
- **Metrics:** Contadores, gauges, histograms para monitoramento (latency, throughput, errors)
- **Logging:** Estruturado já implementado com Zap

Soluções atuais (Zap logging) não fornecem:
- Distributed tracing across services
- Automatic metrics collection (latency percentiles, error rates)
- Integration com observability platforms (Datadog, New Relic, Grafana)

## Decisão Planejada
Integrar OpenTelemetry (OTel) para tracing e metrics, mantendo Zap para logging.

## Arquitetura Proposta

```
┌─────────────────────────────────────────────────┐
│              Application Code                    │
│  (handlers, services, repositories)           │
└──────────────────┬──────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────┐
│            OpenTelemetry API                    │
│  - Tracer for creating spans                   │
│  - Meter for creating metrics                   │
└──────────────────┬──────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────┐
│           OTel SDK (Go)                         │
│  - Span processors                             │
│  - Batch exporters                            │
└──────────────────┬──────────────────────────────┘
                   │
        ┌──────────┴──────────┐
        ▼                     ▼
┌──────────────┐      ┌──────────────┐
│   Traces     │      │   Metrics    │
│  Exporter    │      │  Exporter    │
│  (OTLP/HTTP) │      │  (OTLP/HTTP) │
└──────┬───────┘      └──────┬───────┘
       │                     │
       └──────────┬──────────┘
                  ▼
        ┌────────────────────┐
        │  Observability    │
        │  Platform        │
        │  (Grafana Cloud, │
        │   Datadog, etc)  │
        └──────────────────┘
```

## Componentes

### 1. Tracing

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

// In handler
func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
    ctx, span := otel.Tracer("handlers").Start(r.Context(), "CreateItem")
    defer span.End()

    // Pass context to service
    item, err := h.service.CreateItem(ctx, item)
}

// In service
func (s *ItemService) CreateItem(ctx context.Context, item Item) (Item, error) {
    ctx, span := otel.Tracer("service").Start(ctx, "CreateItem")
    defer span.End()

    // Pass context to repository
    return s.repo.Create(ctx, item)
}
```

**Traces gerados:**
```
POST /item
├─ handler.CreateItem (2ms)
├─ service.CreateItem (1ms)
└─ repository.Create (0.5ms)
```

### 2. Metrics

```go
import (
    "go.opentelemetry.io/otel/metric"
)

var (
    requestCounter  metric.Int64Counter
    requestLatency metric.Float64Histogram
)

func init() {
    meter := otel.Meter("list-manager-api")
    requestCounter, _ = meter.Int64Counter("http_requests_total")
    requestLatency, _ = meter.Float64Histogram("http_request_duration_ms")
}

func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    defer func() {
        duration := float64(time.Since(start).Milliseconds())
        requestLatency.Record(r.Context(), duration)
        requestCounter.Add(r.Context(), 1)
    }()
}
```

**Metrics coletadas:**
- `http_requests_total{method="POST", path="/item", status="201"}` ← Counter
- `http_request_duration_ms{method="POST", path="/item"}` ← Histogram (p50, p95, p99)
- `active_connections` ← Gauge

### 3. Logging (Zap + OTel Bridge)

```go
import "go.opentelemetry.io/contrib/bridges/otelzap"

func main() {
    logger := zap.NewProduction()
    logger = otelzap.NewLogger(logger, otelzap.WithTraceIDField())
}
```

**Logs enriquecidos com trace context:**
```json
{
  "level": "info",
  "msg": "item created",
  "trace_id": "7f8a9d...",
  "span_id": "3b2e1f...",
  "item_id": "123"
}
```

## Exporters

### OTLP over HTTP (Recommended)
```go
import (
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

exporter, _ := otlptracehttp.New(ctx,
    otlptracehttp.WithEndpoint("otlp.nr-data.net:4318"), // New Relic
    otlptracehttp.WithHeaders(map[string]string{"api-key": "..."}),
)
```

### Console Exporter (Development)
```go
exporter, _ := stdout.New(stdout.WithPrettyPrint())
```

## Instrumentation Manual vs Automatic

**Manual (Escolhido):**
- Span creation explícito no código
- Maior controle sobre o que é traced
- Baixo overhead
- Trade-off: mais código boilerplate

**Automatic (Consideração futura):**
- `http.Handler` wrapper automático
- Zero code changes para HTTP layer
- Trade-off: menos granular, mais overhead

## Roadmap

### Phase 1: Setup (Sprint X)
- [ ] Adicionar dependencies OTel
- [ ] Configurar tracer provider
- [ ] Instrumentar handlers manualmente
- [ ] Console exporter para dev

### Phase 2: Service Layer (Sprint Y)
- [ ] Instrumentar service methods
- [ ] Instrumentar repository methods
- [ ] Validar trace propagation

### Phase 3: Metrics (Sprint Z)
- [ ] Configurar meter provider
- [ ] Adicionar HTTP metrics (latency, errors, throughput)
- [ ] Adicionar business metrics (items created, updated, deleted)

### Phase 4: Production (Sprint W)
- [ ] Configurar OTLP exporter para Grafana Cloud
- [ ] Dashboard de observability
- [ ] Alerts (p95 latency threshold, error rate SLO)

## Dependencies

```go
// go.mod
require (
    go.opentelemetry.io/otel v1.x.x
    go.opentelemetry.io/otel/trace v1.x.x
    go.opentelemetry.io/otel/metric v1.x.x
    go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.x.x
    go.opentelemetry.io/contrib/bridges/otelzap v1.x.x
)
```

## Alternativas Consideradas

### Cloud-native solutions (Datadog Agent, New Relic Agent)
✅ Menos configuração
❌ Vendor lock-in
❌ Agents em cada host/pod

### Prometheus direct (sem OTel)
✅ Simples, bem estabelecido
❌ Sem distributed tracing (sem bridge com Jaeger)
❌ Exporters específicos por vendor

### OTel: Universal, vendor-agnostic
✅ Standard em aberto
✅ Switch observability platforms sem mudar código
✅ OpenTelemetry tem momentum como padrão

## Consequências

### Positivas
- Distributed tracing para debug de requests through services
- Metrics automático para SLO monitoring
- Vendor lock-in eliminado
- Logs, traces, metrics correlacionados via trace ID

### Negativas
- Curva de aprendizado para OTel concepts
- Slight performance overhead (~5%)
- Mais boilerplate no código (span creation)

## Custo

**Grafana Cloud Free Tier:**
- 50 GB logs retention
- 10 GB metrics
- 50 GB traces

**Escalando:**
- $49/month para 100GB logs

## References
- [OpenTelemetry Go Documentation](https://opentelemetry.io/docs/instrumentation/go/)
- [OpenTelemetry Specification](https://opentelemetry.io/docs/reference/specification/)
- [Semantic Conventions](https://opentelemetry.io/docs/reference/specification/trace/semantic_conventions/)
