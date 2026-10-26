# Agent instructions — Belot

> **MANDATORY FIRST READ: [CONSTITUTION.md](CONSTITUTION.md)** — the governing product vision and operating contract for Hermes, Codex and all development agents. Read it **before** skills, plans, tool use or code changes; follow its rules for Rado's nontechnical Telegram workflow, 80/20 delivery, rules verification, and explicit MAIN approval. If another instruction conflicts, surface the conflict rather than silently overriding this Constitution.


Read `README.md` and all four files in `docs/` before structural changes. Player experience: four **real** signed-in humans; home Start Random Game / Create a Room. **No bot fill** in real games. `prototype/index.html` is engineering-only. Confirm owner decisions D01–D04; D05–D09 are autonomous defaults in `docs/01-game-product.md`.

## Locked stack

- Backend: **Go**, server-authoritative. Code under `services/game`; test with `cd services/game && go test ./...`.
- UI: **React + Vite + TypeScript** under `apps/web`; build via npm in `apps/web`.
- `prototype/`, `main/` and `demo/` directory snapshots are legacy references; Git branches `main` and `demo` define real environments.
- GitHub repository is **public**. Never commit credentials or private customer/player data.

## Work policy

- Use `skills/feature-delivery/SKILL.md` on every code change; add `skills/gameplay-qa/SKILL.md` for game state/rules/matchmaking and `skills/admin-ops/SKILL.md` for moderation/incident/admin.
- Work in feature branch from `demo`; tests before PR; no direct production merge or production secrets.
- Treat score and hidden card integrity as security requirements. Server must validate every action, own all concealed cards and serialize turns. Do not trust user-supplied team, seat, hand or rules outcome.
- Keep docs and acceptance tests synchronized with implementation. When game variants are uncertain, add golden fixtures and state exactly what remains unverified; do not silently guess.
- `main` and `demo` folder HTML is an identical *initial static prototype snapshot*. Prefer editing `prototype/index.html` and running `npm run sync:prototype`, rather than maintaining two separate files. Actual environments are Git branches, not directories, once production app exists.
- Any update sent to the owner must distinguish **implemented**, **tested**, **deployed on DEMO**, and **approved for MAIN**. Do not imply deployment just because a ZIP or commit exists.

## Ask Niki — mentoring perspective for Rado

- Discover and use `.hermes/skills/ask-niki/SKILL.md` automatically when Rado asks "what would Niki do?", "ask Niki", wants to understand the AI-native product-building philosophy, or needs help thinking as a product owner rather than a programmer.
- The complete original explanation to Rado is preserved in `.hermes/skills/ask-niki/references/message-to-rado.md`. Do not overwhelm him with the whole message unless he requests it; give a concise principle and an actionable next step instead.
- This skill provides **Niki's documented perspective**, not direct contact with Niki or an invented personal approval. Combine it with `understand-rado`, `product-80-20`, and `teach-while-building` as appropriate. Never replace actual development work with motivational text.
