# VPS handoff — dedicated Rado Hermes profile

**Not executed on the VPS by this repository author.** Run after cloning the public repository on your VPS, using the existing Hermes installation. Follow current official Hermes docs before changing live config.

1. `hermes profile create belot-rado` — separate config, sessions, memory, tools and skills from other Hermes profiles. Then `belot-rado setup` (credentials via secure local terminal, not Git or Telegram).
2. Clone `https://github.com/Ntrakiyski/belot-hermes.git` and launch the dedicated profile with the repo as its working directory.
3. Copy `ops/hermes/SOUL.md` into **that profile's** `SOUL.md` as appropriate for your Hermes version. Preserve a backup. The profile's default model/tools must support editing files, GitHub PRs and executing tests.
4. Trust only this reviewed checkout: from the repository root, run the profile equivalent of `hermes skills trust` (for some profiles the alias is `belot-rado skills trust`). Hermes discovers repo skills in `.hermes/skills/`, **not** the old root `skills/` folder. Confirm skills are listed.
5. Install discovery support on the **VPS/profile, not blindly in this session**. Official Find Skills source: `https://www.skills.sh/vercel-labs/skills/find-skills` — CLI: `npx skills add https://github.com/vercel-labs/skills --skill find-skills`. Check target agent support and location. Hermes may instead use its built-in skills system plus `npx skills find <query>` through `skill-discovery`. Inspect any newly suggested SKILL.md/scripts before using them; do not grant production tokens to untrusted packages.
6. Create a **separate Telegram bot** using @BotFather; configure its token and **Rado's numeric Telegram ID allowlist** in the profile's private `.env` through the Hermes Telegram setup wizard (`belot-rado gateway setup`, if available). Never reuse another Hermes profile's token; profiles need distinct bot tokens. Start gateway and verify voice transcription and private text; check users outside allowlist are rejected.
7. Configure GitHub auth for feature branch + DEMO PR, separate DEMO hosting credentials and scoped secrets; don't give the agent MAIN deployment credentials. Set up protection rules and environments in GitHub UI.
8. Test: Telegram vague feature request -> skill chosen -> plan -> code+tests -> PR -> DEMO deployed URL -> Rado checks; require explicit promotion approval. Verify secrets don't appear in public repo, chat logs or generated reports.

## Important
- This repo contains no Telegram token, VPS address, secret, deployment or verified integration.
- `skills.sh` is an ecosystem discovery mechanism, **not** permission to execute everything in its catalog. The desired `find-skills` skill complements custom policies; it does not replace them.
- Hermes can auto-select skills via project-local discovery after checkout trust and appropriate working directory. Triggering by natural language is the default; Rado should never type `/skill`.

Official docs: https://hermes-agent.nousresearch.com/docs/user-guide/features/skills and https://hermes-agent.nousresearch.com/docs/user-guide/messaging/telegram

## Verify Niki's mentor skill in the Rado profile
After the profile checks out the development workspace and trusts its project-local skills, confirm that `ask-niki` appears in the loaded skill list. From Rado's authorized Telegram session, try: **"What would Niki tell me if I don't know how to code but want to improve Belot?"** The agent should respond with a short product-owner principle and a practical next step, not a technical lecture or a fabricated message to Niki. Ask for **"the full Niki message"** to confirm it reads `.hermes/skills/ask-niki/references/message-to-rado.md`. If skill discovery requires a reload/restart, do it only for the isolated Rado profile and verify the main Hermes profile remains unaffected.
