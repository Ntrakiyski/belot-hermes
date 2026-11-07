# DEMO and MAIN release contract

| | DEMO | MAIN |
|---|---|---|
| Git branch | `demo` | `main` |
| Audience | Rado, testers, Hermes | Real players |
| Database, accounts, secrets | Sandbox only | Production only |
| Debug inspector | Authorized, server-gated | Physically excluded/disabled server-side |
| Release | Agent-generated PR, checks, smoke tests | Explicit Rado approval + green promotion PR |
| Failure | Roll back Demo and share diagnostic | Roll back version and alert owner |

Use one codebase; do not maintain copy-pasted `main/` and `demo/` source trees. DEMO has a Coolify Git application; MAIN hosting is not verified. See [the operations runbook](12-coolify-operations.md) for deployment evidence, Compose migration and release gaps.

Minimum production checks: isolated environment variables, TLS, access control, no hidden hand leaks, database compatibility for in-progress games, health/readiness, verified rollback, error monitoring and documented owner approval. Hermes can prepare, never silently approve production.
