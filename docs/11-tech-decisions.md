# Engineering decisions and assumptions

- **2026-10-24 — Go engine:** Go chosen for approachable concurrency, server deployment and tests. Pure rules package; network layer thin. Rust deferred.
- **2026-10-24 — React + Vite + TypeScript:** mobile-first SPA; server-authoritative gameplay. Avoid Next.js/SSR without proven need.
- **2026-10-24 — One repo, two branches:** `main` protected stable; `demo` experimental. Same codebase and independent data/config.
- **2026-10-24 — One Hermes profile:** Rado gets one recognizable Telegram agent. Capabilities are composed through project-local skills, not many confusing bot personas.
- **2026-10-24 — No bots in real matches:** production matchmaking always waits for four distinct humans.
- **2026-10-24 — Rule correctness:** benchmark via golden fixtures + replay explanation; P0 rules must be independently verified before publicly playing.
- **2026-10-24 — Find Skills:** external discovery with source/script review and restricted permissions, not auto-execution of arbitrary community scripts.

Changes to architecture should add a dated entry describing trade-off, rejected options, security impact and reversible migration.
