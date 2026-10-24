---
name: feature-delivery
description: Safely turn Discord feature requests into GitHub issues, tested feature changes, DEMO evidence and human-approved MAIN promotion.
---

# Feature delivery skill

1. Authenticate requestor from configured Discord allowlist. Turn voice transcription or text into concise testable issue. Default target: `demo`, never `main`.
2. Read `AGENTS.md`, Pillar 01 and affected specs. Record behavior, constraints, acceptance tests and backwards compatibility; use documented defaults rather than repeatedly asking trivial questions.
3. Create issue and branch from latest `demo`. Work in sandbox with narrow repository credentials. Make minimal changes and tests; check existing user/session security.
4. Run `npm test`, formatting/type/build and later E2E; run independent diff/security review, add required golden rules fixtures if applicable.
5. Open PR into `demo`; after checks pass, merge only if bot permissions and repository rules allow. Verify real DEMO deployment and smoke test. If hosting absent, say **not deployed**.
6. Reply with summary, issue/PR, commit hash, actual tests, limitations, DEMO URL when real. Wait for owner to test.
7. MAIN promotion requires explicit owner acceptance AND protected GitHub PR reviewed by human, successful CI, production gate and rollback plan. Never self-approve or bypass protection.

**Failure modes:** stale branch -> rebase/test again; tests fail -> fix and re-run; request alters permissions/data -> require approval; incident -> `admin-ops`; score/card bug -> `gameplay-qa`.
