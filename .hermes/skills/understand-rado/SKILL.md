---
name: understand-rado
description: Convert Rado's vague Telegram voice or text requests into scoped, testable Belot changes without technical interrogation.
---
# Interpret Rado's requests
Trigger automatically for unclear, emotional, fragmented or incomplete statements. Treat speech-to-text as possibly mistaken.

1. Reflect the **player problem** in one short sentence. Look at current behavior and project docs; do not invent context.
2. Infer the likely action, user-visible result and expected benefit. Separate observations from assumptions. Offer one small sensible default.
3. Ask **at most one question** only if two plausible interpretations create materially different outcomes, a rule change, data loss, spending or security impact. Never ask micro-decisions Rado can reasonably delegate.
4. Write 2–5 `Given/When/Then` acceptance cases; include one failure/edge case. State what is deferred.
5. Load `product-80-20` for scope and `feature-planning` for implementation. Explain the proposed feature in plain English without architecture trivia.

Example: "Cards feel slow" => ask whether delays are in selecting or opponents responding only if measurements cannot distinguish; otherwise profile latency, improve the dominant wait, test on DEMO.
