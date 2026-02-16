# ADR-007: Mock-based Testing

## Status
Aceito

## Contexto
Testes unitários production-ready requerem:
- Isolamento completo de dependências externas
- Execução rápida (milissegundos, não segundos)
- Determinismo (mesmos inputs → mesmos outputs)
- Ausência de requisitos ambientais (Docker, network, databases)

Testes que dependem de banco de dados real (mesmo em containers) violam esses princípios:
- Setup/teardown lento
- Flakiness potential (network, race conditions)
- Difícil simular cenários de erro (timeout, connection refused)

## Decisão
Adotar testes baseados em mocks utilizando `stretchr/testify/mock` para isolar completamente unidades de código sob teste.

## Estratégia de Testes

### Pirâmide de Testes

```
        /\
       /  \        E2E Tests (5%)
      /____\       - Reais dependencies
     /      \      - Lentos, caros
    /        \     - Poucos cenários
   /  Unit    \
  /  Tests    \    Unit Tests (80%)
 /____________\   - Rápidos, isolados
                  - Mocks
                  - Cobrem todos os casos
```

### Camada de Testes

**Unit Tests (80%)**
- Repository layer: Mock MongoDB client operations
- Service layer: Mock repository interfaces
- Handlers layer: Mock service interfaces
- **Sem Docker, sem network, sem latência**

**Integration Tests (15%)**
- Repository implementation com MongoDB container
- Testes de e2e para fluxos críticos
- Executados apenas em CI/PR, não em cada run local

**Manual/E2E (5%)**
- Smoke tests manuais antes de deploy
- Testes exploratórios de UX

## Implementação com testify/mock

### Repository Mock

```go
// internal/repository/mock.go
type MockItemRepository struct {
    mock.Mock
}

func (m *MockItemRepository) Create(ctx context.Context, item Item) (Item, error) {
    args := m.Called(ctx, item)
    return args.Get(0).(Item), args.Error(1)
}

// Usage in service test
func TestItemService_CreateItem(t *testing.T) {
    mockRepo := new(MockItemRepository)
    mockRepo.On("Create", mock.Anything, mock.Anything).
        Return(domain.Item{ID: "123"}, nil)

    service := NewItemService(mockRepo)
    result, err := service.CreateItem(ctx, item)

    assert.NoError(t, err)
    assert.Equal(t, "123", result.ID)
    mockRepo.AssertExpectations(t)
}
```

### Service Mock

```go
// internal/service/mock.go
type MockItemService struct {
    mock.Mock
}

func (m *MockItemService) CreateItem(ctx context.Context, item domain.Item) (domain.Item, error) {
    args := m.Called(ctx, item)
    return args.Get(0).(domain.Item), args.Error(1)
}

// Usage in handler test
func TestHandler_CreateItem(t *testing.T) {
    mockService := new(MockItemService)
    mockService.On("CreateItem", mock.Anything, mock.Anything).
        Return(domain.Item{ID: "123"}, nil)

    handler := NewItemHandler(mockService)
    // Test HTTP request handling...
}
```

## Vantagens

**Velocidade:** Testes unitários executam em <10ms vs >1s com Docker

**Determinismo:** Mocks retornam valores exatos, não há flakiness

**Cenários de Erro:** Fácil simular errors que seriam difíceis de reproduzir com banco real:
```go
mockRepo.On("Create", mock.Anything, mock.Anything).
    Return(domain.Item{}, errors.New("connection timeout"))
```

**Paralelização:** Testes podem rodar em paralelo sem conflitos de estado

**CI/CD:** PR checks completam em segundos vs minutos

## Trade-offs

**Mais código:** Cada mock requer definição de interface

**Fragilidade:** Mocks precisam ser atualizados quando interfaces mudam

**Over-testing:** Testes muito acoplados à implementação podem quebrar com refactors

## Mitigações

**Auto-generated Mocks:** Ferramentas como `mockgen` podem gerar mocks automaticamente

**Testes de Integração:** Mantém um pequeno conjunto de testes com dependencies reais para validar comportamento end-to-end

**Black-box Testing:** Testes devem focar em comportamento observável (inputs → outputs), não detalhes de implementação

## Consequências

### Positivas
- Feedback loop rápido (desenvolvedor não espera por testes lentos)
- Test coverage pode atingir >80% sem impactar velocity
- CI/CD eficiente (PR checks em segundos)
- Facilita TDD (red-green-refactor cycle rápido)

### Negativas
- Curva de aprendizado para mocks
- Testes podem enganar se mocks não refletem comportamento real
- Requer disciplina para manter mocks sincronizados com interfaces

## Referências
- [testify/mock Documentation](https://github.com/stretchr/testify/blob/master/mock/mock.go)
- [The Mythical Man-Month: No Silver Bullet](https://en.wikipedia.org/wiki/No_Silver_Bullet)
