# Task Plan

## Goal
Run Belot locally, verify current DEMO operations, prepare a minimal Docker Compose deployment, and identify what Rado's Telegram Hermes workflow still needs.

## Constraints
- Start from demo; compare main without promoting or changing it. MAIN requires Rado's explicit version approval and human review.
- Preserve the existing Go + React single-image architecture. No database/Redis until implemented.
- Read-only live infrastructure checks. No gateway restart, credentials, push, deployment or public messaging without concrete review and authorization.
- Never print tokens, private player data or hidden cards from live games.
- Working checkout is belot-hermes/ because the supplied workspace already has its own .git directory.

## Steps
- [x] Clone demo, read governing instructions and compare main.
- [x] Run existing tests and the game locally; inspect actual browser behavior.
- [x] Add minimal Compose and align deployment/Hermes documentation with verified topology.
- [x] Fix reproduced practice API/UI response mismatch with a regression test.
- [x] Validate Compose, review diff and record operational gaps and exact migration steps.

## Verification
- [x] Main output checked
- [x] Relevant source/data reviewed
- [x] Tool results verified
- [x] Manual flow checked
- [x] Edge cases checked
- [x] Stakeholder-facing result reviewed
- [x] Risks or limitations documented

## Review
Local preparation complete; live migration awaits operator approval under the supplied publishing/live-change policy.

- Cloned demo into belot-hermes/, compared origin/main (8 files behind demo), worked only on codex/compose-ops. No changes to MAIN, public repository or live infrastructure.
- Added one-service Compose using the existing image build, loopback-only local override, runtime bots off by default, health check and restart policy. Excluded nested generated build files and env files from Docker context.
- Browser reproduced a blank page after filling bots. Failing regression test proved the action returned only {added}; fixed to return the authenticated requester's view. Rechecked create -> bots -> bid -> card play on a 384x853 view with no error/warning console logs. Screenshot /tmp/belot-local-practice.png. Local practice remains available at http://localhost:3000.
- Passed Go race tests, 10 Node tests (including safe persona sync failure/success), React production build, Compose config validation and image build. HTTP smoke passes through eight tricks and next deal, with bots disabled and explicitly enabled. New CI job added but not run on GitHub.
- Corrected Pages/Discord/path/stack contradictions in current operating instructions and added docs/12-coolify-operations.md. Persona sync assembles a complete file before replacing live SOUL; tested only with temporary fixtures.
- Read-only live evidence: Coolify demo container healthy at 5be6781; temporary domain serves API/UI; intended HTTPS returns 503 because running labels only name the temporary domain. Rado Hermes service active, Telegram connected/getMe verified, two allowlisted users/unrestricted off, repo trusted, sandbox Go/Node/npm/Git available. Live SOUL still points to Pages. MAIN rejects force pushes/deletes but requires zero reviews and no CI checks. Independent reviewer found no material issues in final changes.
- Remaining: publish/review/integrate DEMO PR; migrate Coolify and verify TLS/custom domains/deployment commit; preserve remote checkout's untracked work; sync/restart only Rado persona when idle; verify text+voice and unauthorized-user behavior; establish repo-scoped PR API access and strengthen release gates with operator approval. Do not infer these are complete from config or health.
- Risks: restart loses all in-memory tables; this is guest/practice engineering software, not authenticated four-human multiplayer or independently certified Belot. Bot button still appears when server route is disabled. No production release or Rado approval exists for this version.


## Authorized publication and deployment

User authorized publishing the reviewed changes and fixing the existing Coolify Belot project. This authorizes DEMO integration and deployment, not MAIN promotion or unrelated infrastructure changes.

- [ ] Publish feature commit and PR into demo; wait for all CI checks.
- [ ] Merge tested PR into demo and migrate existing Coolify application to Git Compose with service domains.
- [ ] Verify deployed commit, healthy service, normal HTTPS, UI and disposable DEMO game flow.
- [ ] Record final deployment evidence and remaining Hermes/release limitations.

Live replacement loses in-memory tables; preserve the old image/config for rollback. Keep existing DEMO practice setting unless intentionally required otherwise. Do not expose credentials or send Telegram messages.
