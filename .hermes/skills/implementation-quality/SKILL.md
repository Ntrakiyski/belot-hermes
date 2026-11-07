---
name: implementation-quality
description: Deliver Go and React changes with correctness, security, tests, review and usable Demo evidence.
---
# Engineering execution
1. Confirm branch from `demo` and scope from a plan; map changes across Go engine, React UI and persistence.
2. For game behavior, add failing deterministic Go tests/fixtures **before** implementation. For UX, define action/state and browser acceptance check.
3. Enforce server authority: authenticate seat and turn, validate events, reject duplicates, redact other hands. Handle reconnect races.
4. Implement the smallest cohesive change. Use the working directories and installed sandbox toolchain recipe in `belot-maintenance`: Go tests/vet in `services/game`, `npm run check` and frontend build from repo root. Race tests run in `go-engine` CI when the sandbox lacks a C compiler; Compose is covered by `compose-smoke`. Do not run Go tests at repo root.
5. Self-review actual diff for concurrency, hidden data, rules variants, backward compatibility; fix failures and update docs.
6. Never claim deployed before testing a real DEMO URL; produce a reproducible evidence summary for Rado.
