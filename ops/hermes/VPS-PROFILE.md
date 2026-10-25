# belot-rado — operations (no secrets)
- Profile: `~/.hermes/profiles/belot-rado` (wrapper `belot-rado`). Own gateway: `gateway.standalone: true`, user service `hermes-gateway-belot-rado`.
- Telegram: @Rado_agent_jp_bot; token only in profile `.env`; allowlist = Rado + Niki.
- Sandbox: docker backend, image `belot-rado-sandbox:1`, only `/srv/rado` mounted as `/workspace`, 2 CPU / 4 GB. Repo at `/workspace/projects/belot`. Host symlink `/workspace -> /srv/rado` exists only so Hermes can validate cwd.
- GitHub: repo-only SSH deploy key in `/srv/rado/.ssh`. `main`: PR required, no force push/delete, admins enforced.
- DEMO: Cloudflare Pages `belot-demo` -> https://belot-demo.trakiyski.work. Rado's Pages-edit token is account-wide for Pages (Cloudflare can't scope per project); instructions limit it to `belot-demo`.
- Deploy: `wrangler pages deploy <dir> --project-name belot-demo --branch demo`.
- Persona sync: after editing `ops/hermes/SOUL.md`, the operator runs `ops/hermes/sync-profile.sh`, then restarts the service. The agent cannot edit its live profile (sandbox).
- Skills: `.hermes/skills/` trusted via `belot-rado skills trust /srv/rado/projects/belot`, plus `find-skills`.
