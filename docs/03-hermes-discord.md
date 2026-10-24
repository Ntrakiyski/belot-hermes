# Pillar 03 — Hermes + Telegram autonomous development

**Canonical communications:** Rado talks to one dedicated **Belot Builder Hermes profile** through Telegram, by voice or text in plain English. Discord was considered previously but is not the v1 primary channel. **Integration is planned, not installed on the user's VPS by this repository.**

## 1. End-to-end workflow

`Telegram voice/text -> Hermes transcript and intent interpretation -> high-impact clarification only if needed -> 80/20 scope -> GitHub issue/plan -> feature branch from demo -> implement Go/React + tests -> independent review -> PR to demo -> CI -> deployed DEMO with restricted rules inspector -> Rado tests on phone -> explicit MAIN approval -> reviewed promotion PR -> production deploy and smoke test`.

Rado need not use GitHub or know Git terminology. Hermes should translate product language into acceptance checks, select skills automatically, educate Rado gradually and offer sound recommendations.

## 2. Skills and automatic routing

Project-local skills in `.hermes/skills/`: `understand-rado`, `product-80-20`, `feature-planning`, `implementation-quality`, `belot-rules-verification`, `demo-inspector`, `review-and-release`, `teach-while-building`, `skill-discovery`. The agent proactively loads relevant skills from natural-language intent. `AGENTS.md` is the routing and security contract; `ops/hermes/SOUL.md` defines Rado-aware behavior. Existing root `skills/` are retained as legacy reference and are not sufficient for Hermes project discovery.

**Find Skills**: the official CLI can discover additional skills; install/enable on the dedicated VPS profile only after source/code/permission review. Skills cannot override Belot game integrity or MAIN approval rules.

## 3. Boundaries and approvals

Hermes may decide reversible MVP defaults, create feature branches, code, tests, PRs into DEMO and run sandbox tests. It must obtain meaningful operator approval before executing newly discovered untrusted scripts or granting credentials. MAIN promotions, material data deletion, new production permissions, major spend and policy changes require explicit human authorization. Public GitHub repo: never add credentials or private player data.

- Do not grant the Telegram runtime production DB credentials or release authority.
- Rado's Telegram identity must be allowlisted; ordinary player content must not become agent instructions.
- Rule changes require independently checked deterministic fixtures and a documented ruleset version.
- Demo inspector data must be isolated, staff-authorized and impossible to expose in MAIN through a client flag.
- Feature status must distinguish **coded, tested, in GitHub, deployed to Demo and live in Main**.

## 4. Feature update format for Rado

- **I understood:** one sentence about the user problem.
- **I changed:** the minimal behavior.
- **Try it in Demo:** only a *verified* URL with 2–3 tap steps, else say "not deployed".
- **Verified:** tests, replay evidence and known gaps.
- **One useful concept:** optional, one short explanation.
- **Next:** follow-up or request explicit approval for a specified production release.

## 5. Setup prerequisites and proof

Dedicated Hermes VPS profile, private Telegram token/allowlist, trusted project checkout, GitHub feature-write permissions, restricted DEMO deployment secrets, GitHub CI/branch protections and isolated environment/database. Follow `docs/09-hermes-vps-setup.md` and `docs/05-development-playbook.md`. Required acceptance: Rado sends a vague voice note, Hermes creates/test/reviews a feature, reports real DEMO evidence, and MAIN stays unchanged until a clear human-approved promotion.
