# Product roadmap — 80/20 thin slices

**Version 1 = engineering workspace**, not multiplayer launch.

## v1 (this handoff)
- Canonical business pillars, autonomous engineering playbook, Rado interaction policy, nine focused project skills, VPS profile setup and migration notes.
- Repo with `main` / `demo`, Go backend card-deck starter, React/Vite UI scaffold, tests and CI definition.
- No claim of live Telegram integration or deployed environments. The owner sets up the isolated VPS Hermes profile and bot with credentials outside Git.

## Smallest useful implementation sequence for Hermes
**Slice A — rules before multiplayer:** pure Go deal/bid/turn/score engine and golden tests. Success: reproducible full deal and point explanation without browser.
**Slice B — playable local four-seat fixture:** React table reading server-controlled state, deterministic test players only; authorized DEMO inspector with step/replay/legality reasons.
**Slice C — first four humans:** minimal accounts, one match lobby, WebSocket turn events, reconnect and anti-leak tests. The first public-real multiplayer proof is 4 separate authenticated devices completing a match.
**Slice D — main entry flows:** Start Random Game, solo/duo, balance 5–15s, countdown; then private rooms, username/friends invites.
**Slice E — operations/release:** separate production database, moderation, monitoring, secure sandbox controls, versioned promotion/rollback.

Do not build rankings, tournaments, cosmetics, voice chat, payments or ML opponent scoring before slices A–E meet acceptance. Product metrics: completed real-human matches, rule fixture pass rate, clean reconnects, queue waits, Rado accepted demo releases.
