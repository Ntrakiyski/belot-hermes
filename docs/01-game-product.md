# Pillar 01 — Belot game, rules, matchmaking and player UX

**Revision:** 2026-10-24 · **Authority:** owner's confirmed Decisions 01–04, assistant-chosen defaults for Decisions 05–09 per owner's delegated authority. These defaults are implementation decisions, not user-approved preferences. No further micro-approvals needed for the MVP. Changes remain editable later.

## 1. Product direction and entry points

**Four real players first**, mobile browser first, Bulgarian UI first. User signs up/signs in (unique username). The signed-in home screen has two dominant actions, in order:

1. **Start Random Game** — default, no room creation, matchmaking finds other people automatically. Supports solo queue or inviting **one accepted friend** as a two-person party. No 3–4 person public parties.
2. **Create a Room** — private four-human table, invitations by **username**, existing friends, and link/code.

Public and private games **must never silently substitute bots**. Engineering prototype has 1 human + 3 bots and is explicitly **not** this product. A friend's profile opens from username/avatar during play or after a match. Offer **Add Friend**, friendship status, report/block, and eventually stats. Neither clients nor the profile UI can reveal hidden cards.

## 2. Matchmaking — Decisions 01–09

| ID | Status | Decision |
|---|---|---|
| D01 | **User approved** | **Solo OR duo** queue, two accepted friends are permanently opposite-seat teammates in that match. |
| D02 | **User approved** | Prefer balanced composition: four solo players split into two teams, or duo vs duo; permit **one duo vs two solos** as fallback. |
| D03 | **User approved** | Composition relaxation window **5–15 seconds**, never guarantee a match without four real connected people. |
| D04 | **User approved** | On four connected participants, show teams and a **server-timed automatic 3-second countdown**; no Ready button. |
| D05 | **Default selected** | Disconnect during countdown: **pause, allow 5-second reconnect grace**, then replace missing account/atomic duo party, preserving eligible players' queue priority. Restart complete 3-second countdown after all four reconnect. Do not automatically split a duo or substitute a bot. If not replaceable within 20 seconds, release remaining tickets into queue preserving age; inform them. |
| D06 | **Default selected** | Queue stays active until player cancels or loses authenticated session. Show queue duration; at 60 seconds show **Still searching / Cancel**, but continue automatically. Cancel is instantaneous/idempotent and withdraws an entire duo ticket if either member cancels. If client connection is lost before matching, give 10 seconds to restore; after that remove ticket to prevent ghost matches. |
| D07 | **Default selected** | First release has **no ELO/rank gating**. Choose earliest eligible waiting tickets, preferentially matching compatible compositions, with regional latency preference if measured; never exclude enough players to cause an indefinite wait. Teams/seats are server-assigned; duo opposite; solo partners selected randomly with a seeded, auditable server process. |
| D08 | **Default selected** | After a match begins, reserve the disconnected account's exact seat/hand for **60-second reconnect grace**, temporarily pause gameplay and show countdown. Reconnect restores private hand and current public event sequence. After timeout, **forfeit the disconnected player's team** (not a scored normal Belot completion); record reason separately and protect teammates against abuse. No hidden bot takeover. |
| D09 | **Default selected** | Private rooms: 4 distinct signed-in humans, host invites by username/friend/link, no bots, host chooses seats/teams prior to start, all participants may view teams, host presses **Start** only after four clients connect; then shared 3-second countdown. If host leaves before start, host privilege transfers to longest-present player; empty room expires after 30 minutes. Host cannot override a player's card actions. |

### Composition algorithm (D02/D03)

- 0–4.999 sec from oldest eligible ticket: prefer **duo vs duo** or **four solos**. Do not form mixed duo + two solos solely for speed.
- From 5 sec: one duo + two solos becomes eligible when a balanced compatible match is not already available for those tickets. If a balanced match is available, prefer it.
- At 15 sec: composition preference must not block any available valid four-human match. **Underpopulation still waits**; 15 seconds is not a deadline to spawn bots.
- Tickets are an **atomic** solo (1 account) or duo (2 distinct accounts). Only one active party/queue/match per account. Queue priority is measured server-side. Never split duo friends into opposite teams.
- Server locks a match reservation before sending any start events. Four authenticated connected players must accept reservation implicitly through active connections; if connection fails, apply D05. Clients cannot forge countdown or create duplicate matches.

### Public user journey

`Sign up/in -> Home -> Start Random Game -> [solo OR invite one friend, wait acceptance] -> Searching + Cancel -> Server matches 4 humans -> Team confirmation and 3s countdown -> Bid/play eight tricks per deal -> Score -> Next deal until match conclusion -> Results -> Add friend/rematch`.

### Private-room user journey

`Sign in -> Create Room -> Invite username/friend/link -> Four join -> Host assigns two opposing teams -> Host Start -> 3s countdown -> Belot match -> Results/rematch`.

### Edge cases / expected outcomes

