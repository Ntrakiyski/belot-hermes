# Belot — Hermes-managed Bulgarian multiplayer card game

**Engineering seed v0.2.0 (2026-10-24) · Public repo · NOT yet online multiplayer.**

## Product contract

- Home: **Start Random Game** (four real players), then **Create a Room** (four invited real players).
- Random queue accepts one solo or a duo of friends; prefer duo-v-duo and four-solo matches, with 5–15s fallback to mixed. Three-second countdown once four connected accounts are matched.
- Private rooms: unique usernames, friends, invitations, four human seats and host-managed teams.
- Profiles: tap a player, add friend, block/report. No hidden hand exposure.
- Hermes works from Discord into the `demo` Git branch and its deployment; production `main` requires explicit human review.

## Locked technology stack

| Component | Choice |
|---|---|
| Game rules/server | **Go 1.23+** (server authoritative, WebSocket to be implemented) |
| Frontend | **React + Vite + TypeScript** (mobile first) |
| Persistence target | PostgreSQL; Redis when horizontal scale requires it |
| Quality | Go tests, React build, Node matcher-policy tests via GitHub Actions |
| Environments | `demo` and `main` **Git branches**, separate deployment settings/DB; not two frontend codebases |

## What's implemented today

- Working **Go 32-card deck model**, cryptographically secure shuffle, exact four-hand deal with validation and unit tests; `/healthz` HTTP endpoint reports server process only. The bidding-stage five-plus-three deal and legal move/scoring rules are **not** yet implemented.
- React + Vite + TypeScript home screen scaffold with the two intended primary buttons, **disabled** until real authentication/matchmaking/rooms exist. Connects to Go `/healthz` in local dev.
- Four authoritative product/architecture/agent/workflow documents and three Hermes SKILL.md starter files.
- Standalone JavaScript matchmaking composition-policy example with automated tests.
- Separate legacy single-human/three-bot HTML reference under `prototype/`, `main/`, `demo/`. These are **local snapshots, not deploys**.
- CI definition; GitHub setup/rules and live deployment are separate operations.

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

## Next milestones

1. Implement the **Go rules engine** (deal five-plus-three, bids/contracts, legal moves, declarations, scoring and validated golden fixtures).
2. Implement authenticated human identities, Go WebSocket event broker and server-owned matches.
3. Connect React Start Random Game / Create Room to real Go endpoints, friends and invites; no bot fill.
4. Provision **separate DEMO and MAIN** hosting/data, CI protections and Discord Hermes agent workflow.
5. Release only after rule, multiplayer, security and reconnection acceptance criteria pass.
