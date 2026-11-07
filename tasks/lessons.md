# Lessons

## 2026-11-07 - Inspect workspace before cloning

Mistake: Attempted cloning into a nonempty workspace root.
Why it happened: File search did not show the hidden .git directory.
Rule for next time: Inspect hidden entries and Git state before selecting a clone destination.
Example check: Run ls -la and git status before git clone.

## 2026-11-07 - Observe route status before asserting it

Mistake: Expected disabled bot POST to return 405; the actual Go file server returns 404.
Why it happened: Assumed method handling before testing the packaged runtime.
Rule for next time: Verify endpoint behavior against the container before documenting exact status codes.
Example check: Run the HTTP smoke against Compose, not an inferred router behavior.

## 2026-11-07 - Validate the action response, not only the following GET

Mistake: The first practice smoke checked a GET after adding bots and missed the malformed POST response that crashed React.
Why it happened: The smoke did not follow the UI's response consumption path.
Rule for next time: Assert action-response contracts as well as persisted state, and run the browser flow.
Example check: Filling bots must itself return seats and the requester's hand before React renders.

## 2026-11-07 - Respect module working directories

Mistake: Ran a Go package test from the repository root, outside the Go module.
Why it happened: Used the edit command's cwd for a separate test invocation.
Rule for next time: Run Go commands from services/game as AGENTS.md specifies.
Example check: Verify the command cwd contains go.mod.

## 2026-11-07 - Distinguish Coolify API from connector support

Mistake: Tried the documented is_auto_deploy_enabled setting through a connector that filters out that field.
Why it happened: The tool accepts extra keys syntactically but its implementation does not forward them.
Rule for next time: Inspect returned warnings and verify the saved state; use the native supported action when the connector lacks it.
Example check: Confirm settings.is_auto_deploy_enabled after migration setup instead of assuming an accepted API field was applied.

## 2026-11-07 - Check each configured public hostname independently

Mistake: Initial combined curl command stopped at the www DNS failure before running the canonical-domain game smoke.
Why it happened: Independent hostname checks were chained as dependencies.
Rule for next time: Separate optional aliases from canonical application verification and report DNS, proxy and app health distinctly.
Example check: Verify canonical HTTPS + smoke even when the optional www alias has no record.

## 2026-11-07 - Prove the game experience, not only HTTP success

Mistake: Previous deployment was technically playable through APIs, but its immediate bot loop hid turns/tricks and the simplified rules did not deliver a convincing Belot game.
Why it happened: Infrastructure smoke tests and a single card click were treated as enough experience evidence.
Rule for next time: Verify the complete player loop against a reference: entry, bidding, paced turns, readable trick winner, declarations, score and next deal.
Example check: Play a full visible deal on a phone-sized viewport, checking who played and why points were awarded.

## 2026-11-07 - Clear one-shot gameplay selections after accepted actions

Mistake: First v0.3 UI draft kept first-card declaration IDs selected for later card plays, which the authoritative server correctly rejects as late declarations.
Why it happened: Selection reset only on table/deal change, not after the one-shot action succeeded.
Rule for next time: Clear one-shot input state after accepted actions and filter submitted options against the current server capability list.
Example check: Announce with the first card, then play the second card in the same deal without a declaration error.

## 2026-11-07 - Distinguish missing CLI from missing GitHub access

Mistake: Earlier audit highlighted absent gh/profile tokens without first establishing whether Hermes uses a different configured GitHub integration.
Why it happened: Assumed one authentication/tool path was the only available path.
Rule for next time: Respect confirmed access, inspect the actual integration, and test a harmless permitted operation instead of inferring a blocker from one missing binary.
Example check: Confirm the dedicated agent can inspect this repository/PRs using its existing tools; do not create or request replacement credentials.

## 2026-11-07 - Check infrastructure capacity before adding optional components

Mistake: Started a disposable transcription environment before reviewing the disk-usage command's result; the VPS had no normal free space.
Why it happened: Capacity check and installation were placed in one command instead of making installation conditional on the observed result.
Rule for next time: Inspect capacity and confirm the current feature scope before any dependency/model installation; keep optional additions isolated and removable.
Example check: Read df/available-memory output first, then decide whether a requested component should be installed at all.
