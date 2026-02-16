# ADR-008: MongoDB Transactions

## Status
Aceito

## Contexto
Operações de negócio frequentemente requerem atomicidade across múltiplos documentos:
- Criar item e registrar who/when criou
- Atualizar item e log de auditoria
- Operações compostas que devem suceder ou falhar juntas

Sem transações,operações parciais podem deixar dados inconsistentes:
- Item criado mas audit log falhou
- Item atualizado mas contador de modificações não

## Decisão
Utilizar MongoDB Transactions (multi-document ACID) para operações que requerem atomicidade.

## Implementação

### MongoDB Session e Transaction

```go
func (r *MongoRepository) CreateItemWithUser(ctx context.Context, item Item, user User) error {
    session, err := r.client.StartSession()
    if err != nil {
        return err
    }
    defer session.EndSession(ctx)

    // Inicia transaction
    result, err := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
        // Operação 1: Insert item
        if _, err := r.itemsCollection.InsertOne(sessCtx, item); err != nil {
            return nil, err
        }

        // Operação 2: Insert user/creator
        if _, err := r.usersCollection.InsertOne(sessCtx, user); err != nil {
            return nil, err
        }

        return nil, nil
    })

    return err
}
```

### Propriedades ACID

**Atomicity:** Ambos inserts sucemem ou nenhum

**Consistency:** Estado final é válido (item + user criados)

**Isolation:** Outras operações não veem estado parcial

**Durability:** Commit garante persistência

## Quando Usar Transactions

### ✅ Use Transactions
- Operações multi-document/collection
- Requisitos de consistência estrita
- Auditoria que deve atômica
- Operações financeiras ou críticas

### ❌ Evite Transactions
- Single document operations (já são atômicas)
- Operações em collections diferentes sem relacionamento
- Performance é crítico e consistência eventual é aceitável

## Performance Considerations

Transactions têm overhead:
- Require session management
- Lock duration maior
- Retry logic built-in

**Guideline:** Usar apenas quando necessário. Single document writes não precisam de transactions.

## Retry Logic

MongoDB driver implementa automaticamente retryable writes para transient errors:
- Network glitches
- Primary stepdown during replica set election
- Transaction conflicts

```go
result, err := session.WithTransaction(ctx, callback, options.Transaction().SetRetry(true))
```

## Limitações

**Transaction Size:** Limite de operações por transaction (configurável, default é baixo)

**Cross-shard:** Transactions não funcionam across shards (unsharded collections)

**Session Timeout:** Transactions tem timeout (configurável via `maxTimeMS`)

## Alternativas

**Two-Phase Commit:** Pattern mais complexo para MongoDB pre-4.0

**Application-level Compensation:** Executar rollback manual (não recomendado, error-prone)

**Event Sourcing:** Log de eventos para reconstruir estado (maior complexidade)

## Consequências

### Positivas
- Consistência forte garantida para operações críticas
- Simplifica error handling (rollback automático)
- Padronização com práticas enterprise RDBMS

### Negativas
- Performance overhead (deve ser usado com cautela)
- Requer replica sets (não funcionam em standalone MongoDB)
- Complexidade adicional de session management

## Status da Implementação

No projeto atual, transactions estão disponíveis via `mongo.SessionContext` no repository layer. Operações simples de CRUD não utilizam transactions (single document operations são já atômicas). Operações complexas podem ser implementadas com transactions quando necessário.

## Referências
- [MongoDB Transactions Documentation](https://www.mongodb.com/docs/manual/core/transactions/)
- [Go MongoDB Driver - Sessions](https://www.mongodb.com/docs/drivers/go/current/fundamentals/transactions/)
