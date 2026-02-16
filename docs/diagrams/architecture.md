# Architecture Diagrams - List Manager API

## Clean Architecture Layer Diagram

```mermaid
graph TB
    subgraph "Handlers Layer"
        H1[item.go<br/>HTTP Request/Response]
        H2[middleware.go<br/>CORS, Logging, Errors]
        H3[model.go<br/>API Models]
    end

    subgraph "Service Layer"
        S1[item.go<br/>Business Logic]
        S2[service.go<br/>Service Interfaces]
        S3[parser.go<br/>Data Validation]
    end

    subgraph "Repository Layer"
        R1[repository.go<br/>Repository Interfaces]
        R2[mongodb/repository.go<br/>MongoDB Implementation]
        R3[local/service.go<br/>In-Memory Implementation]
    end

    subgraph "Domain Layer"
        D1[item.go<br/>Domain Entities]
        D2[item_test.go<br/>Domain Tests]
    end

    subgraph "Infrastructure"
        DB[(MongoDB<br/>Primary Database)]
        LOG[(Zap Logger<br/>Structured Logging)]
    end

    H1 --> S1
    H2 --> H1
    S1 --> R1
    R1 --> R2
    R1 --> R3
    R2 --> DB
    R3 --> D1
    S1 --> D1
    S1 --> LOG

    style D1 fill:#e1f5ff
    style R1 fill:#fff3e0
    style S1 fill:#f3e5f5
    style H1 fill:#e8f5e9
    style DB fill:#ffebee
```

## Package Dependency Graph

```mermaid
graph LR
    subgraph "cmd/api"
        A[main.go]
        B[handlers/]
        C[server/]
    end

    subgraph "internal"
        D[service/]
        E[repository/]
        F[domain/]
        G[database/mongodb/]
    end

    A --> C
    C --> B
    B --> D
    D --> E
    D --> F
    E --> F
    E --> G
```

## Data Flow - Create Item Operation

```mermaid
sequenceDiagram
    participant Client
    participant CORS as CORS Middleware
    participant Logging as Logging Middleware
    participant Router as Mux Router
    participant Handler as ItemHandler
    participant Service as ItemService
    participant Repository as MongoRepository
    participant MongoDB as MongoDB

    Client->>CORS: POST /item (JSON)
    Note over CORS: Check Origin, Add Headers
    CORS->>Logging: Pass through
    Note over Logging: Log Request Started
    Logging->>Router: Pass through
    Note over Router: Route to CreateItem
    Router->>Handler: CreateItem(w, r)
    Note over Handler: Parse JSON, Validate
    Handler->>Service: CreateItem(ctx, item)
    Note over Service: Business Logic Validation
    Service->>Repository: Create(ctx, domainItem)
    Repository->>MongoDB: InsertOne(document)
    MongoDB-->>Repository: {id, ...}
    Repository-->>Service: domain.Item
    Service-->>Handler: domain.Item
    Handler-->>Router: JSON Response (201)
    Router-->>Logging: Pass through
    Note over Logging: Log Request Completed
    Logging-->>CORS: Pass through
    Note over CORS: Add CORS Headers
    CORS-->>Client: 201 Created
```

## Component Interaction Diagram

```mermaid
graph TB
    HTTP[HTTP Client]

    subgraph "API Server"
        MUX[Mux Router]

        subgraph "Middlewares"
            CORS[CORS Middleware]
            LOG[Logging Middleware]
            ERR[Error Handling Wrapper]
        end

        subgraph "Handlers"
            IH[ItemHandler]
            HH[HealthHandler]
        end

        subgraph "Services"
            IS[ItemService]
        end

        subgraph "Repositories"
            MR[MongoRepository]
            LR[LocalRepository]
        end
    end

    subgraph "External Dependencies"
        MONGODB[(MongoDB)]
        ZAP[Zap Logger]
    end

    HTTP --> CORS
    CORS --> LOG
    LOG --> MUX
    MUX --> ERR
    ERR --> IH
    ERR --> HH

    IH --> IS
    HH --> IS
    HH --> MR

    IS --> MR
    IS --> LR

    MR --> MONGODB
    IH --> ZAP
    IS --> ZAP

    style MONGODB fill:#ffecb3
    style ZAP fill:#b3e5fc
    style CORS fill:#c8e6c9
    style LOG fill:#c8e6c9
    style ERR fill:#ffccbc
```

## Request Lifecycle Overview

```mermaid
stateDiagram-v2
    [*] --> RequestReceived: HTTP Request
    RequestReceived --> CORSCheck: CORS Middleware
    CORSCheck --> LogRequest: Logging Middleware
    LogRequest --> Routing: Mux Router
    Routing --> ErrorHandling: Wrap Handler
    ErrorHandling --> HandlerExec: Execute Handler
    HandlerExec --> ServiceCall: Service Layer
    ServiceCall --> RepositoryCall: Repository Layer
    RepositoryCall --> DatabaseOp: MongoDB

    DatabaseOp --> ResponseBuild: Success
    RepositoryCall --> ErrorPropagate: Error
    ServiceCall --> ErrorPropagate: Error
    HandlerExec --> ErrorPropagate: Error

    ResponseBuild --> LogResponse: Logging Middleware
    LogResponse --> AddCORSHeaders: CORS Middleware
    AddCORSHeaders --> [*]: HTTP Response

    ErrorPropagate --> ErrorResponse: JSON Error Response
    ErrorResponse --> LogResponse
