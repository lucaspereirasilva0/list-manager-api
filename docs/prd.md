# Product Requirements Document - List Manager API

## Visão Geral

API RESTful para gerenciamento de listas de itens (produtos/usuarios) com suporte a operações em massa. A aplicação fornece endpoints para CRUD completo, health checks, e atualizações em lote do status ativo dos itens.

## Requisitos Funcionais

### RF001: Criar Item
**Descrição:** Sistema deve permitir criação de novos itens com nome, status ativo e observação opcional.

**Critérios de Aceite:**
- Item deve possuir ID gerado automaticamente (UUID/ObjectID)
- Campo `name` é obrigatório e não pode ser vazio
- Campo `active` é obrigatório (default: true)
- Campo `observation` é opcional
- Timestamps `createdAt` e `updatedAt` devem ser gerados automaticamente

**Endpoint:** `POST /item`

---

### RF002: Atualizar Item
**Descrição:** Sistema deve permitir atualização de item existente por ID.

**Critérios de Aceite:**
- Todos os campos podem ser atualizados exceto ID
- `updatedAt` deve ser atualizado automaticamente
- Retornar 404 se item não existe

**Endpoint:** `PUT /item?id={id}`

---

### RF003: Deletar Item
**Descrição:** Sistema deve permitir remoção de item por ID.

**Critérios de Aceite:**
- Item deve ser permanentemente removido do banco de dados
- Retornar 404 se item não existe
- Retornar 204 No Content em sucesso

**Endpoint:** `DELETE /item?id={id}`

---

### RF004: Buscar Item por ID
**Descrição:** Sistema deve permitir consulta de item específico.

**Critérios de Aceite:**
- Busca deve ser por ID exato
- Retornar 404 se item não existe
- Retornar item completo em JSON

**Endpoint:** `GET /item?id={id}`

---

### RF005: Listar Todos os Itens
**Descrição:** Sistema deve permitir listagem completa de itens.

**Critérios de Aceite:**
- Retornar array com todos os itens
- Array vazio se não houver itens
- Itens ordenados por createdAt descendente (mais recentes primeiro)

**Endpoint:** `GET /items`

---

### RF006: Atualizar Status Ativo em Massa
**Descrição:** Sistema deve permitir atualização do campo `active` de todos os itens em uma única operação.

**Critérios de Aceite:**
- Todos os itens devem ser atualizados com o mesmo status
- Retornar contagem de itens matched e modified
- Operação deve ser atômica (transaction)

**Endpoint:** `PUT /items/active`

**Request:**
```json
{
  "active": true
}
```

**Response:**
```json
{
  "matchedCount": 150,
  "modifiedCount": 42
}
```

---

### RF007: Health Check
**Descrição:** Sistema deve fornecer endpoint de verificação de saúde.

**Critérios de Aceite:**
- Verificar status da aplicação (up/down)
- Verificar conexão com MongoDB (connected/disconnected)
- Retornar timestamp da verificação
- Incluir checks adicionais conforme necessário

**Endpoint:** `GET /healthz`

**Response:**
```json
{
  "status": "up",
  "server": "up",
  "database": {
    "status": "connected"
  },
  "timestamp": "2026-02-13T10:30:00Z",
  "checks": {
    "mongodb": {
      "status": "passed"
    }
  }
}
```

---

### RF008: Versão da Aplicação
**Descrição:** Sistema deve fornecer endpoint de versão para PWA auto-update.

**Critérios de Aceite:**
- Retornar versão atual da aplicação
- Formato JSON compatível com PWA version checking

**Endpoint:** `GET /_app/version.json`

**Response:**
```json
{
  "version": "1.0.0"
}
```

---

## Requisitos Não-Funcionais

### RNF001: Performance
**Descrição:** Operações simples devem retornar em menos de 100ms.

**Métricas:**
- CRUD operations: < 100ms (p95)
- Health check: < 50ms (p95)
- List items (primeiras 100): < 200ms (p95)

### RNF002: Escalabilidade
**Descrição:** Sistema deve suportar múltiplos clientes concorrentes.

**Requisitos:**
- Suportar mínimo de 50 concurrent requests sem degradação significativa
- Connection pooling configurado para MongoDB
- Stateless design para horizontal scaling

### RNF003: Disponibilidade
**Descrição:** Sistema deve fornecer health check para monitoramento.

**Requisitos:**
- Health check endpoint sempre disponível
- Health check deve verificar conectividade com banco de dados
- Suportar Kubernetes/Docker health checks

### RNF004: Segurança
**Descrição:** Sistema deve implementar segurança básica via configuração.

**Requisitos:**
- CORS configurável via variáveis de ambiente
- Environment variables para secrets (MongoDB credentials)
- Logs não devem conter informações sensíveis

### RNF005: Observabilidade
**Descrição:** Sistema deve implementar logging estruturado.

**Requisitos:**
- Logs estruturados com Zap
- Níveis de log configuráveis (debug, info, warn, error)
- Contexto propagado através de request tracing

---

## Modelo de Dados

### Item Entity

```json
{
  "id": "string (UUID/ObjectID)",
  "name": "string (required, non-empty)",
  "active": "boolean (required, default: true)",
  "observation": "string | null (optional)",
  "createdAt": "datetime (ISO 8601)",
  "updatedAt": "datetime (ISO 8601)"
}
```

