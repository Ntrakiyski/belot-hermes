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

- [x] Publish feature commit and PR into demo; wait for all CI checks.
- [x] Merge tested PR into demo and migrate existing Coolify application to Git Compose with service domains.
- [x] Verify deployed commit, healthy service, normal HTTPS, UI and disposable DEMO game flow.
- [x] Record final deployment evidence and remaining Hermes/release limitations.

Live replacement loses in-memory tables; preserve the old image/config for rollback. Keep existing DEMO practice setting unless intentionally required otherwise. Do not expose credentials or send Telegram messages.

Publication: PR #5 merged to demo at 12c75ae98dbec57fe09eff1b0eb473b5c55093fc after all four CI jobs passed. Old image retained as belot-hermes-demo:pre-compose-20261107. Native Coolify loadComposeFile loaded the Git manifest; service belot has the two custom HTTPS domains and the existing temporary hostname, internal port 3000. Deployment pqmnucis6w6goi1jsaafclqq queued. Auto-deploy paused only for migration; restore after verified deployment.


### Publication and deployment result

- PR https://github.com/Ntrakiyski/belot-hermes/pull/5 merged into demo at 12c75ae98dbec57fe09eff1b0eb473b5c55093fc. PR checks and post-merge CI https://github.com/Ntrakiyski/belot-hermes/actions/runs/42406468136 passed all four jobs.
- Coolify deployment pqmnucis6w6goi1jsaafclqq finished on that exact commit. Build pack dockercompose, Git manifest /docker-compose.yaml, service belot on internal port 3000, container healthy. Original auto-deploy setting restored to true. Existing DEMO practice setting 1 retained. Old image retained with a rollback tag; config backup is private on the Coolify host.
- https://belot-demo.trakiyski.work/healthz and the same origin with verified TLS return status ok. HTTP redirects to HTTPS. Live API smoke passed authorization rejection, isolated seat projections, eight tricks, next deal and explicit practice. Browser create -> fill bots -> bid -> play passed without error/warning console logs; screenshot /tmp/belot-live-compose.png.
- Optional www.belot-demo.trakiyski.work has no DNS record. No connected DNS provider or valid DNS-write credential was available; do not imply the alias works. Canonical domain works.
- MAIN not changed or released. No live Hermes persona sync/restart, credential change or Telegram message sent. Rado text/voice, unauthorized-sender rejection and repo-scoped PR API access remain unverified. Container replacement reset in-memory guest tables as documented.
- Deployment evidence is also recorded in PR #5. Local tasks/ notes contain post-deployment operational evidence; product code is published and deployed.


# Playable v0.3 update

## Goal
Compare the supplied iOS Belot reference with our game, verify Bulgarian rules separately, and deliver a coherent playable Go/React DEMO version.

## Constraints
- Bulgarian Belot remains authoritative; reference is Croatian, not a rules specification for this product.
- Preserve Go validation, hidden hands, explicit practice opt-in, existing guest room flow and protected MAIN.
- No new accounts/matchmaker infrastructure, hidden bot substitutions or unrelated dependencies.
- Reference is inspected only; license and asset provenance checked before reuse. Prefer original UI and source-backed rule fixes.
- Work from latest demo on codex/playable-v0-3; preserve previous local operational notes.

## Steps
- [x] Compare reference flow and current rules/UI, record confirmed defects and variant differences.
- [x] Define v0.3 API/gameplay contract and acceptance cases before implementation.
- [x] Implement tested rules fixes, paced server transitions and a usable Bulgarian card-table UI.
- [x] Verify full deals/matches, four-seat visibility, invalid actions, reloads and phone/desktop flow.
- [ ] Review, version, publish and deploy the authorized DEMO update after green checks; leave MAIN unchanged.

## Verification
- [x] Independent source comparison and regression fixtures
- [x] Go race tests, Node checks, frontend build and Compose smoke
- [x] Real browser practice and guest room flow; screenshots and console health
- [ ] Actual DEMO commit/UI/API checked if deployed
- [ ] Remaining rule variants/production gaps documented

## Risks
Scoring corrections change the DEMO ruleset. Deployment resets in-memory tables. Current guest rooms do not authenticate distinct human accounts. Source reference is a different national variant, and a playable release must not imply independently certified rules.

## Review
In progress.


### v0.3 verification before publication

- Rule audit and reference review separated Bulgarian rules from Croatian UI ideas. No reference code/assets used. Independent review found and helped fix one-shot declaration selection persistence.
- Go race tests and vet pass; 10,000 deterministic full deals verify exact totals/card/turn invariants and40virtual complete practice matches verify bot/pacing/declaration safety.10Node checks and frontend build pass. Compose builds/health and full HTTP smoke pass.
- Browser390×844: first-card terce selected, next human play succeeds; individual bot states and four-card hold observed; final result stays hidden during collection. Completed8tricks, inspected source-consistent ledger(170+88 cards,0+20 declarations,170/108 raw,17/11 awards), next deal retains17/11 and rotates dealer. Session recovery after local table reset returned useful home error. Screenshot /tmp/belot-v03-mobile.png. Desktop/guest final checks in progress.
- No MAIN changes. Live pre-v0.3 image retained as belot-hermes-demo:pre-v030-20261107 for rollback.

Browser checks completed: guest invite prefilled the table code, a second tab joined without bots, both lobbies updated2/4 and reload restored the same seat.390×844 and1280×900 had no horizontal page overflow or console errors. Screenshots /tmp/belot-v03-mobile.png and /tmp/belot-v03-desktop.png. Belot toggle and rare four-plus-shorter-sequence UI were not encountered in random browser hands; deterministic engine fixtures cover those rules, not their visual interactions.
