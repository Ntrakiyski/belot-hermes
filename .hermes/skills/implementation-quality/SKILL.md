---
name: implementation-quality
description: Deliver Go and React changes with correctness, security, tests, review and usable Demo evidence.
---
# Engineering execution
1. Confirm branch from `demo` and scope from a plan; map changes across Go engine, React UI and persistence.
2. For game behavior, add failing deterministic Go tests/fixtures **before** implementation. For UX, define action/state and browser acceptance check.
3. Enforce server authority: authenticate seat and turn, validate events, reject duplicates, redact other hands. Handle reconnect races.
4. Implement smallest cohesive change; format/test Go (`go test ./...`), Node tests, TypeScript build, lint/security scans if configured.
5. Self-review actual diff for concurrency, hidden data, rules variants, backward compatibility; fix failures and update docs.
6. Never claim deployed before testing a real DEMO URL; produce a reproducible evidence summary for Rado.
