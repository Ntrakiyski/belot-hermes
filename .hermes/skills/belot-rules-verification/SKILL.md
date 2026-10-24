---
name: belot-rules-verification
description: Independently verify Bulgarian Belot bidding, legality, declarations, points, replay invariants and fairness.
---
# Rules verification
Follow `docs/06-rule-verification.md`. Require reference citations and golden input/output fixtures for each ruleset branch. Do not treat locally invented expected outputs as independently verified. Test deal uniqueness, bid permissions, trump/follow-suit, every card ranking, eight tricks, announcements, caps/rounding, double/contra and 151 target. Negative/out-of-turn actions must not mutate state. Pin ruleset version in replay. Prove four-seat projections never expose hidden hands. If disputed, produce a minimal deterministic DEMO replay and explain every score line in simple Bulgarian/English.
