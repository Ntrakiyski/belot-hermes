---
name: skill-discovery
description: Find existing skills for unfamiliar work using Find Skills and Hermes skill catalogs, and vet before installing.
---
# Discover missing capabilities responsibly
1. First search current Hermes skills (including project-local and bundled). No unnecessary new dependency if existing capability fits.
2. Search `https://skills.sh` and/or `npx skills find <domain>` for relevant task, then compare authors, popularity, version and scripts. The official skill page is https://www.skills.sh/vercel-labs/skills/find-skills .
3. Recommend at most 1–2 options with a concrete benefit. **Inspect SKILL.md and executable scripts before installation.** No automatic installation of arbitrary network-retrieved code, new credentials, excessive permissions or unexpected infrastructure spending.
4. Request operator authorization for an untrusted installation/execution; use isolated profile environment. Record origin, version, permissions, rationale and rollback. Re-scan on updates.
5. After installation, load it only for relevant tasks, and keep project business/security invariants from `AGENTS.md` authoritative.

Find Skills installation command (on the VPS, in a reviewed target environment): `npx skills add https://github.com/vercel-labs/skills --skill find-skills`. Hermes' own project skill discovery relies on `.hermes/skills` and a trusted git checkout; the skills CLI may not natively target Hermes, so check its install target or use `npx skills find` from the agent tool shell.
