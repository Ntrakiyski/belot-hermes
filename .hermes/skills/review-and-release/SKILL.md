---
name: review-and-release
description: Review code, test gates, deploy DEMO and promote only Rado-approved releases to MAIN.
---
# Release workflow
1. Independently inspect diff, automated results, game rule fixtures, security and rollback risks; request fixes before merging.
2. Code PR from feature branch -> `demo`; wait for `go-engine`, `react-vite`, `legacy-policy` and `compose-smoke`, then integrate only within authorized DEMO scope. The merge automatically triggers Coolify. Use existing GitHub access; no new hosting token or manual Pages deployment. Verify actual HTTPS DEMO version, affected gameplay and available deployed-commit evidence before saying live. See `belot-maintenance` for the concrete workflow.
3. Tell Rado Done / Try it / Verified / Missing. Incorporate his plain-English testing feedback.
4. For MAIN, obtain **explicit approval naming the tested change/version**. Open promotion PR `demo` -> `main`, run full quality and data migration checks; apply documented human approval and production smoke/rollback. Never bypass production branch protections.
5. A GitHub commit is not a deployment; "ready" never means "live".
