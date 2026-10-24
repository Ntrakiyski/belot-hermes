# Pillar 02 — Architecture, environments, administration and release

**Status:** target design, not an implemented online system. **Invariant:** ONE GitHub repository with `main` (owner-approved) and `demo` (development/QA) branches. Demo and Main start at exactly the same initial commit, evolve separately only through changes reviewed/promoted via PR. The `main/` and `demo/` directories in seed v0.1.0 are **identical local snapshots**, not deployments or permanent duplicate applications.

## 1. Implementation architecture (target)

```text
Mobile/desktop browser (Bulgarian React/Next.js or Vite SPA)
  -> HTTP authentication + profile/friends/invitation API
  -> persistent WebSocket multiplayer connection
  -> authoritative Node.js game server + pure TypeScript Belot rules package
  -> PostgreSQL (accounts, games, scores, friend events, replay history)
  -> Redis (distributed matchmaking queue, presence, reconnect tokens, optional)
  -> structured audit/telemetry + admin dashboard (/admin, RBAC + MFA)
```

First implement single regional server process so atomic queue and game events are correct; Redis becomes appropriate with multiple instances. Stateless frontends can use Vercel or equivalent, but websocket server must run in a persistent container/VPS. Choose deployment provider after tests; **do not infer that this repository is deployed**.

## 2. Trust and data isolation

- Authentication: signed-in accounts only; globally unique username, unguessable account ID, secure session storage, HTTPS. Rate-limit invite, username lookup, match requests and login. Validate reported/banned state server-side.
- Matchmaking: tickets are atomic solo/duo; one active queue or match/account; server-atomic reservations; server selects teams and timestamps D03/D04 boundaries. No client authority for playing as another seat.
- Game: rules engine is a pure function over validated events; ordered event log with sequence/idempotency keys. Shuffle, secret hands and legal actions on server only. Never send other hands to a player (including JSON debug fields or live admin endpoints).
- Reconnection: authenticated session + match-scoped tokens. Preserve seat and private hand, reject stale prior socket actions; heartbeat and D05/D06/D08 timers maintained by server monotonic clock.
- Persistence: account, friendships, party/invite, queue metadata, room/seat, match/deal/trick/action, rulesetVersion, score reason, audit. Retention/anonymization policies before launch. Never log raw tokens or unseen hands in client-facing logs.
- Admin: MFA, per-action RBAC and audit. Operators can inspect public match event timelines without spying on unrevealed hands. Admin can block users, review reports, inspect match status, trigger safe remediation; cannot edit finished scores without recorded audited process.

## 3. Environments and release controls

| Property | DEMO | MAIN |
|---|---|---|
| Git branch | `demo` | `main` |
| Audience | Owner, testers, Hermes | Public users |
| Change entry | feature PR -> `demo` after green checks | Promotion PR from tested `demo` -> `main`, approved by human |
| Database, Redis, secrets | Isolated sandbox | Separate production resources |
| Telegram Hermes write | Dev branch and demo deployment only | No direct merge/deploy bypass |
| Domain | `demo.<chosen-domain>` after hosting | `<chosen-domain>` after hosting |
| Target code | Same app + environment config | Same app + environment config |

Create a **private GitHub repo** first. Default branch `main`. Keep all changes on feature branches, review before merging to `demo`; enable repository ruleset/protection on `main` requiring green CI, human approval and no force pushes. Create deployment environments `demo` and `production` in hosting; production secrets and write credentials never accessible to Hermes runtime. A successful Git commit is NOT automatically a successful deployment.

### CI/CD pipeline

1. On PR: formatter/lint/type check, unit rules + queue tests, deterministic full-match fixtures, security & secret checks, build, browser smoke, license/dependency review.
2. On approved merge to `demo`: deploy DEMO; verify health/readiness, websocket handshake, isolated sample match; report URL + commit and automated evidence to Telegram.
3. Owner tests and explicitly authorizes promotion; open PR against `main`, verify required checks (including migrations/backward compatibility), get human approval.
4. On `main` release: deploy production with backward-compatible migrations, smoke test, monitor, rollback on regression; preserve event/ruleset compatibility for in-progress matches.

## 4. Data/state model (target)

- `Account`, `FriendRequest`, `Friendship`, `BlockedAccount`, `UserReport`.
- `Party`, `PartyInvitation`, `QueueTicket`, `MatchReservation` (four account IDs and countdown state).
- `PrivateRoom`, `RoomInvitation`, `Seat` (ownership and team positions).
- `Match`, `Deal`, `Bid`, `Trick`, `Announcement`, `ScoreLedger`, `GameEvent` (ordered and pinned ruleset version).
- `ReconnectSession`, `AdminAuditEntry`, `ReleaseArtifact`, `FeatureFlag`.

## 5. MVP admin features

Dashboard: online users, queue depth and latency, active matches, disconnect/forfeit counts, version/health. Search: account / match ID; view sanitized action replay, exact score explanation and rulesetVersion. Moderation: reports, blocks, suspensions. Feature flags must be auditable and cannot change rules midway through a deal.

## 6. Repo seed v0.1.0 vs target

Present now: static offline practice prototype (three bots), four docs, nine project skills, isolated match-selection policy example with tests, CI template. Missing: signed-in user accounts, four-human matchmaking server, WebSocket play, private rooms, friends, Admin, actual cloud/demo deployments and Hermes integration. **Do not enable public launch until all P0 items pass.**
