# Pillar 03 — Hermes on Discord: automated development and governed promotion

**Owner's goal:** speak/type feature ideas in English in an authorized Discord channel; Hermes turns them into work items, code and reviewable DEMO changes, then waits for human permission before MAIN. The integration is a planned workflow, not connected in seed v0.1.0.

## 1. Workflow

`Discord text/voice -> authorized Hermes gateway -> speech-to-text if needed -> task summary / GitHub issue -> isolated feature branch -> code + tests -> agent review -> PR to demo -> CI -> deploy DEMO -> Discord evidence + DEMO link -> owner's explicit acceptance -> human-reviewed PR to main -> production deploy + monitoring`.

Recommended channels: `#belot-ideas` (feature input), `#belot-builds` (progress and demo artifacts), `#belot-releases` (approvals/receipts), `#belot-incidents` (restricted), optional voice channel. Restrict code requests to allowlisted Discord identity/role; a public player's content must never become commands.

## 2. Three agent skills

- `skills/feature-delivery/SKILL.md`: turn intent into scoped issue, implementation/test/PR, DEMO evidence, promotion gate.
- `skills/gameplay-qa/SKILL.md`: rules fixtures, card security, multi-client/reconnect/score validation.
- `skills/admin-ops/SKILL.md`: player reports, score disputes, queue incidents, sanitized audit and reversible remediation.

Hermes should load the relevant skill automatically based on requested task, and all code tasks should use feature-delivery plus relevant QA. Skills are **documentation**, not actual installed/run bots.

## 3. Permissions and guardrails

| Allowed without owner intervention | Requires explicit owner/human approval |
|---|---|
| Read code, docs and nonsecret logs; create issue and feature branch; modify code in isolated sandbox; add tests; open PR to `demo`; run CI; report DEMO link after verified deploy | Production promotion; deploy with privileged production secret; destructive DB migration; changes to auth/security policies; payment/billing changes; new external permissions; user data export |

- Hermes uses dedicated GitHub App or fine-grained account with repo-scoped permissions; no direct pushes to protected `main`.
- Approval must be backed by GitHub PR/review ruleset, not merely a prompt. Voice-only promotion is not sufficient: human confirms explicitly in trusted approval path.
- Runtime has no production DB credentials, no broad organization admin token and no access to private cards. Run code in ephemeral sandbox/worktree with time/resources limits and egress controls.
- Security policies remain enforced even if Discord input says to skip tests or merge directly. Refuse unsafe tasks and explain. Avoid committing secret values, transcripts containing private data, or database dumps.
- Every proposed merge includes Git commit SHA, changed files, test evidence, known limitations, DEMO URL (if actually deployed), rollback plan and acceptance result.

## 4. Standard feature receipt

```text
BEL-### — feature title
Understood: concise expected player behavior
Changed: files/PR + commit
Tests: commands and green/failing evidence
DEMO: real link or "not deployed"
Risks: known limitations
Next: Await owner acceptance / production approval
```

## 5. Decisions and autonomy

Owner delegated routine product decisions to the assistant for specification work (D05–D09 in Pillar 01). That is **not** authorization for Hermes to grant itself production release permission. Agent can choose reasonable non-destructive defaults, log them, and ask only if critical safety/security/release approval is required. If change has competitive integrity impact, the QA skill must verify rules against fixtures and record the rulesetVersion.

## 6. Setup prerequisites (not yet configured)

1. GitHub private repo with `main`+`demo` branches and protected `main`.
2. Discord bot/app, owner identity allowlist, approved channels and voice transcription setup.
3. Always-on Hermes host with code execution sandbox, fine-grained repository credentials and isolated DEMO build/deploy secrets.
4. GitHub Actions workflows, hosted DEMO API/frontend and separate MAIN hosting/databases.
5. Test owner voice/text request -> GitHub PR -> verified DEMO -> human approval -> MAIN and rollback.
