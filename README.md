# Belot — Hermes-managed Bulgarian multiplayer card game

**v0.3.0 · Bulgarian Belot guest tables and explicit practice · Public repo · Authenticated multiplayer remains a separate milestone.**

## Product contract

- Home: **Start Random Game** (four real players), then **Create a Room** (four invited real players).
- Random queue accepts one solo or a duo of friends; prefer duo-v-duo and four-solo matches, with 5–15s fallback to mixed. Three-second countdown once four connected accounts are matched.
- Private rooms: unique usernames, friends, invitations, four human seats and host-managed teams.
- Profiles: tap a player, add friend, block/report. No hidden hand exposure.
- Rado talks to Hermes from Telegram; Hermes works into the `demo` Git branch and its deployment; production `main` requires explicit human review.

## Locked technology stack

| Component | Choice |
|---|---|
| Game rules/server | **Go 1.23+** (server authoritative, WebSocket to be implemented) |
| Frontend | **React + Vite + TypeScript** (mobile first) |
| Persistence target | PostgreSQL; Redis when horizontal scale requires it |
| Quality | Go tests, React build, Node matcher-policy tests via GitHub Actions |
| Environments | `demo` and `main` **Git branches**, separate deployment settings/DB; not two frontend codebases |

## What's implemented today

- Go-authoritative Bulgarian ruleset `bg-0.3.0`: 3+2+3 dealing, six contracts, contra/recontra, legal moves, first-card declarations, belot, hanging carry, contract-aware rounding and explained scoring to 151. Reference comparison, exact fixtures and 10,000 seeded deal invariants are documented in [the v0.3 review](docs/13-playable-v0.3.md); this is not certification of every house variant.
- Bulgarian phone/desktop table with one-click explicit practice, friend-room invite links, paced bot turns, visible four-card trick/winner, legal-action controls, declarations, score ledger and match history. Tables remain in memory and disappear on restart. Accounts, public matchmaking, friends, reconnect/forfeit policy and a privileged replay inspector remain unimplemented.
- Optional **explicit practice bots**, enabled server-side only with `ENABLE_PRACTICE_BOTS=1` on DEMO. Default is off; never enable on MAIN. The UI reads the runtime capability and disables practice when unavailable; guest friend tables never silently fill with bots.
- One Docker image serves React and Go on port 3000; Docker Compose reuses that image build. No database or Redis is used yet.
- Go race tests, deterministic rules/declaration fixtures, complete-match timing checks, React build, legacy Node checks and a Compose HTTP smoke through a full paced practice deal. CI configuration is not evidence that a particular run passed.
- Dedicated Telegram Hermes profile documented under `ops/hermes/`; see the [operations runbook](docs/12-coolify-operations.md) for verified state, gaps and migration.

## Local development

```bash
# Terminal 1: Go HTTP health server
cd services/game && go run ./cmd/server

# Terminal 2: Vite app (requires npm access for first install)
cd apps/web && npm install && npm run dev
# http://localhost:5173

# Go engine tests
cd services/game && go test ./...

# Legacy JavaScript matcher policy and snapshot tests (repo root)
npm run check
```

No production credentials or player data may be committed to this **public** repository.

## Repo layout

```text
apps/web/                 React + Vite frontend
services/game/            Go service and server-owned Belot deck package
packages/matchmaker/      Legacy JS matchmaking decision utility (policy prototype)
prototype/                Legacy bot-play reference (NOT the target multiplayer UI)
main/ and demo/           Legacy identical HTML snapshots (NOT real deployments)
docs/                     Four Markdown pillars and decisions D01–D09
skills/                   Hermes skill instructions (delivery, gameplay QA, admin)
AGENTS.md                 Agent work policy and security constraints
.github/workflows/ci.yml  Go + React + Node validation
```

## Docker Compose

```bash
# Local game, bots disabled; UI and API on http://localhost:3000
docker compose -f docker-compose.yaml -f docker-compose.local.yaml up --build -d --wait
npm run smoke:game

# Optional local practice, recreates the container and loses existing tables
ENABLE_PRACTICE_BOTS=1 docker compose -f docker-compose.yaml -f docker-compose.local.yaml up -d --wait
npm run smoke:game -- --practice
```

For Coolify use **only `docker-compose.yaml`**, assign domains to service `belot` at internal port 3000, and follow [the migration and rollback procedure](docs/12-coolify-operations.md). The local override binds only to loopback. Compose orchestrates the service; the existing Dockerfile still builds it.

## Next milestones

1. Independently verify rules with golden fixtures and a secured DEMO inspector.
2. Implement authenticated identities, real multiplayer connectivity and reconnect handling.
3. Build public matchmaking, private room permissions, friends and invitations.
4. Verify Rado's Telegram text/voice → feature → tests → DEMO workflow and release gates.
5. Promote a tested version to MAIN only after explicit Rado approval and human review.
