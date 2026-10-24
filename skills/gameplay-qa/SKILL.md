---
name: gameplay-qa
description: Test Bulgarian Belot rules, matchmaking constraints, secret cards, server-side turns, replays, reconnect and scoring.
---

# Gameplay QA skill

1. Read Pillar 01. Identify which decision (D01–D09), bidding, legal play, announcements or scoring is affected.
2. Derive deterministic fixtures: 32 unique cards, five-plus-three dealing, valid contract sequence, each legal/illegal play, trick winner, announcements, contra/recontra, last trick, capot, carried scores, match end, variant restrictions.
3. For matchmaking fixtures test: solo+solo+solo+solo; duo+duo; duo+solo+solo at 4.999/5/15 seconds; fewer than four cannot start; duplicate account; wrong seat; disconnected user during countdown; atomic parties.
4. Add four-device integration tests to prove only each player's own hand is visible; server rejects forged card/seat/team actions, double clicks and out-of-order moves. Reconnect must replay ordered public events and return only player's private state.
5. Run unit, property-based / seeded 10k-deal simulations, websocket and mobile browser E2E as available. Separate executed tests from proposed ones.
6. Produce concise evidence with exact command results, fixture names, rulesetVersion and unresolved Bulgaria house-rule questions. Block DEMO merge if correctness/secret-card invariant fails.
