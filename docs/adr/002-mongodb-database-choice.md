# ADR-002: MongoDB como Banco de Dados

## Status
Aceito

## Contexto
O projeto List Manager API necessita persistir dados de itens (produtos/usuarios) com as seguintes características:
- Schema flexível para evolução do modelo de dados
- Operações de leitura/escrita com baixa latência
- Suporte a transações para operações atômicas
- Escalabilidade horizontal para crescimento futuro

## Decisão
Utilizar MongoDB como banco de dados principal para persistência de todos os dados da aplicação.

## Rationale
**Flexibilidade de Schema:** Documentos JSON/BSON permitem evolução do schema sem migrações complexas. Campos como `observation` podem ser adicionados sem impactar documentos existentes.

**Modelo de Dados Natural:** Estrutura de itens com propriedades variáveis se alinha perfeitamente com modelo de documentos. Um item = um documento.

**Performance:** Operações de leitura/escrita otimizadas com índices. Expressões de consulta poderosas para filtros complexos.

**Transações ACID:** MongoDB 4.0+ suporta multi-document transactions, garantindo atomicidade para operações complexas como `CreateItemWithUser`.

**Escalabilidade:** Sharding nativo para crescimento horizontal. Replica sets para alta disponibilidade.

**Ecossistema Cloud:** Managed services disponíveis (MongoDB Atlas, AWS DocumentDB) para deploy sem overhead de manutenção.

**Integração Go:** Driver oficial `go.mongodb.org/mongo-driver` maduro com suporte completo a features modernas.

## Consequências

### Positivas
- Desenvolvimento ágil sem schema migrations
- Índices automáticos no `_id` para performance O(1) em buscas
- Aggregation framework para consultas complexas
- Change streams para reactive updates
- TTL indexes para expiração automática de dados
- Backup/restore simplificado com ferramentas como `mongodump`

### Negativas
- Consistency eventual em setups distribuídos (configuração padrão)
- Consumo de memória maior que bancos relacionais (WiredTiger cache)
- Limitações em joins comparado a SQL (requer application-level joins ou $lookup)
- Document size limit de 16MB (não é problema para o caso de uso)

## Alternativas Consideradas
- **PostgreSQL:** Relacional maduro mas schema rígido requer migrations para cada mudança
- **DynamoDB:** Serverless mas custo elevado e limitações de query
- **Redis:** In-memory excelente para cache mas não persistente por padrão
- **Cassandra:** Escalabilidade excepcional mas modelagem complexa e latência maior
