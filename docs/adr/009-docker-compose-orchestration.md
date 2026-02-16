# ADR-009: Docker Compose Orchestration

## Status
Aceito

## Contexto
Ambiente de desenvolvimento local necessita:
- Serviços dependentes (MongoDB)
- Facilidade de setup (developers não devem instalar MongoDB manualmente)
- Consistência entre ambientes (dev, staging, prod)
- Isolamento de dependências
- Facilidade de teardown e cleanup

## Decisão
Utilizar Docker Compose para orquestrar serviços dependentes da aplicação.

## Serviços Orquestrados

### 1. MongoDB
```yaml
mongodb:
  image: mongo:8.0
  container_name: list-manager-mongodb
  ports:
    - "27017:27017"
  environment:
    MONGO_INITDB_ROOT_USERNAME: root
    MONGO_INITDB_ROOT_PASSWORD: password
  volumes:
    - mongodb_data:/data/db
```

### 2. Mongo Express (UI)
```yaml
mongo-express:
  image: mongo-express:1.0
  container_name: list-manager-mongo-express
  ports:
    - "8081:8081"
  environment:
    ME_CONFIG_MONGODB_URL: mongodb://root:password@mongodb:27017/
  depends_on:
    - mongodb
```

### 3. API (Opção Futura)
```yaml
api:
  build:
    context: .
    dockerfile: Dockerfile
  ports:
    - "8085:8085"
  environment:
    MONGO_URI: mongodb://root:password@mongodb:27017/
    MONGO_DB_NAME: listmanager
    PORT: 8085
  depends_on:
    - mongodb
```

## Rationale

### Local Development
**Zero Setup:** Clone repo → `docker-compose up` → ambiente completo

**Consistência:** Todos developers usam mesma versão de MongoDB

**Isolamento:** Sem conflitos com MongoDB local instalado

**UI de Administração:** Mongo Express para visualizar dados sem CLI

### Paridade com Produção
Docker images utilizadas em dev podem ser as mesmas de staging/production, reduzindo surpresas de "works on my machine".

### Volumes
```yaml
volumes:
  mongodb_data:
```
Persiste dados entre container restarts, permitindo desenvolvedor manter dados locais durante work-in-progress.

## Comandos Principais

```bash
# Start services
docker-compose up -d

# View logs
docker-compose logs -f mongodb

# Stop services
docker-compose down

# Stop e remove volumes (limpeza completa)
docker-compose down -v

# Rebuild after changes
docker-compose up -d --build
```

## Network

Docker Compose cria network bridge `list-manager-api_default` automaticamente. Serviços se comunicam via service names como hostnames:
- API connecta em `mongodb:27017` (não `localhost:27017`)

## Alternativas Consideradas

### Docker Manual
- ❌ Requer múltiplos `docker run` commands
- ❌ Network configuration manual
- ❌ Orchestration manual

### MongoDB Instalado Localmente
- ❌ Setup manual por developer
- ❌ Versões不一致 entre devs
- ❌ Difícil cleanup e switch entre versões

### Kubernetes
- ❌ Overkill para desenvolvimento local
- ❌ Curva de aprendizado íngreme
- ❌ Resource-intensive

## Consequências

### Positivas
- Onboarding instantâneo para novos developers
- Ambientes consistentes
- Facilidade de reset de dados
- CI/CD simplificado (mesmo compose pode rodar em GitHub Actions)

### Negativas
- Requer Docker instalado (resource overhead)
- Slight latency adicional comparado a native (mas não é significativo para desenvolvimento)

## Melhorias Futuras

**Healthchecks:**
```yaml
healthcheck:
  test: echo 'db.runCommand("ping").ok' | mongosh localhost:27017/test --quiet
  interval: 10s
  timeout: 5s
  retries: 5
```

**Automatic Seeding:**
Script de inicialização para popular dados de teste automaticamente.

**Hot Reload:** Para desenvolvimento com rebuild automático da API.

## Referências
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [MongoDB Official Docker Image](https://hub.docker.com/_/mongo)
