# Request Flow Diagrams - List Manager API

## Middleware Chain Flow

```mermaid
flowchart TD
    A[Incoming HTTP Request] --> B{CORS Middleware}
    B -->|Preflight OPTIONS| C[Return 200 OK]
    B -->|Regular Request| D{Logging Middleware}
    D --> E[Log: Method, Path, RemoteAddr]
    E --> F[Mux Router]
    F --> G{Route Matching}
    G -->|/healthz| H[Health Handler + Error Wrapper]
    G -->|POST /item| I[Create Handler + Error Wrapper]
    G -->|GET /item| J[Get Handler + Error Wrapper]
    G -->|PUT /item| K[Update Handler + Error Wrapper]
    G -->|DELETE /item| L[Delete Handler + Error Wrapper]
    G -->|GET /items| M[List Handler + Error Wrapper]
    G -->|PUT /items/active| N[Bulk Update Handler + Error Wrapper]
    G -->|/_app/version.json| O[Version Handler - No Wrapper]

    H --> P[Handler Execution]
    I --> P
    J --> P
    K --> P
    L --> P
    M --> P
    N --> P
    O --> P

    P --> Q[Return Response]
    Q --> R{Logging Middleware - Response}
    R --> S[Log: Status, Duration]
    S --> T{CORS Middleware - Response}
    T --> U[Add CORS Headers]
    U --> V[HTTP Response to Client]

    style B fill:#c8e6c9
    style D fill:#b3e5fc
    style F fill:#fff9c4
    style P fill:#e1bee7
    style V fill:#ffcdd2
```

## Create Item Flow (Detailed)

```mermaid
flowchart TD
    A[Client: POST /item<br/>JSON body] --> B[CORS Middleware<br/>Check origin, add headers]
    B --> C[Logging Middleware<br/>Log incoming request]
    C --> D[Mux Router<br/>Match POST /item]
    D --> E[Error Handling Wrapper<br/>Setup panic recovery]
    E --> F[ItemHandler.CreateItem]

    subgraph "Handler Layer"
        F --> G[Parse JSON body<br/>Validate fields]
        G --> H{Valid?}
        H -->|No| I[Return 400 Bad Request<br/>Invalid input]
        H -->|Yes| J[Call Service.CreateItem]
    end

    subgraph "Service Layer"
        J --> K[ItemService.CreateItem]
        K --> L[Business Logic<br/>Validate name not empty]
        L --> M{Valid?}
        M -->|No| N[Return Service Error]
        M -->|Yes| O[Call Repository.Create]
    end

    subgraph "Repository Layer"
        O --> P[MongoRepository.Create]
        P --> Q[BSON tag mapping<br/>domain → mongo]
        Q --> R[Create MongoDB session]
        R --> S[InsertOne operation]
        S --> T{Success?}
        T -->|No| U[Return Repository Error]
        T -->|Yes| V[Return created Item]
    end

    V --> W[Service returns domain.Item]
    W --> X[Handler maps to API model<br/>domain → handlers]
    X --> Y[Write JSON 201 Created<br/>with location header]
    Y --> AA[Logging Middleware<br/>Log completed request]
    AA --> AB[CORS Middleware<br/>Add CORS headers]
    AB --> AC[Client receives 201 Created]

    I --> AC
    N --> AC
    U --> AC

    style A fill:#e3f2fd
    style K fill:#f3e5f5
    style P fill:#fff3e0
    style V fill:#c8e6c9
    style AC fill:#a5d6a7
```

## Error Flow Path

