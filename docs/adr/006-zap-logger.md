# ADR-006: Zap Logger

## Status
Aceito

## Contexto
Aplicações production-ready necessitam logging estruturado para:
- Debug eficiente em produção
- Análise de comportamento da aplicação
- Monitoramento de erros e performance
- Facilitar troubleshooting de incidentes

Logging baseado em `fmt.Printf` ou `log.Printf` não fornece:
- Estrutura semântica (nível, campos, contexto)
- Performance adequada para alta throughput
- Integração com sistemas de log agregados (ELK, Loki, CloudWatch)

## Decisão
Utilizar Uber Zap (`go.uber.org/zap`) como solução de logging estruturado.

## Rationale

### Comparativo de Loggers

| Feature | Zap | Logrus | Zerolog | stdlib log |
|---------|-----|--------|---------|-------------|
| Structured logging | ✅ | ✅ | ✅ | ❌ |
| Performance | 🟢 Mais rápido | 🟡 Lento | 🟢 Rápido | 🟢 Rápido (mas básico) |
| Zero-allocation | ✅ | ❌ | ✅ | ❌ |
| Field-based API | ✅ | ✅ | ✅ | ❌ |
| Leveled logging | ✅ | ✅ | ✅ | Parcial |
| Development activity | 🟢 Manutenção | 🟡 Baixa | 🟢 Ativo | 🟢 Stdlib |

### Benchmarks (Operações/segundo, maior é melhor)
```
Zap:        86,912 ops/sec
Zerolog:    80,004 ops/sec
Stdlib:     54,875 ops/sec
Logrus:     12,279 ops/sec
```

*Fonte: [Zap Documentation](https://github.com/uber-go/zap#performance)*

### Features do Zap

**Structured Fields:** Logs como dados estruturados, não strings
```go
logger.Info("item created",
    zap.String("id", item.ID),
    zap.String("name", item.Name),
    zap.Bool("active", item.Active))
```

**Zero-Allocation:** Minimiza pressão no garbage collector

**Configuração Flexível:** Development (console) vs Production (JSON)

**Leveled Logging:** Debug, Info, Warn, Error, Fatal, Panic

## Consequências

### Positivas
- Logs query-friendly em sistemas agregados (ex: Loki, CloudWatch)
- Performance excelente para high-throughput APIs
- Fields tipados reduzem erros de parsing
- Output JSON em produção para integração com observability stack

### Negativas
- Verbosidade maior que `fmt.Printf` para logs simples
- Mais uma dependência externa
- Syntax field-based pode parecer verbose inicialmente

## Implementação

### Development Logger (console)
```go
logger, _ := zap.NewDevelopment()
defer logger.Sync()
```

### Production Logger (JSON)
```go
logger, _ := zap.NewProduction()
defer logger.Sync()
```

### Uso na Aplicação

**Middleware de Logging:**
```go
func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            logger.Info("incoming request",
                zap.String("method", r.Method),
                zap.String("path", r.URL.Path))
            // ...
        })
    }
}
```

**Error Logging:**
```go
if err != nil {
    logger.Error("failed to create item",
        zap.Error(err),
        zap.String("name", item.Name))
    return err
}
```

## Integração Futura
Zap suporta OpenTelemetry via `zapcore`. Isso facilitará integração com distributed tracing planejada no projeto.

## Alternativas Consideradas

**Logrus:** Mais popular mas significativamente mais lento. Adequado para aplicações com menor throughput.

**Zerolog:** Performance similar a Zap com API diferente. Escolha foi baseada em maturidade e adoção em projetos enterprise da Uber.

**stdlib log:** Insuficiente para logging estruturado production-ready.

## Referências
- [Zap GitHub](https://github.com/uber-go/zap)
- [Zap Performance Benchmarks](https://github.com/uber-go/zap#performance)
