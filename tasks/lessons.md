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
