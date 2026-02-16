# ADR-010: Middleware Order

## Status
Aceito

## Contexto
HTTP requests passam por múltiplos middlewares antes de reaching handlers. A ordem de execução é crítico para:
- CORS preflight requests funcionarem corretamente
- Requests being logged corretamente
- Error handling capturando errors de middlewares anteriores

Ordem incorreta pode resultar em:
- CORS headers não sendo enviados
- Logs não sendo registrados
- Errors não sendo capturados

## Decisão
Estabelecer ordem explícita de middlewares: **CORS → Logging → Router → Error Wrapping**

## Ordem de Execução

```
Incoming Request
       │
       ▼
┌─────────────────────┐
│  CORS Middleware    │ ← Primeiro: handle preflight, add headers
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│  Logging Middleware │ ← Segundo: log request details
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│      Router         │ ← Terceiro: route to correct handler
│   (gorilla/mux)     │
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│ Error Handling Wrap  │ ← Ao redor de cada handler: catch errors
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│     Handler         │ ← Execute business logic
└─────────────────────┘
```

## Implementação

```go
// cmd/api/server/server.go

// 1. Setup router com handlers
router := mux.NewRouter()
router.Handle("/healthz", middleware.ErrorHandlingMiddleware(s.healthHandler.HealthCheck)).Methods("GET")
router.Handle("/item", middleware.ErrorHandlingMiddleware(s.handler.CreateItem)).Methods("POST")
// ... outras rotas

// 2. Criar logging middleware
loggingMiddleware := middleware.LoggingMiddleware(s.logger)

// 3. Criar CORS middleware
corsMiddleware := middleware.CORSMiddleware([]string{"*"})

// 4. Aplicar middlewares (ordem é importante!)
//    Handler wrapped por ErrorHandling → wrapped por Logging → wrapped por CORS
s.server.Handler = corsMiddleware(loggingMiddleware(router))
```

## Rationale por Ordem

### 1. CORS First
**Por que:** CORS preflight requests (OPTIONS) não devem passar por logging ou business logic.

**Se depois de logging:** Preflight requests seriam logados desnecessariamente.

**Se depois de router:** Router pode não encontrar route para OPTIONS, retornando 404.

```go
func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            // Handle preflight
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

### 2. Logging Second
**Por que:** Queremos logar requests que passaram CORS check. Não logamos preflight requests.

**Middleware:**
```go
func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            logger.Info("incoming request",
                zap.String("method", r.Method),
                zap.String("path", r.URL.Path),
                zap.String("remote_addr", r.RemoteAddr))
            next.ServeHTTP(w, r)
            logger.Info("request completed",
                zap.Duration("duration", time.Since(start)))
        })
    }
}
```

### 3. Router Third
**Por que:** Routing é o core de dispatch para handlers apropriados.

**Se antes de CORS/logging:** Preflight requests seriam roteadas incorretamente.

### 4. Error Handling Wraps Each Handler
**Por que:** Errors podem ocorrer em qualquer handler. Wrapping individual handlers permite error handling granular.

```go
func ErrorHandlingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                handlers.SendJSONError(w, "Internal server error", http.StatusInternalServerError)
            }
        }()
        next(w, r)
    }
}
```

## Response Flow (Reverse Order)

```
Response
   │
   ▼
Handler returns response
   │
   ▼
Error Handling Middleware (check for errors)
   │
   ▼
Logging Middleware (log response status/duration)
   │
   ▼
CORS Middleware (add CORS headers to response)
   │
   ▼
HTTP Response sent to client
```

## Exemplo Completo de Request/Response

```
Request: POST /item
├─ CORS Middleware: ✅ Add headers, not OPTIONS, continue
├─ Logging Middleware: 📝 Log "POST /item"
├─ Router: 📍 Route to CreateItem handler
├─ Error Wrapper: ✅ No panic, execute handler
└─ CreateItem Handler: 💼 Create item, return 201

Response (reverse order):
├─ CreateItem Handler: ← Return JSON response
├─ Error Wrapper: ← No error, pass through
├─ Logging Middleware: 📝 Log "201 Created in 50ms"
└─ CORS Middleware: ✅ Add CORS headers to response
   └─ Client receives: 201 Created with CORS headers
```

## Consequências

### Positivas
- CORS funciona corretamente para preflight e regular requests
- Logs limpos (sem preflight noise)
- Errors são capturados e retornados como JSON
- Ordem explícita facilita debug

### Negativas
- Ordem deve ser mantida manualmente (não há validação automática)
- Adicionar novo middleware requer entender onde posicionar

## Testando Middleware Order

```go
func TestMiddlewareOrder(t *testing.T) {
    // Preflight request deve parar no CORS middleware
    req := httptest.NewRequest("OPTIONS", "/item", nil)
    rec := httptest.NewRecorder()

    router.ServeHTTP(rec, req)

    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
}
```

## Referências
- [Go net/http Middleware Pattern](https://www.alexedwards.net/blog/middleware-and-chaining)
- [CORS MDN Documentation](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
