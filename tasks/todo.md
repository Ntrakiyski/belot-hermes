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
- [x] Review, version, publish and deploy the authorized DEMO update after green checks; leave MAIN unchanged.

## Verification
- [x] Independent source comparison and regression fixtures
- [x] Go race tests, Node checks, frontend build and Compose smoke
- [x] Real browser practice and guest room flow; screenshots and console health
- [x] Actual DEMO commit/UI/API checked if deployed
- [x] Remaining rule variants/production gaps documented

## Risks
Scoring corrections change the DEMO ruleset. Deployment resets in-memory tables. Current guest rooms do not authenticate distinct human accounts. Source reference is a different national variant, and a playable release must not imply independently certified rules.

## Review
Released v0.3.0 to DEMO through PR #6 after green PR checks and independent review. Confirmed rule omissions and unreadable synchronous bot flow were corrected together; no reference code/assets copied. All user-facing results and limits below are evidence-based.


### v0.3 verification before publication

- Rule audit and reference review separated Bulgarian rules from Croatian UI ideas. No reference code/assets used. Independent review found and helped fix one-shot declaration selection persistence.
- Go race tests and vet pass; 10,000 deterministic full deals verify exact totals/card/turn invariants and40virtual complete practice matches verify bot/pacing/declaration safety.10Node checks and frontend build pass. Compose builds/health and full HTTP smoke pass.
- Browser390×844: first-card terce selected, next human play succeeds; individual bot states and four-card hold observed; final result stays hidden during collection. Completed8tricks, inspected source-consistent ledger(170+88 cards,0+20 declarations,170/108 raw,17/11 awards), next deal retains17/11 and rotates dealer. Session recovery after local table reset returned useful home error. Screenshot /tmp/belot-v03-mobile.png. Desktop/guest final checks in progress.
- No MAIN changes. Live pre-v0.3 image retained as belot-hermes-demo:pre-v030-20261107 for rollback.

Browser checks completed: guest invite prefilled the table code, a second tab joined without bots, both lobbies updated2/4 and reload restored the same seat.390×844 and1280×900 had no horizontal page overflow or console errors. Screenshots /tmp/belot-v03-mobile.png and /tmp/belot-v03-desktop.png. Belot toggle and rare four-plus-shorter-sequence UI were not encountered in random browser hands; deterministic engine fixtures cover those rules, not their visual interactions.


### v0.3 release evidence

- Feature eb752331603ef2565daa00c72dc1f1a2f91e6f38 published and PR https://github.com/Ntrakiyski/belot-hermes/pull/6 merged into demo at2998b7a4433b1151bad72e5a4a4a9477d404f879. All four PR CI jobs passed, including the full Compose gameplay check. Independent reviewer found no remaining blockers after declaration selection fix.
- Coolify webhook deployment959ekwkk6vusw8a8crfoyrrn finished on that exact merge commit. Docker Compose service is healthy; HTTPS health reports version0.3.0 and practice enabled. Previous live image retained as belot-hermes-demo:pre-v030-20261107. MAIN unchanged.
- Live full-deal HTTP smoke passed: private seats/declaration options, invalid action rejection, held tricks, next deal and paced practice completion. Default bots-off local Compose smoke also passed.
- Live browser: new home/version, immediate explicit practice, paced bidding, all-trump hand, explicit belot checkbox on first Q-spades awarded exactly20; later K-spades played successfully without duplicate belot. Completed four-card trick showed winner Elena/+17, controls disabled while collecting, no console warnings/errors. Screenshot /tmp/belot-v030-live.png. Phone screenshot /tmp/belot-v03-mobile.png and desktop /tmp/belot-v03-desktop.png reflect the same tested bundle.
- Ordinary first-card declaration, full visible8trick deal/result/next deal, guest invite join and reload verified locally.40complete practice matches were simulated under virtual deadlines and10,000deal invariants pass. Rare four-plus-shorter-sequence UI interaction was not manually encountered; engine fixtures cover valid/invalid combinations. No every-house-rule certification, authenticated accounts/public matchmaking, persistence or privileged replay inspector is claimed.
- Release notes/evidence added to PR #6. Operational notes after publication remain local task records. No live Hermes persona change or Telegram messages in this release.


# Hermes Telegram operations verification

## Goal
Make the dedicated Rado Hermes profile capable of maintaining and improving v0.3 from Telegram, using its existing GitHub access and merge-triggered Coolify deployment.

