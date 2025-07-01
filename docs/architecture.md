# Documento de Arquitetura do Projeto: List Manager API

## 1. Introdução

Este documento descreve a arquitetura da "List Manager API", um serviço de backend robusto para gerenciar listas de itens, especificamente produtos e usuários. O projeto foca na integração com MongoDB como o principal armazenamento de dados, implementando operações CRUD e aderindo aos princípios da Arquitetura Limpa.

## 2. Metas e Requisitos

Os principais objetivos e requisitos deste projeto incluem:

- Implementar uma API de gerenciamento de listas em Go.
- Integrar MongoDB para persistência de dados.
- Definir modelos de dados para `Products` e `Users` com campos como `_id`, `name`, `active` para produtos e `_id`, `created_by` para usuários.
- Implementar operações transacionais quando aplicável.
- Aderir à Arquitetura Limpa: separação de handlers/controllers, services/use cases, repositories/data access e modelos de domínio.
- Seguir práticas idiomáticas de Go, design modular, testabilidade e melhores práticas de desenvolvimento backend.
- Garantir tratamento de erros adequado, injeção de dependência e propagação de contexto.
- Desenvolver testes de unidade e integração abrangentes.

## 3. Visão Geral da Arquitetura

A arquitetura do projeto segue os princípios da Arquitetura Limpa (Clean Architecture), promovendo a separação de preocupações e a testabilidade. O sistema é dividido em camadas distintas:

- **Handlers (ou Controllers)**: Responsáveis por receber requisições HTTP, parsear entradas e chamar a camada de Serviço.
- **Services (ou Use Cases)**: Contêm a lógica de negócio principal do aplicativo, orquestrando operações e interagindo com a camada de Repositório.
- **Repositories (ou Data Access)**: Abstraem a complexidade da persistência de dados, fornecendo interfaces para operações CRUD.
- **Domain**: Define as entidades de negócio e suas regras.
- **Database (MongoDB)**: O armazenamento de dados persistente.

A comunicação entre as camadas é feita através de interfaces, garantindo baixo acoplamento.

## 4. Componentes Chave

- **`cmd/api/handlers`**: Lida com as requisições HTTP, roteamento e validação de entrada. Contém `item.go` para operações relacionadas a itens, `errors.go` para tratamento de erros HTTP, `middleware.go` para middlewares, e `model.go` para modelos de requisição/resposta.
- **`internal/service`**: Contém a lógica de negócio e os casos de uso. `item.go` define os serviços para operações de item. `parser.go` para parsear e validar dados de entrada.
- **`internal/repository`**: Define interfaces para abstrair o armazenamento de dados.
  - **`internal/repository/mongodb`**: Implementação concreta das interfaces de repositório usando MongoDB. Inclui `repository.go` para operações CRUD de `Product` e `User`.
  - **`internal/repository/local`**: Uma implementação de repositório local (em memória) que pode ser substituída ou aumentada pela implementação MongoDB.
- **`internal/domain`**: Contém as definições de modelos de dados, como `item.go` para as estruturas `Product` e `User`.
- **`internal/database/mongodb`**: Gerencia a conexão e as operações de baixo nível com o cliente MongoDB, incluindo `client.go` e `interfaces.go`.

## 4.1 Visão Geral dos Componentes Chave e Estrutura de Pacotes

A estrutura do projeto é organizada para seguir os princípios da Arquitetura Limpa, com uma clara separação de responsabilidades. Abaixo, detalhamos os principais componentes e a organização de seus respectivos pacotes.

### 4.1.1 Estrutura de Pacotes e Resumo

O projeto está organizado nos seguintes diretórios e pacotes, cada um com uma função específica:

- **`cmd/`**: Contém as entradas principais da aplicação.
  - **`cmd/api/`**: Onde a aplicação principal da API é definida.
    - **`cmd/api/handlers/`**: Pacote responsável por lidar com as requisições HTTP, roteamento e validação de entrada.
      - `cors.go`: Configurações de Cross-Origin Resource Sharing (CORS).
      - `errors.go`: Definições para tratamento de erros HTTP.
      - `handlers.go`: Definições gerais dos manipuladores de requisição.
      - `handlers_test.go`: Testes unitários para os manipuladores HTTP.
      - `item.go`: Manipuladores específicos para operações de itens (produtos e usuários).
      - `middleware.go`: Implementa middlewares HTTP para funcionalidades como autenticação e logging.
      - `model.go`: Definições de modelos de dados para requisições e respostas HTTP.
      - `parser.go`: Utilitários para parsear dados das requisições.
      - `version.go`: Informações de versão da API.
    - **`cmd/api/main.go`**: O ponto de entrada principal da aplicação API.
    - **`cmd/api/server/`**: Contém a configuração e inicialização do servidor HTTP.
      - `server.go`: Responsável por configurar e iniciar o servidor.
- **`docs/`**: Contém a documentação do projeto.
  - `architecture.md`: Este documento de arquitetura.