```mermaid
flowchart TD
    A[Handler Execution] --> B{Error Occurs?}
    B -->|No Error| C[Normal Response Path]
    B -->|Error Returned| D[Error Handling Wrapper]
    B -->|Panic| E[Recover Middleware]

    D --> F{Error Type?}
    F -->|Service Error| G[Map to HTTP Status<br/>400, 404, etc.]
    F -->|Repository Error| H[Map to HTTP Status<br/>500 or 404]
    F -->|Unknown Error| I[Return 500 Internal Error]

    E --> J[Log Stack Trace]
    J --> I

    G --> K[JSON Error Response<br/>{message: "..."}]
    H --> K
    I --> K

    K --> L[Logging Middleware<br/>Log error response]
    L --> M[CORS Middleware<br/>Add headers]
    M --> N[Client receives error response]

    C --> N

    style D fill:#ffccbc
    style E fill:#ef9a9a
    style I fill:#ef5350
    style K fill:#ffca28
```

## Health Check Flow

```mermaid
flowchart TD
    A[Client: GET /healthz] --> B[CORS Middleware]
    B --> C[Logging Middleware]
    C --> D[Mux Router → HealthCheck]
    D --> E[Error Handling Wrapper]
    E --> F[HealthHandler.HealthCheck]

    subgraph "Health Check Logic"
        F --> G[Check Server Status<br/>Always up]
        G --> H[Ping MongoDB<br/>Ping command]
        H --> I{MongoDB OK?}
        I -->|Yes| J[Set database: connected<br/>Set status: up]
        I -->|No| K[Set database: disconnected<br/>Set status: degraded]
        J --> L[Build HealthCheckResponse<br/>with timestamp and checks]
        K --> L
    end

    L --> M[Return JSON 200 OK<br/>Health check data]
    M --> N[Logging Middleware<br/>Log health check completion]
    N --> O[CORS Middleware<br/>Add headers]
    O --> P[Client receives health status]

    style H fill:#fff3e0
    style I fill:#ffe082
    style L fill:#c8e6c9
    style P fill:#a5d6a7
```

## Bulk Update Flow

```mermaid
flowchart TD
    A[Client: PUT /items/active<br/>{active: true}] --> B[CORS]
    B --> C[Logging]
    C --> D[Mux Router → BulkUpdateActive]
    D --> E[Error Wrapper]
    E --> F[Handler.BulkUpdateActive]

    F --> G[Parse JSON<br/>Extract active bool]
    G --> H{Valid?}
    H -->|No| I[Return 400]
    H -->|Yes| J[Service.BulkUpdateActive]

    J --> K[Repository.BulkUpdateActive<br/>with transaction]
    K --> L[Start MongoDB Session]
    L --> M[WithTransaction]
    M --> N[UpdateMany {active: <new_value>}]
    N --> O[Get MatchedCount & ModifiedCount]
    O --> P[Return counts to Service]

    P --> Q[Service returns counts]
    Q --> R[Handler maps to BulkActiveResponse]
    R --> S[Return JSON 200 OK<br/>{matchedCount, modifiedCount}]
    S --> T[Logging - Log completion]
    T --> U[CORS - Add headers]
    U --> V[Client receives 200 with counts]

    I --> V

    style K fill:#fff3e0
    style M fill:#ce93d8
    style S fill:#c8e6c9
```

## Complete Request/Response Lifecycle

```mermaid
sequenceDiagram
    participant C as Client
    participant CORS as CORS Middleware
    participant LOG as Logging Middleware
    participant R as Mux Router
    participant ERR as Error Wrapper
    participant H as Handler
    participant S as Service
    participant REPO as Repository
    participant DB as MongoDB

    Note over C,DB: Forward Path
    C->>CORS: POST /item
    CORS->>LOG: Request with headers
    LOG->>R: Request logged
    R->>ERR: Route matched
    ERR->>H: Wrapped handler call
    H->>S: Service call
    S->>REPO: Repository call
    REPO->>DB: Database operation
    DB-->>REPO: Result
    REPO-->>S: Domain entity
    S-->>H: Service result
    H-->>ERR: Response data

    Note over C,DB: Reverse Path
    ERR-->>LOG: Response to log
    LOG-->>CORS: Response with logs
    CORS-->>C: Final response

    Note over C,CORS: Headers added: Access-Control-Allow-Origin, etc.
    Note over LOG,LOG: Logged: method, path, status, duration
