# ADR-001: Go como Linguagem Principal

## Status
Aceito

## Contexto
O projeto List Manager API requer uma linguagem que ofereça:
- Alto desempenho para operações CRUD
- Suporte nativo a concorrência para múltiplos clientes
- Tempo de execução rápido sem overhead de JVM
- Forte sistema de tipos e segurança em tempo de compilação
- Ecossistema robusto para desenvolvimento web

## Decisão
Utilizar Go (Golang) 1.24+ como linguagem principal para o projeto.

## Rationale
**Performance:** Go é compilado para código de máquina nativo, eliminando overhead de runtime e fornecendo latência baixa ideal para APIs.

**Concorrência:** Goroutines e channels permitem lidar com milhares de requisições concorrentes com eficiência de memória superior a threads tradicionais.

**Simplicidade:** Sintaxe minimalista com poucas keywords facilita onboarding de novos desenvolvedres e manutenção de código.

**Tooling:** `go mod`, `go test`, `gofmt`, `golint` fornecem ferramentas completas out-of-the-box para desenvolvimento profissional.

**Deployment:** Binário único sem dependências externas simplifica deployment em containers e servidores.

**Ecosistema:** Bibliotecas robustas para HTTP (net/http), gorilla/mux para routing, zap para logging, drivers MongoDB oficiais.

## Consequências

### Positivas
- Binário compilado com performance nativa
- Memory footprint reduzido comparado a Java/Node.js
- Startup time rápido para containers serverless
- Type safety em tempo de compilação
- Ferramentas de profiling integradas (pprof)
- Cross-compilation simplificada

### Negativas
- Curva de aprendizado para desenvolvedores vindos de linguagens OO tradicionais (sem classes/herança)
- Ecossistema menor comparado a JavaScript/Python
- Generics introduzidos recentemente (Go 1.18) podem ter limitações
- Falta de frameworks "batteries-included" como Django/Rails (requer mais decisões arquiteturais)

## Alternativas Consideradas
- **Java/Spring:** Ecossistema maduro mas overhead de JVM maior
- **Node.js:** Excelente para I/O bound mas single-threaded event loop pode ser limitante
- **Python:** Desenvolvimento rápido mas performance inferior para APIs de alta throughput
- **Rust:** Performance excepcional mas curva de aprendizado íngreme e tempo de desenvolvimento maior
