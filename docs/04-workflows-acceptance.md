# Pillar 04 — Personas, workflows, skills and release acceptance

**Revision:** 2026-10-24. Matches the owner's real-human-first product direction. The local 1-human/3-bot prototype is exclusively an engineering reference, NOT the planned public user journey.

## Persona 1 — Player: fastest possible real-human Belot

**Random solo:** `sign in -> Start Random Game -> queue -> 4 humans found -> teams + 3s countdown -> 8 tricks/deal -> score -> 151-point match result -> rematch/friends`. System handles D02/D03 5–15-second fallback, D05 countdown reconnect, D08 in-match reconnect. Never add bots.

**Random duo:** `profile/friends -> invite one friend -> accept -> atomic duo ticket -> server finds opposing humans -> opposite-seat partners -> countdown -> match`. Duo never split or silently downgraded.

**Private:** `Create a Room -> invite via username/friends/link -> four human guests join -> host adjusts seats -> Start -> countdown -> match`. Host transfer and room expiry per D09.

**Social:** tap player avatar while playing -> open public profile -> Add Friend -> recipient accepts/declines -> permitted invitations. Opening profile must not pause server authority, reveal private cards or compromise turn integrity.

**Success signals:** successful match start rate, queue duration distributions, reconnection success, legal-action errors, match completion rates, user-reported scoring defects.

## Persona 2 — Administrator: fair, safe operations

**Example:** player disputes score. `Admin MFA -> search match ID -> review redacted ordered events and pinned ruleset -> reproduce with golden fixture -> compare ledger -> mark resolved or open issue -> deploy correction on DEMO -> independently verify -> owner-approved promotion`. No ad-hoc score rewrite.

**Example:** sudden queue delays. `Review queue count and ticket age -> check disconnected ghosts / capacity / WebSocket health -> inspect composition D02/D03 behavior -> remediate without inventing bots -> record actions`. Can suspend misbehaving account and moderate reports using RBAC/audit.

**Success signals:** mean time to detect/recover, queue SLA, scoring divergence rate, audit completeness, release rollback success.

## Persona 3 — Owner and Hermes developer: voice-to-demo-to-main

`Rado Telegram message -> Hermes task issue -> isolated branch -> implement -> tests incl. gameplay QA -> review -> demo PR -> CI -> deploy DEMO -> report link/test proof -> owner accepts -> protected main promotion PR and human review -> deploy MAIN -> smoke/rollback-ready`.

**Contract:** Hermes may autonomously implement/test DEMO within prescribed permissions, but **cannot self-authorize MAIN**. No link must be reported as deployed before verification.

## Acceptance checkpoints

| Gate | Pass condition |
|---|---|
| Product | 4 humans, solo/duo, private rooms, friends + profiles; home two primary buttons |
| Gameplay | Complete server-authoritative deals with verified legal moves/score and repeatable ruleset fixture suite |
| Reliability | D05 countdown replacement, D06 queue cancellation, D08 reconnect/forfeit tested with dropped sockets |
| Security | Each connection sees only own hidden hand; no duplicated account/seats, role escalation or unauthorized moves |
| Operator | Admin can investigate sanitized match audit, queue incidents and moderation via RBAC |
| Development | Hermes request produces real issue/PR, independent review, tests and verifiable DEMO deployment |
| Release | `main` ruleset protected and requires human approval, green CI, rollback plan and verified production deploy |

## Build milestones (in execution order)

- **M0 — Repository seed v0.1.0 (this delivery):** four authoritative specs + five defaulted decisions, nine project skills, AGENTS.md, static browser-only dev prototype, deterministic matchmaking-policy utility and Node tests, CI template, local git initial history. No human multiplayer yet.
- **M1 — Shared online engine and identity:** TypeScript authoritative Belot engine + deterministic fixtures, authenticated users + usernames, WebSocket game server, four human seats. Server-owned shuffle/actions/reconnect.
- **M2 — Entry flows/social:** signed-in home, solo/duo public matchmaker D01–D08, Create Room D09, friendship, profile overlay and username invites. End-to-end real-device multiplayer acceptance.
- **M3 — DevOps and agent:** protected GitHub `main`/`demo`, isolated deployment and databases, CI E2E, admin dashboard, Telegram Hermes + 9 project skills, feature PR to DEMO workflow.
- **M4 — Public beta and release:** automated load/latency/security tests, moderation, observability, audited production promotion and rollback. Actual MAIN release only after owner approval.

## Definition of done for an individual change

Behavior and edge cases documented; deterministic tests added; ruleset compatibility considered; security and hidden data checked; PR reviewed and CI green; DEMO URL verified (or explicitly absent); owner can test; change is **not MAIN** until separately approved and deployed. Any reported UI screenshot, build URL or deployment must be a real verified artifact, not a plan.
