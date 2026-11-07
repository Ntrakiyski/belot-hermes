---
name: belot-rules-verification
description: Independently verify Bulgarian Belot bidding, legality, declarations, points, replay invariants and fairness.
---
# Rules verification
Follow `docs/06-rule-verification.md`. Require reference citations and golden input/output fixtures for each ruleset branch. Do not treat locally invented expected outputs as independently verified. Test deal uniqueness, bid permissions, trump/follow-suit, every card ranking, eight tricks, announcements, caps/rounding, double/contra and 151 target. Negative/out-of-turn actions must not mutate state. Pin ruleset version in replay. Prove four-seat projections never expose hidden hands. If disputed, reproduce with the current deterministic Go fixtures and explain the visible trick/score ledger. A privileged DEMO replay inspector is not yet implemented; never invent replay output. Consult `belot-maintenance` and `docs/13-playable-v0.3.md` for the implemented evidence path.
