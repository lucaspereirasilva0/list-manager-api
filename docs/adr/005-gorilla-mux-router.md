# ADR-005: Gorilla/Mux como Router HTTP

## Status
Aceito

## Contexto
A aplicação necessita de um router HTTP que ofereça:
- Roteamento baseado em métodos (GET, POST, PUT, DELETE)
- Suporte a path variables e query parameters
- Middleware chain flexível
- Performance adequada para throughput esperado
- Estabilidade e maturidade no ecossistema Go

## Decisão
Utilizar Gorilla Mux (`github.com/gorilla/mux`) como router HTTP principal da aplicação.

## Rationale

### Features Específicas para API REST
**Method-based Routing:** Restrição explícita de métodos HTTP por rota:
```go
router.Handle("/item", handler.CreateItem).Methods("POST")
router.Handle("/item", handler.GetItem).Methods("GET")
router.Handle("/item", handler.UpdateItem).Methods("PUT")
router.Handle("/item", handler.DeleteItem).Methods("DELETE")
```

**Path Variables:** Suporte a variáveis de_path (útil para futuras expansões):
```go
router.Handle("/items/{id}", handler.GetItem).Methods("GET")
```

**Query Parameters:** Parsing automático via `http.Request` para endpoints como `/item?id=123`.

### Matriz Comparativa de Routers

| Feature | Gorilla/Mux | Chi Router | Gin | net/http stdlib |
|---------|-------------|------------|-----|-----------------|
| Method restriction | ✅ | ✅ | ✅ | Manual |
| Path variables | ✅ | ✅ | ✅ | Manual |
| Middleware support | ✅ | ✅ | ✅ | Manual |
| Performance | Alto | Mais alto | Mais alto | N/A (baseline) |
| Maturity | 🟢 Estável | 🟡 Relativamente novo | 🟢 Estável | 🟢 Stdlib |
| Learning curve | Baixa | Baixa | Média | Baixa (manual) |
| Dependencies | Mínimas | Mínimas | Mínimas | Zero |

### Por que não usar Gin ou Chi?

**Gin:** Framework completo com template engine, binding, etc. Maior que necessário para API simples.

**Chi:** Excelente performance mas ainda relativamente novo (2020+). Mux é estável desde 2012.

**stdlib net/http:** Roteamento manual seria necessário ou uso de `http.ServeMux` que não suporta method restriction nativamente antes de Go 1.22.

### Matriz de Decisão

| Critério | Peso | Gorilla/Mux | Chi | Gin | stdlib |
|-----------|-------|-------------|-----|-----|---------|
| Maturidade | 5 | 5 | 3 | 5 | 5 |
| Performance | 4 | 4 | 5 | 5 | 5 |
| Features REST | 5 | 5 | 5 | 5 | 1 |
| Simplicidade | 4 | 4 | 5 | 3 | 5 |
| **Score** | - | **90** | **85** | **82** | **70** |

## Consequências

### Positivas
- Syntax limpa e intuitiva para definição de rotas
- Router altamente testado e estável (11+ anos)
- Middleware chain flexível (CORS → Logging → Router)
- Zero configuração adicional necessária
- Upgrade path para Go 1.22+ `http.ServeMux` seria simples

### Negativas
- Router framework adicional (aumentam superfície de dependências)
- Performance ligeiramente inferior a Chi/Gin (mas não é bottleneck para API CRUD)

## Implementação

```go
// cmd/api/server/server.go
router := mux.NewRouter()

// Health check
router.Handle("/healthz", middleware.ErrorHandlingMiddleware(s.healthHandler.HealthCheck)).Methods("GET")

// Item routes
router.Handle("/item", middleware.ErrorHandlingMiddleware(s.handler.CreateItem)).Methods("POST")
router.Handle("/item", middleware.ErrorHandlingMiddleware(s.handler.GetItem)).Methods("GET")
router.Handle("/item", middleware.ErrorHandlingMiddleware(s.handler.UpdateItem)).Methods("PUT")
router.Handle("/item", middleware.ErrorHandlingMiddleware(s.handler.DeleteItem)).Methods("DELETE")
router.Handle("/items", middleware.ErrorHandlingMiddleware(s.handler.ListItems)).Methods("GET")
router.Handle("/items/active", middleware.ErrorHandlingMiddleware(s.handler.BulkUpdateActive)).Methods("PUT")

// PWA version endpoint (no middleware)
router.HandleFunc("/_app/version.json", handlers.GetVersion).Methods("GET")
```

## Alternativa Futura: Go 1.22+ net/http
Go 1.22 introduz melhorias em `http.ServeMux` incluindo method matching e wildcards. Para projetos novos hoje, stdlib seria preferível. Para este projeto, Mux continua adequado e não há urgência para migration.

## Referências
- [Gorilla Mux Documentation](https://github.com/gorilla/mux)
- [Go 1.22 Release Notes - Enhanced net/http](https://go.dev/doc/go1.22#net/http)
