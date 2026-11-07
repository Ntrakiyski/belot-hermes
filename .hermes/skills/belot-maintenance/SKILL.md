---
name: belot-maintenance
description: Maintain, debug, test and improve this Bulgarian Belot Go/React game from Rado's Telegram requests; use for gameplay, UI, bugs, current status and DEMO delivery.
---
# Maintain Belot from Telegram

## Orient from the actual checkout

Rado supplies the goal and tests the result on his phone; you own engineering. Read `CONSTITUTION.md`, `AGENTS.md` and current `README.md` from `/workspace/projects/belot`. Check `git status -sb`, `git log -1` and fetch `origin` before starting a branch from `origin/demo`. Preserve uncommitted/untracked work; never reset or overwrite another task. Do not infer installed features from old memories or profile skill copies.

Current code map (inspect before editing):
- `services/game/internal/belot/`: authoritative rules, declarations and source-linked fixtures. `docs/13-playable-v0.3.md` explains the reference comparison and pinned Bulgarian variant; do not import Croatian rules.
- `services/game/cmd/server/tables.go`, `table_flow.go`: seat-private API, bot deadlines, trick collection and score history. Keep the collection gate; do not restore synchronous bot loops.
- `apps/web/src/App.tsx`, `styles.css`: Bulgarian phone UI, session recovery, legal controls, one-shot declarations and explained score. Read `PRODUCT.md`/`DESIGN.md` for UI work.
- `scripts/smoke-game.mjs`: disposable four-seat and paced-practice HTTP flow. Never run against MAIN.

Use `understand-rado`/`product-80-20` for unclear requests, `feature-planning`/`implementation-quality` for changes, `belot-rules-verification` for rules, `review-and-release` for delivery. Load the applicable repo skill automatically; Rado should never need a skill command.

## Prove a change in the sandbox

Run from the repository root, using the installed tools (do not install host packages):

```bash
export PATH="/workspace/home/go123/go/bin:$PATH"
export GOTMPDIR=/workspace/home/gotmp
mkdir -p "$GOTMPDIR"
(cd services/game && CGO_ENABLED=0 go test ./... && go vet ./...)
npm run check
npm --prefix apps/web run build
```

If frontend dependencies are missing, use `npm --prefix apps/web install --no-audit --no-fund` first. The sandbox image has Go 1.19; the existing workspace Go 1.23 toolchain above is required. Its missing C compiler prevents local race tests. Report that distinction: ordinary tests/vet run in the sandbox, while CI’s `go-engine` runs race tests/vet. Never call an unrun race check a local pass.

For a playable local check, start this with the terminal tool's **background mode**, not a foreground blocking command:

```bash
cd /workspace/projects/belot/services/game
export PATH="/workspace/home/go123/go/bin:$PATH" GOTMPDIR=/workspace/home/gotmp CGO_ENABLED=0
PORT=8081 STATIC_DIR=/workspace/projects/belot/apps/web/dist ENABLE_PRACTICE_BOTS=1 go run ./cmd/server
```

Wait for the background server to become ready, then run smoke from the repo root:

```bash
curl --fail --silent --show-error --retry 30 --retry-connrefused --retry-delay 1 --retry-max-time 60 --max-time 2 http://127.0.0.1:8081/healthz &&
BELOT_SMOKE_URL=http://127.0.0.1:8081 npm run smoke:game -- --practice
```

If readiness times out, inspect the background process and its logs before retrying; do not report a gameplay failure from a server that never started.

The sandbox does not need the host Docker socket or Coolify credentials. Existing GitHub `compose-smoke` CI tests the Compose packaging. Use available browser/screenshot tools for a changed UI; verify practice -> bid -> play -> held four-card winner -> score -> next deal, and guest invites/reload when affected. If a browser cannot reach the sandbox-local server, distinguish HTTP/local verification from browser verification on the deployed DEMO. Never claim a screenshot or replay you did not inspect.

For declarations, test first-card announcement **and a later play** so stale selections cannot block the game. Score/card/hidden-hand changes require deterministic regression fixtures and rejected-action checks. The player-visible ledger is implemented; the privileged replay inspector is still a target, not an existing tool.

## Deliver without making Rado an operator

Use the existing configured GitHub integration for this repository. Do not request new credentials because one CLI is absent. Open a feature PR into `demo`, inspect the diff and wait for `go-engine`, `react-vite`, `legacy-policy` and `compose-smoke`. Integrate only within the authorized DEMO scope and repository rules. **A merge into `demo` automatically triggers Coolify.** Do not run Wrangler/Pages or manually deploy to another service.

After merge, verify HTTPS at `https://belot-demo.trakiyski.work`, `/healthz`, the expected version and affected gameplay. Use existing deployment evidence (GitHub/Coolify status if available); health/version alone does not prove the exact deployed commit when versions are reused. Report any missing commit evidence plainly. Retry startup health briefly; if deployment fails, use `admin-ops` and report the failing stage instead of claiming success.

Reply to Rado with what changed, a real DEMO link and 1–2 tap steps, what you verified and the remaining limit. Accept feedback in Telegram and repeat. MAIN promotion still requires Rado's explicit approval naming the tested version and the documented human-review gate; do not silently merge/deploy MAIN.