## Constraints
- User confirms GitHub access; retain credentials and access boundaries, do not create new auth or treat absent gh as proof of no access.
- Only the dedicated Belot profile/workspace may change; preserve remote working branches and untracked work.
- Keep Telegram numeric allowlist, sandbox isolation and MAIN approval boundary intact.
- Inspect only necessary sanitized configuration/state; never print tokens, conversation history or player data.
- Live instructions need to use current repo/skills and automatic Coolify deployment, not Pages/manual hosting.

## Steps
- [ ] Inspect live profile, skills discovery, sandbox tools, Telegram text route and gateway health.
- [ ] Correct minimal skills/persona/test/release guidance for the actual v0.3 workflow.
- [ ] Apply reviewed instructions to the dedicated live profile/workspace safely; restart only if needed and idle.
- [ ] Test actual agent tool/skill behavior and Telegram request/response path; distinguish simulated checks from real delivery.
- [ ] Publish reviewed changes, verify deployment continuity and record evidence/remaining blocks.

## Verification
- [ ] Profile gateway and Telegram connection/allowlist
- [ ] Current repo, trusted skills and runtime tools
- [ ] Go/Node/frontend tests runnable inside sandbox
- [ ] Agent follows request -> plan -> test -> PR -> CI -> merged DEMO -> live verification
- [ ] Text actual channel test or explicit limitation

## Review
In progress. Initial direct SSH route100.123.30.5 timed out; investigate connectivity before modifying live state.


## Steering: no voice support; investigate disk usage

User explicitly removed voice support from scope and prioritized why the VPS has no space. No voice configuration was activated. Remove only the disposable STT environment/fixtures created during this audit. Continue the previously authorized Hermes text/skills maintenance after the disk investigation; do not install or configure transcription.

- [x] Identify disk/inode usage, Docker cache/images/volumes, logs and deleted-open files using metadata only.
- [x] Explain the primary consumers and growth mechanism; distinguish old buildup from this audit's additions.
- [x] Verify reclaimed capacity and running services; recommend only targeted further changes.

### Disk investigation evidence

- Root filesystem is 75GiB, 70GiB used, 2.3GiB available (97%); inodes 34%. It was already full before this audit's temporary speech environment.
- Largest physical consumers: containerd 18GiB (snapshots ~14.6GiB, blobs ~3.4GiB), root user tools/caches ~21.8GiB, other project home directories ~11.4GiB. Docker volumes ~2.1GiB; logs ~0.5GiB. No material deleted-open file leak. Docker image/build figures overlap physical storage and must not be summed as independent consumption.
- Deleted only unused build cache older than 7 days then older than 24 hours: Docker reported 2.252GB + 984.7MB reclaimed. Removed the isolated 474MiB speech environment/fixtures created during this audit; no transcription config was enabled.
- BuildKit uses permissive default retention (up to ~55.9GiB cache, one policy retaining 60 days); no scheduled prune was found. Old editor/browser/package caches and project checkouts compound the pressure. No databases, volumes, project checkouts, conversations or running containers deleted.
- Further image/cache pruning needs a targeted retention decision; unreferenced images may be rollback images and unused volumes may hold data. Remaining 2.3GiB headroom is insufficient for comfortable builds.

### Hermes verification evidence

- Replaced stale copied MVP/Pages knowledge with two managed profile bridges and a repo-owned maintenance skill. Backed up the live persona/config/old skills and archived the competing untracked working-session skill outside discovery. Preserved untracked frontend lockfile and old feature branch.
- Native Hermes config loader confirmed drop_pending_on_cold_boot=false before an idle-only dedicated gateway restart. Service active and Telegram connected afterward; allowlist and credentials unchanged.
- Actual configured Hermes headless probe (source=tool, no messages/edits) automatically read maintenance/release/profile skills, found v0.3.0 and the current checkout, described correct tests/automatic Coolify path/MAIN gate, and correctly flagged CI race guidance pending publication. This proves tool/skill routing, not Telegram delivery or PR creation.
- Actual sandbox: Go 1.23 ordinary tests and vet, all 10 Node checks, React TypeScript/Vite build passed. Local updated profile-sync regression and Node checks passed. Race checks run in CI, not the compiler-free sandbox.
- Text-only steering applied to current charter, persona, user guide, setup and intent skill. Voice is excluded. Real Telegram user request/reply and unauthorized sender rejection remain unverified; no chat history read or Telegram messages sent by this audit.
