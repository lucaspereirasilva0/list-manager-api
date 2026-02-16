# ADR-003: Clean Architecture

## Status
Aceito

## Contexto
Aplicações enterprise crescem em complexidade ao longo do tempo. Codebases que não separam interesses tornam-se difíceis de:
- Testar (dependências externas acopladas)
- Evoluir (mudanças em uma camada impactam outras)
- Escalar (lógica de negócio misturada com infraestrutura)
- Entender (responsabilidades mistas em um único arquivo)

## Decisão
Adotar Clean Architecture com separação explícita em 4 camadas: **domain → repository → service → handlers**.

## Arquitetura

```
┌─────────────────────────────────────────────────┐
│                  Handlers                       │
│            (HTTP, Request/Response)             │
└──────────────────┬──────────────────────────────┘
                   │ depends on
┌──────────────────▼──────────────────────────────┐
│                  Services                       │
│         (Business Logic, Use Cases)              │
└──────────────────┬──────────────────────────────┘
                   │ depends on
┌──────────────────▼──────────────────────────────┐
│               Repositories                      │
│         (Data Access Abstraction)               │
└──────────────────┬──────────────────────────────┘
                   │ depends on
┌──────────────────▼──────────────────────────────┐
│                  Domain                         │
│           (Business Entities)                   │
└─────────────────────────────────────────────────┘
```

### Camadas

**1. Domain (`internal/domain/`)**
- Entidades puras sem dependências externas
- Regras de negócio core (ex: `IsActive()`, `IsEmpty()`)
- Estruturas de dados com `bson` tags para mapeamento

**2. Repository (`internal/repository/`)**
- Interfaces definindo contratos de persistência
- Implementações: `local/` (in-memory), `mongodb/` (production)
- Abstração completa da tecnologia de banco

**3. Service (`internal/service/`)**
- Casos de uso e orquestração de lógica
- Validação de regras de negócio
- Transações e operações compostas

**4. Handlers (`cmd/api/handlers/`)**
- HTTP request/response handling
- Parsing e validação de input
- Error handling e status codes

### Dependency Rule
**Dependências apontam para dentro:** Camadas externas dependem de interfaces definidas em camadas internas. Nenhuma camada interna sabe sobre HTTP ou MongoDB.

## Rationale
**Testabilidade:** Cada camada pode ser testada isoladamente com mocks. Services podem ser testados sem banco de dados real.

**Independência de Tecnologia:** Mudar de MongoDB para PostgreSQL requer apenas nova implementação de Repository. Services e Handlers permanecem inalterados.

**Manutenibilidade:** Responsabilidades claras facilitam encontrar onde implementar mudanças.

**Escalabilidade:** Camadas podem ser escaladas independentemente (ex: handlers como microserviços separados).

**Reuso:** Domain entities podem ser reutilizadas em diferentes contextos (CLI, HTTP, gRPC).

## Consequências

### Positivas
- Testes unitários rápidos sem dependências externas
- Troca de implementações (ex: MongoDB → PostgreSQL) sem impacto em business logic
- Codebase organizado facilita onboarding
- Baixo acoplamento entre componentes
- Fácil adicionar novas interfaces (CLI, gRPC) reaproveitando service layer

### Negativas
- Boilerplate inicial para definir interfaces
- Curva de aprendizado para desenvolvedores não familiarizados
- Indireção adicional pode parecer over-engineering para projetos simples
- Mais arquivos/diretórios para navegar

## Implementação
- Interfaces em `internal/repository/repository.go` e `internal/service/service.go`
- Implementações concretas em subdiretórios (ex: `mongodb/`)
- DI via constructors em `cmd/api/main.go`
- Mocks gerados automaticamente via testify/mock

## Referências
- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
