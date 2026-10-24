# Rado's guide — no coding knowledge required

Talk to your dedicated Hermes bot in Telegram, in voice or text. Describe **the result you want**; code terminology is optional.

Examples:
- "When I tap Start Random Game, I want to see how many players we're waiting for."
- "I think the team got too many points from this trick. Can you show why?"
- "Make it easier to invite the person I just played with."
- "Try a simple version on Demo, and tell me where to tap."
- "I tested Demo and approve moving this exact version to Main." (For production releases only.)

Hermes should reply in ordinary English: a tiny plan, test evidence, an actual Demo link when it exists, and what to try. It will make low-risk decisions without interrogating you, and will teach one useful concept occasionally, for example: **server** = the computer that keeps the official game rules and state.

Main is for real players. Demo is your playground. Seeing a feature in GitHub does **not** mean it is deployed. Never send secrets, bot tokens or private passwords in Telegram or GitHub issues.
