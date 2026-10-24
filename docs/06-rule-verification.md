# Rule verification and DEMO inspector requirements

## Why a second way to view the match exists
A Belot score is only trustworthy when the **Go rules engine** can explain every deal/trick/point using deterministic evidence, including what each real player could legally see. DEMO's inspector is a development/test surface, **not** a public cheating aid.

## P0 golden cases and invariant tests
- Deck: 32 distinct valid cards, seeded deterministic deal in tests, production crypto-safe RNG, no hand exposed across seats.
- Bidding: 5+3 deal staging, valid bid order and pass/contra/recontra permissions; rejected invalid/out-of-turn command never mutates state.
- Every card turn: must follow suit; trump/overtrump/raise with all/no-trump variants; return a **reason code** for every legal/illegal move.
- Trick: compute winning seat and winner's team; exactly four played cards/trick; eight tricks/deal.
- Announcements: belot, sequences and four-of-kind with competing declarations, contract restrictions and point attribution.
- Points: cards + final trick + declarations + capot, no-trump/all-trump and doubles, hanging/failed contract, rounding and 151 target; use independently reviewed Bulgarian reference fixtures.
- State-machine properties: no skipped turns, no duplicate actions, serial event numbers, idempotent retry, reconnect never forks match.
- Matchmaking: four distinct authenticated humans, no bots, keep duo opposite, 5–15s balanced fallback, server countdown.

## DEMO-only inspector MVP (not yet implemented)
A separate route behind DEMO role authorization, e.g. `/dev/matches/:id`:
1. **Timeline:** events by sequence with timestamps and ruleset version; rewind/play/step controls.
2. **Trick view:** four seats, active turn, contract, last trick, legal cards highlighted per selected seat with concise reason text.
3. **Score ledger:** line items explaining every award/rounding and final winners; independently compare to reference fixtures.
4. **Deterministic test controls:** load a canned deal, choose a seeded scenario, auto-run a match and produce a shareable sanitized bug report.
5. **Visibility toggle:** inspect each seat's perspective without sending other hands to game clients; privileged server-side tester data only in isolated sandbox.
6. **Replay assertions:** fixture input -> expected output -> actual -> pass/fail; button to create a reproducible issue.

## Boundaries
- **Never enable full-hand introspection on MAIN**, regardless of query string, build flag supplied by client, or role claimed by browser; gate server-side by deployment environment AND staff authorization.
- DEMO must not contain live player identities/data or production secrets. Spectators and players only receive authorized projections.
- No assertion of "100% correct" absent validated fixtures; show coverage and gaps honestly.
- CI should run golden fixtures and Go fuzz/property tests. A code diff touching rules requires new fixture or explicit explanation.

## Evidence standard
For each bug report save ruleset version, event timeline with hidden hands redacted for public sharing, deterministic seed **only for non-live test games**, expected/actual result and regression test ID. Preserve seed only for authorized sandbox fixtures.
