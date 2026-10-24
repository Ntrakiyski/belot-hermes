---
name: admin-ops
description: Operate Belot admin workflows for queue incidents, player reports, score disputes, version rollout and rollback with least privilege and audit logs.
---

# Admin operations skill

1. Verify admin identity, MFA and RBAC before reading or mutating player state. Never expose another player's unrevealed hand during an active match.
2. For report: find account/match, inspect sanitized event log and rulesetVersion, reproduce through deterministic engine, capture evidence and create an issue when bug exists.
3. For long queue: inspect real queue counts/ticket ages, fairness fallback D02/D03 and D05/D06 ghosts; never invent a bot to solve queue shortages.
4. For disconnect/forfeit: verify D08 grace, resumed event sequence, idempotent timeout and explicit reason, not ad-hoc score changes.
5. For release: check real health dashboard/deployment hash, errors/latency and backwards-compatible migration; rollback only via authorized production procedure.
6. Record all admin reads/mutations according to audit policy. Share anonymized incident summary with owner; never paste production secrets, tokens, private user fields or other hidden game state in Discord.
