# ADR-004: Repository Pattern

## Status
Aceito

## Contexto
A aplicação necessita persistir dados em MongoDB mas precisa:
- Testar business logic sem depender de banco de dados real
- Possibilidade de trocar tecnologia de persistência no futuro
- Isolar complexidade de MongoDB do restante da aplicação
- Suportar múltiplas implementações (in-memory para testes, MongoDB para produção)

## Decisão
Aplicar Repository Pattern com interfaces definindo contratos de persistência e múltiplas implementações concretas.

## Interface Contract

```go
// ItemRepository defines the contract for item persistence
type ItemRepository interface {
    Create(ctx context.Context, item Item) (Item, error)
    Update(ctx context.Context, item Item) (Item, error)
    Delete(ctx context.Context, id string) error
    GetByID(ctx context.Context, id string) (Item, error)
    List(ctx context.Context) ([]Item, error)
    BulkUpdateActive(ctx context.Context, active bool) (int64, int64, error)
}
```

## Implementações

### 1. Local Repository (`internal/repository/local/`)
**Uso:** Desenvolvimento, testes unitários, prototipagem rápida

**Características:**
- In-memory storage com `map[string]Item`
- Sem dependências externas
- Instant startup
- Ideal para testes isolados

```go
type LocalRepository struct {
    items map[string]Item
    mu    sync.RWMutex
}
```

### 2. MongoDB Repository (`internal/repository/mongodb/`)
**Uso:** Produção, ambientes de staging

**Características:**
- Persistência durável com MongoDB
- Suporte a transações para operações atômicas
- Índices configurados para performance
- Connection pooling gerenciado pelo driver

```go
type MongoRepository struct {
    client *mongo.Client
    db     string
}
```

## Rationale
**Inversão de Dependência:** Service layer depende da abstração (interface), não da implementação concreta. Isso permite substituir MongoDB por PostgreSQL sem mudar uma linha de service.

**Testabilidade:** Testes unitários de service podem usar mock repository ou local repository, eliminando necessidade de Docker containers ou dependências externas.

**Isolamento:** Detalhes específicos de MongoDB (BSON tags, `primitive.ObjectID`, `context.Context`) ficam encapsulados na implementação MongoDB. O restante da aplicação trabalha com domain entities puras.

**Flexibilidade:** Implementações adicionais podem ser adicionadas facilmente (Redis cache, Elasticsearch, etc.) sem impactar código existente.

**Performance de Testes:** Testes com in-memory repository executam em milissegundos vs segundos com MongoDB real.

## Consequências

### Positivas
- Testes rápidos e determinísticos
- Fácil migration entre tecnologias de persistência
- Separação clara de responsabilidades
- Possibilidade de estratégias de caching (ex: RedisRepository) sem mudar services
- Mock generation simplificada com interfaces pequenas

### Negativas
- Camada adicional de indireção
- Pode parecer over-engineering para operações CRUD simples
- Trade-off entre flexibilidade e simplicidade

## Exemplo de Uso

```go
// Service layer usando interface
type ItemService struct {
    repo repository.ItemRepository
}

func (s *ItemService) CreateItem(ctx context.Context, item domain.Item) (domain.Item, error) {
    // Business logic aqui
    return s.repo.Create(ctx, item)
}

// Production: MongoDB implementation
itemService := NewItemService(mongodb.NewMongoRepository(...))

// Testing: Local implementation
itemService := NewItemService(local.NewLocalRepository(...))
```

## Referências
- [Patterns of Enterprise Application Architecture - Repository](https://martinfowler.com/eaaCatalog/repository.html)
- [Microsoft Docs - Repository Pattern](https://docs.microsoft.com/en-us/azure/architecture/patterns/repository)