- **`internal/`**: Código interno da aplicação, não destinado a ser exposto publicamente.
  - **`internal/database/`**: Contém a abstração para a interação com o banco de dados.
    - **`internal/database/mongodb/`**: Implementações específicas para o MongoDB.
      - `client.go`: Lógica para conexão e gerenciamento do cliente MongoDB.
      - `interfaces.go`: Interfaces para as operações do cliente MongoDB.
      - `wrappers.go`: Funções de "wrapper" para operações de baixo nível do MongoDB.
  - **`internal/domain/`**: Define as entidades do domínio e as regras de negócio puras.
    - `item.go`: Definição das estruturas de dados `Product` e `User`. Inclui a função `generateID()` para criar IDs compatíveis com `ObjectID` do MongoDB, e métodos de lógica de negócio como `IsEmpty()` e `IsActive()`.
    - `item_test.go`: Testes unitários para as entidades de domínio.
  - **`internal/repository/`**: Camada de abstração para persistência de dados.
    - `errors.go`: Erros específicos da camada de repositório.
    - `mock.go`: Implementações mock para facilitar os testes unitários dos repositórios.
    - `model.go`: Modelos de dados utilizados internamente pela camada de repositório.
    - `repository.go`: Interfaces que definem os contratos para operações de persistência de dados.
    - **`internal/repository/local/`**: Implementação de repositório em memória para desenvolvimento/testes rápidos.
      - `service.go`: Serviço do repositório local.
    - **`internal/repository/mongodb/`**: Implementação concreta das interfaces de repositório para MongoDB.
      - `repository.go`: Lógica para persistência de `Product` e `User` no MongoDB. Inclui a implementação de transações MongoDB para operações multi-documento/coleção, como `CreateItemWithUser`, garantindo a atomicidade.
      - `repository_test.go`: Testes unitários para o repositório MongoDB.
  - **`internal/service/`**: Contém a lógica de negócio principal (casos de uso).
    - `errors.go`: Erros específicos da camada de serviço.
    - `item.go`: Lógica de negócio para operações de itens (produtos e usuários).
    - `mock.go`: Implementações mock para testes unitários dos serviços.
    - `parser.go`: Utilitários para parsear e validar dados dentro da camada de serviço.
    - `service.go`: Definições gerais dos serviços.
    - `service_test.go`: Testes unitários para os serviços.
- **`memory-bank/`**: Diretório que contém documentos de contexto e informações de projeto.
  - `activeContext.md`, `productContext.md`, `progress.md`, `projectbrief.md`, `systemPatterns.md`, `techContext.md`: Documentos diversos de contexto.
- **`docker-compose.yml`**: Arquivo de configuração do Docker Compose para orquestração de serviços (ex: MongoDB).
- **`go.mod`**: Módulo Go, define as dependências do projeto.
- **`go.sum`**: Checksums das dependências do módulo Go.
- **`Makefile`**: Arquivo Makefile para automatizar tarefas de construção, teste e implantação.
- **`README.md`**: Documento de introdução ao projeto.

## 5. Persistência de Dados

O MongoDB é o banco de dados principal, configurado via `docker-compose.yml`. As entidades `Product` e `User` são persistidas com tags `bson` para mapeamento correto. As operações transacionais são implementadas quando necessárias para garantir a integridade dos dados, como visto na função `CreateItemWithUser` no repositório MongoDB.

## 6. Estratégia de Testes

- **Testes de Unidade**: Priorizam testes de unidade extensivos para a camada de repositório usando mocks para dependências do MongoDB. Isso evita a necessidade de uma instância de Docker para testes de unidade e garante isolamento.
- **Testes de Integração**: Serão considerados após a estabilidade dos testes de unidade, com potencial uso de `testcontainers-go` para simular um ambiente de banco de dados real.
- **Testes Manuais**: Validação inicial da funcionalidade através de testes manuais para confirmar operações CRUD básicas.

## 7. Tratamento de Erros

Um tratamento de erros robusto é garantido, com o encapsulamento e propagação de erros para facilitar a rastreabilidade e depuração. Erros específicos do domínio e do sistema são tratados de forma apropriada.

## 8. Injeção de Dependência

A injeção de dependência é realizada através de funções construtoras, garantindo que as dependências sejam passadas de forma explícita e controlada, o que melhora a testabilidade e a modularidade do código.

## 9. Observabilidade

A aplicação utiliza `go.uber.org/zap` para logging estruturado, o que facilita a análise e depuração de logs em ambientes de produção. Os logs são configurados para fornecer informações detalhadas sobre o fluxo da aplicação.

**Futura Integração com OpenTelemetry**: Há planos para integrar o OpenTelemetry para tracing distribuído e métricas. Esta integração permitirá uma visibilidade mais profunda do desempenho e do fluxo das requisições através dos serviços, complementando o logging existente com rastreamento de ponta a ponta e coleta de métricas padronizadas.

## 10. Ordem dos Middlewares

No `cmd/api/server/server.go`, os middlewares são aplicados na seguinte ordem para garantir o comportamento correto:

1.  **`CORSMiddleware`**: Aplicado primeiro para manipular as requisições de pré-voo (preflight requests) do CORS antes de qualquer outra lógica de middleware ou roteamento.
2.  **`LoggingMiddleware`**: Aplicado após o CORS para registrar as requisições que passaram pela verificação de CORS.
3.  **Router (`mux.NewRouter()`)**: O roteador é o último a ser aplicado, garantindo que as requisições sejam logadas e que o CORS seja tratado antes de o roteamento ser realizado.

---
