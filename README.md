# Belot — Hermes-managed Bulgarian multiplayer game

**Engineering seed v0.2.0 · Go backend + React/Vite/TypeScript frontend.**

This repository is the code and specification foundation for a four-human-player Bulgarian Belot game. It is **not yet a deployed online multiplayer application**.

## Product
- **Start Random Game**: solo queue or two-person friend party. Four real human players; never insert bots.
- **Create a Room**: four invited humans via username, friends or invite link.
- Player profiles: add friend, report, block.
- Matchmaking prefers balanced solo/duo teams, relaxes composition after 5–15 seconds, then starts a server-controlled three-second countdown.

## Technology
- `services/game`: **Go** authoritative game engine and HTTP/WebSocket server (WebSocket gameplay not yet implemented).
- `apps/web`: **React + Vite + TypeScript** mobile-first user interface.
- PostgreSQL target persistence; Redis optional as concurrency requires it.
- `demo` branch for Hermes-driven testing; `main` branch for human-approved production changes.

## Local development
```bash
cd services/game && go test ./...
cd services/game && go run ./cmd/server
cd apps/web && npm install && npm run dev
```

The Go endpoint `/healthz` reports server health only, **not multiplayer readiness**. React buttons remain disabled until the actual game flows exist.

## Documentation
- [Gameplay and matchmaking](docs/01-game-product.md)
- [Architecture and releases](docs/02-architecture-release.md)
- [Hermes + Discord](docs/03-hermes-discord.md)
- [Personas, workflows and acceptance](docs/04-workflows-acceptance.md)

Agent instructions: [AGENTS.md](AGENTS.md); skills: `skills/`.

This is a **public** GitHub repository: never commit credentials, transcripts containing personal data, private cards or live database dumps.
