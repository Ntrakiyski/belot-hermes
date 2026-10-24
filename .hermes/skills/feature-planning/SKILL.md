---
name: feature-planning
description: Turn approved product intent into thin-slice code plan, explicit acceptance tests and safe change boundaries.
---
# Feature plan
1. Read all relevant docs and existing code; distinguish specified / implemented / deployed.
2. Record a short issue with outcome, minimal scope, assumptions, `Given/When/Then` acceptance, edge cases, dependencies, rollback and data/security impact.
3. Choose smallest vertical slice; do not design a complete platform for a simple request.
4. On rules or architecture changes, update `docs/11-tech-decisions.md` and tests/fixtures.
5. Create `feat/...` or `fix/...` from `demo`; pass task to `implementation-quality` and eventual `review-and-release`.
