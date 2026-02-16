# Deployment Diagrams - List Manager API

## Docker Compose Architecture

```mermaid
graph TB
    subgraph "Docker Network: list-manager-api_default"
        subgraph "Container: list-manager-api"
            API[Go Application<br/>Port 8085]
        end

        subgraph "Container: list-manager-mongodb"
            MONGO[(MongoDB 8.0<br/>Port 27017)]
            VOL[mongodb_data<br/>Named Volume]
        end

        subgraph "Container: list-manager-mongo-express"
            EXPRESS[Mongo Express UI<br/>Port 8081]
        end
    end

    subgraph "Host Machine"
        CLI[Docker CLI]
        COMPOSE[docker-compose.yml]
    end

    subgraph "External Access"
        USER[Developer/User]
        BROWSER[Web Browser]
    end

    USER -->|:8085| API
    BROWSER -->|:8081| EXPRESS

    API -->|mongodb://mongodb:27017| MONGO
    EXPRESS -->|mongodb://mongodb:27017| MONGO
    MONGO --> VOL

    COMPOSE --> CLI
    CLI --> API
    CLI --> MONGO
    CLI --> EXPRESS

    style MONGO fill:#4db6ac
    style EXPRESS fill:#f48fb1
    style API fill:#4dd0e1
    style VOL fill:#fff59d
```

## Local Development Environment

```mermaid
flowchart LR
    subgraph "Developer Machine"
        DEV[Developer]
        IDE[IDE/Editor<br/>VS Code, GoLand]
        TERM[Terminal]

        subgraph "Project Directory"
            SRC[/cmd/api<br/>/internal<br/>/docs]
            CFG[docker-compose.yml<br/>Makefile<br/>.env]
            GIT[.git/]
        end
    end

    subgraph "Docker Daemon"
        subgraph "Containers"
            C1[list-manager-api<br/>:8085]
            C2[mongodb<br/>:27017]
            C3[mongo-express<br/>:8081]
        end
        subgraph "Volumes"
            V1[mongodb_data]
        end
        subgraph "Networks"
            N1[list-manager-api_default]
        end
    end

    DEV --> IDE
    DEV --> TERM
    IDE --> SRC
    TERM --> SRC
    TERM -->|docker compose up| C1
    TERM -->|docker compose logs| C1
    TERM -->|docker compose down| C1

    C1 --> N1
    C2 --> N1
    C3 --> N1
    C2 --> V1

    SRC --> GIT

    style DEV fill:#e3f2fd
    style SRC fill:#e8f5e9
    style C1 fill:#4dd0e1
    style C2 fill:#4db6ac
    style C3 fill:#f48fb1
```

## Production Deployment (Render)

```mermaid
graph TB
    subgraph "Render Cloud Platform"
        subgraph "Web Service"
            WS[List Manager API<br/>Web Service<br/>RAM: 512MB<br/>CPU: 0.1]
        end

        subgraph "Managed Database"
            RDB[(Render MongoDB<br/>Managed Instance<br/>Auto-backups<br/>High Availability)]
        end

        subgraph "Environment Variables"
            ENV[PORT=8085<br/>MONGO_URI<br/>MONGO_DB_NAME<br/>...]
        end
    end

    subgraph "External World"
        CLIENTS[API Clients<br/>Web Apps, Mobile, PWA]
        INTERNET[Public Internet<br/>HTTPS]
        MONITOR[Monitoring<br/>Logs, Metrics]
    end

    CLIENTS -->|HTTPS| INTERNET
    INTERNET -->|Custom Domain| WS
    WS --> ENV
    WS --> RDB
    WS --> MONITOR

    WS -->|Deploy on git push| GIT[GitHub Repository]
    WS -->|Auto-deploy| BRANCH[main branch]

    style WS fill:#4dd0e1
    style RDB fill:#4db6ac
    style ENV fill:#fff59d
    style INTERNET fill:#e1f5fe
```

## CI/CD Pipeline

