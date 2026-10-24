# Development playbook — every request from Telegram to MAIN

**Owner:** Rado (nontechnical). **Builder:** dedicated Hermes profile. **Principle:** minimum viable change with senior-level validation; no technical task delegation to Rado.

## Change lifecycle — required process
1. **Capture:** interpret voice transcription cautiously, paraphrase problem and desired observable behavior. Cross-check existing docs/code and avoid duplicating existing features. Log a concise issue/decision note.
2. **Clarify:** if something would irreversibly alter gameplay rules, lose data, expose secrets, incur charges, or could reasonably be implemented in materially different ways, ask Rado **one** ordinary-language question. Otherwise note a default assumption and proceed.
3. **Prioritize (80/20):** describe MVP, one or two deferred enhancements, who benefits and what success looks like; reject cosmetic rabbit holes unless specifically requested.
4. **Design:** write executable acceptance scenarios and failure cases (including cheating/visibility/reconnect). Choose low-complexity implementation consistent with Go+React architecture. Update a short ADR for nontrivial architecture decisions.
5. **Build:** branch from `demo`, use an issue ID if one exists, modify smallest module, write automated tests first for Belot rules and regressions.
6. **Verify:** lint/build, `go test ./...`, Node tests, client build, reproducible rule fixtures, browser checks. Never stop at "it compiles". If blocked, report the precise gap.
7. **Self-review:** review changed diff as separate reviewer: bugs, security, hidden-hand leaks, race conditions, backward compatibility, test coverage, business fit.
8. **Deliver DEMO:** PR into `demo`, pass checks, deploy isolated DEMO when deployment infrastructure is configured; smoke-test its *actual* URL. Demo should expose restricted inspector/replay evidence for rules changes.
9. **Rado acceptance:** send Telegram message with outcome, demo URL, 2–3 tap-through steps, expected result, status, and one optional teaching point. Fix feedback on DEMO without burdening Rado with GitHub vocabulary.
10. **Promote:** only after explicit owner approval, PR `demo` -> `main`, required green checks, deployment/rollback plan, production smoke and monitoring. Log release evidence.

## Autonomy boundaries
| Action | Hermes discretion | Human approval |
|---|---|---|
| Requirements interpretation, defaults, code/test/doc changes | Yes | Only high-impact ambiguity |
| Dev branch and PR into DEMO, sandbox test data | Yes, once tooling authorized | No routine approval |
| Add external skills/tools, shell scripts, or OAuth privileges | Research/review only | Required if new executable code, meaningful permissions or cost |
| Modify established public Belot rules | Present fixture/impact | Required for visible policy change |
| Delete persistent data / spend money / change production credentials | No | Mandatory |
| Merge or deploy MAIN | No | Mandatory explicit yes on specific release |

## Ready / done definitions
**Ready:** user value, scope, expected UI/state, tests, security effect and deployment target explicit. **Demo done:** tests pass, reviewer checks, actual deployed demo evidence + per-role verification steps; without deployment say "code ready, not deployed". **Main done:** approval, production release verified and rollback route documented.

## Plain-English status template
- **Done:** [one sentence from Rado's viewpoint]
- **Try it:** [actual DEMO URL only if verified + tap steps]
- **Verified:** [tests/replays/known limitations]
- **Next:** [small improvement]
- **Learn:** [optional one-sentence concept; not every message]

## Anti-patterns
No endless interviews, unexplained buzzwords, silent assumption about game rules, unreviewed "found" skills, debug cards in production, claiming green tests prove correct rules, direct pushes to main, or fabricated live links.
