# belot-rado — operations (no secrets)
- Profile: `~/.hermes/profiles/belot-rado` (wrapper `belot-rado`). Own gateway: `gateway.standalone: true`, user service `hermes-gateway-belot-rado`.
- Telegram: @Rado_agent_jp_bot; token only in profile `.env`; allowlist = Rado + Niki.
- Sandbox: docker backend, image `belot-rado-sandbox:1`, only `/srv/rado` mounted as `/workspace`, 2 CPU / 4 GB. Repo at `/workspace/projects/belot`. Host symlink `/workspace -> /srv/rado` exists only so Hermes can validate cwd.
- GitHub: existing repo-scoped integration (owner-confirmed); retain its configured credentials. `main`: PR required, no force push/delete, admins enforced.
- DEMO: Coolify Git application `belot-hermes:demo`, intended domain https://belot-demo.trakiyski.work. Combined Go + React image on port 3000; Cloudflare Pages is no longer the deployment path.
- Deploy: feature PR -> tests/review -> authorized integration into `demo` -> Coolify -> verify deployed commit, HTTPS UI and API. See `docs/12-coolify-operations.md` for Compose migration, gates and rollback. No broad hosting credential or Docker socket in the agent sandbox.
- Persona sync: the operator runs `ops/hermes/sync-profile.sh` to update SOUL and managed profile skill bridges, then restarts the dedicated service only when idle. The agent cannot edit its live profile (sandbox).
- Skills: `.hermes/skills/` trusted via `belot-rado skills trust /srv/rado/projects/belot`, plus `find-skills`.

## Verification status (2026-11-07)
Gateway service active; sandbox backend/image/mount/cwd match this profile; Telegram token present and a two-entry allowlist configured with unrestricted access off. These checks do not prove text handling, rejection of an unauthorized sender or PR creation. Checkout was safely fast-forwarded to the current DEMO while preserving its old feature branch and untracked work. Profile bridge/SOUL sync and actual text-channel checks are recorded in the current task evidence. Do not reset the checkout or infer delivery from config alone. See the operations runbook for outstanding release checks.