```mermaid
flowchart LR
    subgraph "Git Repository"
        PUSH[git push]
        PR[Pull Request]
        MAIN[main branch]
    end

    subgraph "GitHub Actions CI"
        CHECKOUT[Checkout Code]
        DEPS[go mod download]
        TEST[go test ./...<br/>Unit Tests]
        LINT[golangci-lint<br/>Linting]
        BUILD[go build<br/>Compile Check]
        ARTIFACT[Build Artifact]
    end

    subgraph "Render Deployment"
        DEPLOY[Deploy Hook<br/>auto-deploy on main]
        ROLLBACK[Rollback<br/>if health check fails]
    end

    subgraph "Production"
        RENDER_APP[Render Web Service<br/>Live Application]
        RENDER_DB[Render MongoDB<br/>Production Database]
    end

    PUSH --> PR
    PR --> CHECKOUT
    CHECKOUT --> DEPS
    DEPS --> TEST
    TEST -->|Pass| LINT
    TEST -->|Fail| FAILED[PR Check Failed]
    LINT -->|Pass| BUILD
    LINT -->|Fail| FAILED
    BUILD -->|Pass| ARTIFACT
    BUILD -->|Fail| FAILED

    ARTIFACT -->|PR merged to main| MAIN
    MAIN --> DEPLOY
    DEPLOY --> RENDER_APP
    RENDER_APP --> RENDER_DB

    RENDER_APP -->|Health Check OK| DONE[Deployment Complete]
    RENDER_APP -->|Health Check Fail| ROLLBACK

    style TEST fill:#c8e6c9
    style LINT fill:#fff59d
    style BUILD fill:#4dd0e1
    style FAILED fill:#ef5350
    style RENDER_APP fill:#26a69a
    style DONE fill:#66bb6a
```

## Network Topology (Production)

```mermaid
flowchart TB
    subgraph "Internet"
        USER[End Users]
        CDN[CDN / Load Balancer<br/>Render Managed]
    end

    subgraph "Render Region: Oregon"
        subgraph "Private Network"
            API[list-manager-api<br/>Container 1]
            API2[list-manager-api<br/>Container 2<br/>(Auto-scale)]
            MONGO[(MongoDB<br/>Primary)]
        end
    end

    subgraph "Monitoring & Observability"
        LOGS[Render Logs<br/>Real-time streaming]
        METRICS[Metrics Dashboard<br/>CPU, Memory, Requests]
        ALERTS[Alerts<br/>Response time, Errors]
    end

    USER -->|HTTPS| CDN
    CDN -->|Round Robin| API
    CDN --> API2
    API --> MONGO
    API2 --> MONGO

    API --> LOGS
    API2 --> LOGS
    API --> METRICS
    API2 --> METRICS
    METRICS --> ALERTS

    style MONGO fill:#4db6ac
    style API fill:#4dd0e1
    style API2 fill:#26c6da
    style CDN fill:#fff59d
```

## Infrastructure Components

```mermaid
graph TB
    subgraph "Application Layer"
        APP[Go Application<br/>Binary]
        HANDLER[Handlers<br/>HTTP Layer]
        SERVICE[Services<br/>Business Logic]
        REPO[Repositories<br/>Data Access]
    end

    subgraph "Infrastructure Layer"
        ZAP[Zap Logger<br/>Structured Logs]
        MUX[Mux Router<br/>HTTP Routing]
        MONGO_DRV[MongoDB Driver<br/>Client]
        MONGO_POOL[Connection Pool<br/>Max Pool Size: 100]
    end

    subgraph "External Services"
        MONGO_EXT[(MongoDB Server)]
        RENDER_LOGS[Render Log Stream]
    end

    APP --> HANDLER
    HANDLER --> SERVICE
    SERVICE --> REPO
    REPO --> MONGO_DRV
    MONGO_DRV --> MONGO_POOL
    MONGO_POOL --> MONGO_EXT

    APP --> ZAP
    ZAP --> RENDER_LOGS

    HANDLER --> MUX

    style APP fill:#4dd0e1
    style MONGO_EXT fill:#4db6ac
    style RENDER_LOGS fill:#ffca28
```

## Development vs Production Comparison

```mermaid
graph LR
    subgraph "Local Development"
        L1[Go run cmd/api/main.go<br/>Local Build]
        L2[docker-compose up<br/>MongoDB Container]
        L3[localhost:8085<br/>Unsecured HTTP]
        L4[Local .env file<br/>Clear credentials]
    end

    subgraph "Production"
        P1[Render Build<br/>Optimized Binary]
        P2[Render MongoDB<br/>Managed Service]
        P3[Custom Domain<br/>HTTPS with SSL]
        P4[Render Environment<br/>Encrypted Secrets]
    end

    subgraph "Shared"
        S1[Clean Architecture<br/>Same Code]
        S2[Environment Config<br/>MONGO_URI, PORT]
        S3[Zap Logging<br/>Structured Output]
        S4[Health Check<br/>/healthz endpoint]
    end

    L1 --> S1
    L2 --> S2
    L3 --> S4
    L4 --> S2

    P1 --> S1
    P2 --> S2
    P3 --> S4
    P4 --> S2

    style L1 fill:#e3f2fd
    style L2 fill:#e3f2fd
    style P1 fill:#a5d6a7
    style P2 fill:#a5d6a7
    style S1 fill:#fff59d
    style S2 fill:#fff59d
```