**Validation Rules:**
- `id`: Auto-generated, imutável
- `name`: 1-255 caracteres, trimmed
- `active`: Boolean, obrigatório
- `observation`: Opcional, pode ser null
- `createdAt`: Auto-generated na criação
- `updatedAt`: Auto-generated na criação, atualizado em updates

---

## Casos de Uso

### UC001: Gerente Cadastra Novo Produto
**Ator:** Gerente de Produtos

**Pré-condições:**
- Gerente está autenticado
- Sistema está operacional

**Fluxo Principal:**
1. Gerente acessa interface de criação
2. Gerente insere nome do produto
3. Gerente opcionalmente adiciona observação
4. Gerente confirma criação
5. Sistema valida dados
6. Sistema cria item com ID único
7. Sistema retorna confirmação com item criado

**Pós-condições:**
- Item persistido no banco de dados
- Logs registram criação

**Fluxos Alternativos:**
- 5a. Nome vazio: Sistema retorna 400 Bad Request
- 5b. Erro de banco: Sistema retorna 500 Internal Server Error

---

### UC002: Gerente Atualiza Informações de Produto Existente
**Ator:** Gerente de Produtos

**Pré-condições:**
- Produto existe no sistema
- Gerente está autenticado

**Fluxo Principal:**
1. Gerente seleciona produto para editar
2. Sistema exibe dados atuais do produto
3. Gerente modifica campos desejados
4. Gerente confirma alterações
5. Sistema atualiza produto no banco
6. Sistema retorna produto atualizado

**Pós-condições:**
- Produto atualizado no banco
- `updatedAt` refreshado
- Logs registram modificação

---

### UC003: Gerente Remove Produto da Lista
**Ator:** Gerente de Produtos

**Pré-condições:**
- Produto existe no sistema
- Gerente está autenticado

**Fluxo Principal:**
1. Gerente seleciona produto para remover
2. Sistema solicita confirmação
3. Gerente confirma remoção
4. Sistema deleta produto do banco
5. Sistema retorna 204 No Content

**Pós-condições:**
- Produto permanentemente removido
- Logs registram deleção

---

### UC004: Gerente Consulta Produto Específico
**Ator:** Gerente de Produtos

**Pré-condições:**
- Produto existe no sistema
- Gerente está autenticado

**Fluxo Principal:**
1. Gerente fornece ID do produto
2. Sistema busca produto no banco
3. Sistema retorna dados do produto

**Pós-condições:**
- Nenhuma (leitura não altera estado)

**Fluxos Alternativos:**
- 2a. Produto não existe: Sistema retorna 404 Not Found

---

### UC005: Gerente Visualiza Todos os Produtos
**Ator:** Gerente de Produtos

**Pré-condições:**
- Gerente está autenticado

**Fluxo Principal:**
1. Gerente solicita listagem completa
2. Sistema busca todos os produtos no banco
3. Sistema retorna array ordenado (mais recentes primeiro)

**Pós-condições:**
- Nenhuma (leitura não altera estado)

---

### UC006: Gerente Ativa/Desativa Todos os Produtos de Uma Vez
**Ator:** Gerente de Produtos

**Pré-condições:**
- Gerente está autenticado
- Pelo menos um produto existe

**Fluxo Principal:**
1. Gerente acessa função de atualização em massa
2. Gerente seleciona status desejado (ativo/inativo)
3. Gerente confirma operação
4. Sistema atualiza todos os produtos em transação atômica
5. Sistema retorna estatísticas (matchedCount, modifiedCount)

**Pós-condições:**
- Todos os produtos possuem mesmo status `active`
- Logs registram operação em massa

---

## Definição de Done

### Sprint Backlog Item está completo quando:
- [ ] Code review aprovado
- [ ] Testes unitários passando (coverage > 80%)
- [ ] Testes de integração passando
- [ ] Documentação atualizada (ADR se aplicável, OpenAPI)
- [ ] Merge para main branch

### Feature está completa quando:
- [ ] Todos os requisitos funcionais implementados
- [ ] Todos os requisitos não-funcionais validados
- [ ] Documentação (PRD, OpenAPI) reflete implementação atual
- [ ] Deploy em staging validado
- [ ] Sign-off de product owner

---

## Roadmap

### Versão 1.0 (Current)
- [x] CRUD completo de itens
- [x] MongoDB integration
- [x] Clean Architecture
- [x] Docker Compose setup
- [x] Health check com verificação de MongoDB
- [x] Bulk update endpoint
- [x] Observation field
- [x] Version endpoint for PWA

### Versão 1.1 (Planned)
- [ ] OpenTelemetry integration (tracing, metrics)
- [ ] Enhanced validation com custom error types
- [ ] Pagination para list items
- [ ] Filter por active status
- [ ] Search por nome

### Versão 2.0 (Future)
- [ ] Authentication/Authorization (JWT)
- [ ] Rate limiting
- [ ] GraphQL alternative endpoint
- [ ] WebSocket para realtime updates
- [ ] Multi-language support (i18n)

---

## Referências

- [OpenAPI 3.0 Specification](https://swagger.io/specification/)
- [RFC 7617 - JSON Web Token (JWT)](https://tools.ietf.org/html/rfc7519)
- [OWASP API Security Top 10](https://owasp.org/www-project-api-security/)
