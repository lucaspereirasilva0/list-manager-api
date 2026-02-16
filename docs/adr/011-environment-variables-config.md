# ADR-011: Environment Variables Config

## Status
Aceito

## Contexto
Aplicações production-ready necessitam configuração externalizada:
- Different values per environment (dev, staging, prod)
- Secrets não devem ser committed em código
- Facilidade de deployment em diferentes infraestruturas (local, Docker, cloud)

Hardcoding configuration values é anti-pattern:
- Difficil mudar valores sem recompilar
- Secrets expostos em version control
- Impossível rodar múltiplas instâncias com configurações diferentes

## Decisão
Utilizar variáveis de ambiente para toda configuração externa da aplicação.

## Variáveis de Ambiente

### Obrigatórias

| Variável | Descrição | Exemplo | Default |
|-----------|-----------|---------|---------|
| `MONGO_URI` | Connection string MongoDB | `mongodb://user:pass@localhost:27017/` | - |
| `MONGO_DB_NAME` | Nome do database | `listmanager` | - |
| `PORT` | Porta HTTP server | `8085` | `8080` |

### Opcionais (Futuro)

| Variável | Descrição | Exemplo | Default |
|-----------|-----------|---------|---------|
| `LOG_LEVEL` | Nível de log (debug, info, warn, error) | `info` | `info` |
| `ENVIRONMENT` | Ambiente (dev, staging, prod) | `production` | `development` |
| `CORS_ORIGINS` | CORS allowed origins | `https://example.com` | `*` |

## Implementação

### Environment File (.env)
Para desenvolvimento local:

```bash
# .env (não committed)
MONGO_URI=mongodb://root:password@localhost:27017/
MONGO_DB_NAME=listmanager
PORT=8085
LOG_LEVEL=debug
```

### Loading no Código

```go
// cmd/api/main.go

func main() {
    mongoURI := os.Getenv("MONGO_URI")
    if mongoURI == "" {
        log.Fatal("MONGO_URI environment variable is required")
    }

    dbName := os.Getenv("MONGO_DB_NAME")
    if dbName == "" {
        log.Fatal("MONGO_DB_NAME environment variable is required")
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // default
    }

    // Initialize com configs
    // ...
}
```

### Docker Compose

```yaml
# docker-compose.yml
services:
  api:
    environment:
      - MONGO_URI=mongodb://root:password@mongodb:27017/
      - MONGO_DB_NAME=listmanager
      - PORT=8085
```

### Deployment (Render)

```yaml
# render.yaml
services:
  - type: web
    name: list-manager-api
    envVars:
      - key: MONGO_URI
        sync: false # Secret
      - key: MONGO_DB_NAME
        value: listmanager
      - key: PORT
        value: 8085
```

## Rationale

### Separation of Concerns
Configuration é orthogonal ao código. Mudanças de configuração não devem requerer recompilation.

### Environment Parity
Dev usa Docker Compose com env vars → Production usa cloud platform com env vars. Mesmo mecanismo.

### Security
Secrets (MONGO_URI com senha) estão em environment, não em código.

### Flexibility
Easy para:
- Deploy em múltiplas regiões com configs diferentes
- A/B testing com feature flags via env
- Rollback para versão anterior com config diferente

## Alternativas Consideradas

### Command-line Flags
```go
// -mongo-uri, -db-name, -port flags
```
❌ Requer passing flags em cada run (verboso)
❌ Difficil para managed services (cloud platform config é melhor)

### Configuration Files (JSON, YAML, TOML)
```go
// config.yaml
mongo_uri: "mongodb://..."
```
❌ Arquivo adicional para manage
❌ Difícil para containerized apps (mount config file)
✅ Bo para configurações complexas com muitos aninhamentos

### Configuration Services (Consul, etcd)
❌ Overkill para aplicação simples
✅ Adequado para distributed systems com dynamic config

### Code (constants)
```go
const MongoURI = "mongodb://..."
```
❌ Security issue (secrets committed)
❌ Requer recompilação para mudar
❌ Impossível ter múltiplas instâncias com configs diferentes

## Validation

Application deve fail-fast se config obrigatória está faltando:

```go
func main() {
    requiredEnvVars := []string{"MONGO_URI", "MONGO_DB_NAME"}
    for _, envVar := range requiredEnvVars {
        if os.Getenv(envVar) == "" {
            log.Fatalf("%s environment variable is required", envVar)
        }
    }
}
```

## Consequências

### Positivas
- Fácil deployment em diferentes ambientes
- Secrets fora de version control
- Standard practice em cloud-native applications
- Compatível com 12-factor app methodology

### Negativas
- Typo em nomes de variáveis não é detectado em compile-time
- Requer documentation de quais variáveis são necessárias
- Testes locais requerem setup de environment

## Melhorias Futuras

### Config Struct com Validation
```go
type Config struct {
    MongoURI   string `env:"MONGO_URI,required"`
    MongoDBName string `env:"MONGO_DB_NAME,required"`
    Port        int    `env:"PORT" envDefault:"8080"`
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    if err := env.Parse(cfg); err != nil {
        return nil, err
    }
    return cfg, nil
}
```

### Feature Flags
```go
const FeatureBulkUpdate = os.Getenv("ENABLE_BULK_UPDATE") == "true"
```

## Referências
- [12-Factor App: Config](https://12factor.net/config)
- [The GO Blog: Environment Variables](https://go.googlesource.com/proposal/+/master/design/53606-embedded-config.md)
