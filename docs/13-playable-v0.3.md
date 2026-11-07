# Playable v0.3 plan and reference comparison

## Evidence and scope

Compared our DEMO with [nikomarinovic/Belot-Card-Game](https://github.com/nikomarinovic/Belot-Card-Game), commit `b76753615b5b6725e87f90b42320ea4cab4b8502`. The reference is Croatian iOS Belot, not a Bulgarian rules oracle. Its explicit bot think time, fixed table positions, visible completed tricks, declaration prompts and score presentation are useful interaction references. Its different targets, forced trump and deal structure do not replace our six-contract Bulgarian game. Source review also found inconsistent ordinary-ten ranking in its trick winner and legal-card routines.

No code or assets are copied. The reference has an MIT LICENSE but contradictory restrictive README wording; this update independently implements behaviour using our Go/React stack and original card rendering.

Bulgarian rule fixtures use [Belot.BG's published rules](https://belot.bg/belot/rules/) and [Belot Arena’s published scoring implementation](https://belotarena.com/tochki-v-belota) for paired rounding detail where the summary leaves it implicit. Engine tests pin expected outcomes; no claim of complete certification is made.

## Confirmed problems

1. Synchronous bots play several cards/tricks before React sees the result. Completed tricks disappear; no useful winner/points feedback survives.
2. Declarations and belot are missing, tied hands are treated as failure, rounding is not contract-aware, and recontra ends bidding prematurely.
3. Mobile cards wrap unevenly; seat direction is reversed; lower bids and illegal cards remain clickable; several explanations are English.
4. Practice requires hidden lobby steps; disabled practice still appears; stale sessions silently spin forever.
5. Version and scoring explanations are not sufficient for Rado to validate a match.

## Chosen update

Retain the existing authoritative Go engine and guest-table API. Extend that vertical slice rather than port the iOS app or build accounts/matchmaking infrastructure. Improve rules, table pacing, practice entry and guest-room usability together. UI-only polish would leave wrong scores and invisible tricks; a complete backend rewrite adds unnecessary risk.

- Server advances at most one bot action per scheduled deadline: 800 ms for bidding, 650 ms for play. Hold completed tricks for 1800 ms and reject card/bid/next actions during collection. Use a single store ticker with deterministic clock-driven tests.
- Preserve all four cards, winner, captured points, current deal, bid history and match score sheet. Expose only the requesting seat's private hand/declaration options; published declaration summaries contain no hidden cards.
- Add explicit first-card declaration choices, optional correctly timed belot, source-backed score/carry fixes and ruleset `bg-0.3.0`. Keep all six contracts, five-plus-three bidding and a 151-point match.
- Bulgarian card-table UI, legal/disabled controls, clear friend-room invitations, direct labelled practice when enabled, session recovery and version display.

## API additions

`/healthz`: `version`, `practiceEnabled`. Existing fields remain.

Table view: `seq`, `dealNumber`, `collecting`, `trickWinner`, `trickPoints`, `nextActionAt`, `captured`, `pending`, `bidHistory`, `history`, `declarationOptions` (owner-only), `declared` (public summaries), `canBelot` (own legal cards), `belotPoints`, `practiceEnabled`, `version`.

Play body retains `suit` and `rank`; adds optional `declarations: string[]` and `belot: boolean`. Existing clients may omit declarations. Authentication uses `X-Token`; table codes/invite links carry no seat token.

## Acceptance

- One-click explicit practice starts a labelled table; friend tables never silently fill bots.
- Each bot bid/card is visible, four-card trick/winner persists for collection, and no fifth card or next-deal click bypasses the hold.
- Higher bidding, contra/recontra, declarations, legal plays, hanging carry and explained scores pass deterministic regressions and full-deal simulations.
- A phone player completes a visible deal through score and next deal; refreshed sessions recover, deleted tables show a useful way home and network failures remain retryable.
- All seat projections hide other hands/tokens; malformed/out-of-turn/illegal commands preserve game state.

## Release and limitations

User requested an app version update after authorizing publication and DEMO deployment. Release via reviewed feature PR and green CI to DEMO only; MAIN remains unchanged. Deployment resets in-memory tables. Guest codes do not prove four distinct authenticated humans; persistence, public matchmaking, reconnect/forfeit policy and privileged replay inspector remain separate milestones.


## Source-backed fixture summary

| Concern | Regression coverage |
|---|---|
| Deal and auction | 3+2+3 packets from dealer’s right; recontra stays open; all-pass redeal |
| Cards and turns | 10,000 seeded deals across six contracts, no duplicates, all 32 cards played, eight tricks, exact card-point totals |
| Scores | Hanging carry/accumulation, pair-rounding boundaries, whole-pot doubles, NT doubling before rounding, capot and no-capot match exit |
| Declarations | First-card/NT restrictions, sequences/fours and independent competition, equal sequences cancel, overlapping choices reject atomically, timed belot |
| Runtime | 40 complete paced practice matches using virtual deadlines; early/double inputs reject; completed tricks hold and advance once |

Local tests compare implementations with these published examples. Regional variants beyond the pinned examples need Rado’s review; the reference iOS app itself is not ground truth. No production approval is implied.