- Two players accept competing duo invites: server accepts exactly one active party; second invite invalidates cleanly.
- Same account opens two tabs/devices: one active game identity, reconnection replaces previous connection session safely; never five seats.
- Duo member leaves pre-match: atomic ticket removed, other friend offered solo queue; no silent conversion.
- Player drops during public countdown: D05, entire three-second countdown restarts only with four connected distinct humans.
- Player drops mid-trick: D08 keeps current hand hidden, preserves event order, no bot, later forfeit if not recovered.
- Profile opened mid-turn: legal move clock/timeout behavior must remain server-side; hidden hands stay private.
- Room invite full/expired/blocked/banned: reject server-side with clear state, not a half-joined seat.
- Users decline/reject friend invitation: persistent state updated, inviter notified, no implicit gameplay privileges.

## 3. Bulgarian Belot — provisional rules baseline (engine must be independently certified)

These are targets **not proof of full compliance**; finalize variations via fixtures before real-money, ranking or public launch. Rule source for comparison: https://belot.bg/belot/rules/

- 4 players, fixed opposite partnerships. Deck: 32 cards, suits ♣ ♦ ♥ ♠, ranks 7 8 9 10 J Q K A. 8 tricks per deal, one card from each player per trick.
- Deal 5 cards/player for bidding, then remaining 3 each after contract. Dealer advances each deal; all-pass redeals per chosen rule.
- Bid order `♣ < ♦ < ♥ < ♠ < no trump < all trump`; passes, contra, recontra validated by who is authorized and contract state. Server alone makes decisions about valid bids and start of play.
- Must follow led suit; trumping/overtrumping and rise requirements depend on contract, relative winning team and local rule variant. A player can **only** submit a server-verified legal card.
- Trump rank order: `J 9 A 10 K Q 8 7`; ordinary rank order: `A 10 K Q J 9 8 7`. No-trump uses ordinary order; all-trump uses trump order within all suits.
- Standard card values: trump `J20 9=14 A11 10=10 K4 Q3`; nontrump `A11 10=10 K4 Q3 J2`; last trick +10; capot/valat +90 as applicable.
- Announcements must include belot (K/Q suited requirement), tierce 20, quarte 50, quinte 100, four-of-kind, all-trump specifics, competing declarations and no-trump restrictions. **Explicit declaration timing** and counting/rounding require signed golden fixtures.
- Contract success, tied/hanging hands, contra x2/recontra x4, no-trump doubling, belot, capot, rounding and game-to-151 win conditions need golden fixtures, particularly local house variants. Display scores with detailed explanation, never silently improvise house rules.
- Target match: 151 or more *as defined by verified scoring policy*; equal scores require continuing until a clear winner. The 'cannot win on capot' variation should be configurable only after validating player consensus.

**Non-negotiable:** identical engine used to validate public and private tables; server owns shuffle (cryptographically sound RNG), secret cards, bidding, allowed moves, scores, and event sequence. Clients present views and ask for actions, never author game outcomes. Replayable events + pinned ruleset version are required for disputes.

## 4. UX and social specification

- Bulgarian locale initially; UI text localized rather than mixed-language technical terms, even though owner dictates requests in English.
- During game: local hand bottom, partner opposite, opponents left/right. Show deal/trick number, bidding controls, contract, team score, whose turn, reconnect status. Cards should be comfortably tappable on mobile.
- Profile overlay via player avatar/name: username, avatar, online state if privacy permits, **Add Friend / Pending / Friends** and **Report / Block**. Never show email or unrevealed game information.
- Search by exact username or prefix with throttling/abuse controls; unique canonical usernames; case-insensitive uniqueness. Blocking prevents invites, friend requests and direct social interaction.
- Private invitations require accepting; link joining may require host approval based on room privacy setting. Capacity maximum four. Rematch uses prior room participants only on mutual opt-in.
- Do not build unnecessary ranked, tournament, spectator or cosmetics systems for MVP.

## 5. Acceptance criteria for initial human-multiplayer release

- A signed-in player sees **Start Random Game**, then **Create a Room**, in that order; buttons are functional, not mockups.
- 4 distinct real accounts complete at least one public match on separate browser devices; one duo stays together opposite, four solo matching works, and mixed composition starts after configured relaxation.
- Automated boundary test cases cover 4.999 sec, 5.000 sec, 15.000 sec; fewer than 4 people always continue searching and do not spawn bots.
- Four live connections see identical three-second server-based countdown and a deal cannot start early, twice or with missing players. D05 dropout test exercises pause, replacement and countdown restart.
- Disconnect/reconnect during bidding and tricks preserves exact cards/event ordering; post-grace forfeit is labeled a **forfeit**, not silently scored as a normal contract.
- Private room supports invite-by-username, friend and link/code, correct seat permissions and four-human Start; host transfer tested.
- At least 10,000 seeded deal simulations have no duplicated/missing cards and no illegal moves; golden scoring fixtures cover all contracts, declarations, tie/hanging cases, contra/recontra, capot and rounding.
- Every action authorizes account/seat on server; no hidden card leaks via profile, API events, logs to other players or admin UI.
- Four human players can play a match to valid conclusion on phone browser; game completion, delay and disconnect telemetry reported.

## 6. Scope boundaries / deferred questions

All nine matchmaking decisions above are sufficiently decided for **planning**; exact inactivity timers within legal turns, anti-cheat mitigations, scoring fixtures and certain Bulgarian variants require implementation-level testing (not additional owner micro-approval). MVP prioritizes correctness, connectivity and fair 4-human play, not appearance or bots. The old local bot prototype remains a separate engineering aid only.
