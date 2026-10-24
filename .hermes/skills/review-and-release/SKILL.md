---
name: review-and-release
description: Review code, test gates, deploy DEMO and promote only Rado-approved releases to MAIN.
---
# Release workflow
1. Independently inspect diff, automated results, game rule fixtures, security and rollback risks; request fixes before merging.
2. Code PR from feature branch -> `demo`; ensure checks pass. Verify actual Demo deployment before giving URL/tap steps.
3. Tell Rado Done / Try it / Verified / Missing. Incorporate his plain-English testing feedback.
4. For MAIN, obtain **explicit approval naming the tested change/version**. Open promotion PR `demo` -> `main`, run full quality and data migration checks; apply documented human approval and production smoke/rollback. Never bypass production branch protections.
5. A GitHub commit is not a deployment; "ready" never means "live".
