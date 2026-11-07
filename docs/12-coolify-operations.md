# Belot deployment and Hermes operations

## Scope and current evidence (2026-11-07)

This is an engineering DEMO, not the authenticated four-human product. The application keeps tables in memory: **every container replacement/restart loses active tables**. Neither this runbook nor a healthy process certifies Belot rule correctness.

Read-only checks found:

- Coolify application `lzxn1arb4v7hwmlaswjxizip`, Git source `Ntrakiyski/belot-hermes`, branch `demo`, Dockerfile build, running healthy at deployed commit `5be67819515fa1783880716542a493e132b208b8`. Automatic deployment is enabled; no CI-wait gate was proven.
- Temporary HTTP origin serves React and `/healthz`. The intended `https://belot-demo.trakiyski.work` returns **503**. Running proxy labels contain only the temporary hostname, while Coolify's saved app settings include the custom domains. Apply/redeploy the intended routing and then verify TLS and both domains. Do not change DNS merely to fix stale container labels.
- The dedicated Hermes user service `hermes-gateway-belot-rado` is active on the Hermes VPS. Its Docker sandbox uses `belot-rado-sandbox:1`, `/srv/rado:/workspace`, cwd `/workspace/projects/belot`; Go, Node, npm and Git exist in the image. Docker and `gh` do not. Keep the Docker socket out of the sandbox.
- Telegram `getMe` identifies `Rado_agent_jp_bot`; Hermes 0.18.0 reports Telegram connected. A two-entry numeric allowlist is configured, unrestricted access is off, and the repo is in `trusted_project_dirs`. Live SOUL still says Cloudflare Pages. The remote checkout is on `feature/playable-engine` at `7a2b65b` with untracked work; preserve it.
- MAIN protection enforces PRs for administrators and rejects force-push/deletion, but has **zero required approvals and no required CI checks**. This does not enforce the Constitution's human-review/green-CI gate. No MAIN application was found in the inspected Coolify instance.

These are timestamped observations, not permanent guarantees. Profile configuration alone does not prove delivery or isolation. No Telegram message was sent, no live gateway restarted, no public change pushed and no hosting configuration modified during this audit.

## Local Compose

```bash
docker compose -f docker-compose.yaml -f docker-compose.local.yaml up --build -d --wait
npm run smoke:game
```

Open http://localhost:3000. The override publishes only `127.0.0.1:3000`; the deployment manifest itself exposes container port 3000 without publishing a host port. The existing Dockerfile builds React and Go together; Compose does not replace the image build.

Practice is opt-in and recreates the container:

```bash
ENABLE_PRACTICE_BOTS=1 docker compose -f docker-compose.yaml -f docker-compose.local.yaml up -d --wait
npm run smoke:game -- --practice
```

The smoke creates disposable tables and checks UI HTTP, health, authorization rejection, four-seat projections, full eight-trick play and next-deal flow. Use it only locally or against an explicitly authorized DEMO URL via `BELOT_SMOKE_URL`; never against MAIN. It does not prove four distinct authenticated humans or certify rules. Browser testing separately verifies rendering and interaction.

## Operator migration to Compose

Prepare and review the feature PR into `demo` first. The manifest is currently local until approved publication/integration. Do not promote to MAIN as part of this migration.

1. Wait for testers to finish. Record the old Dockerfile configuration, domains, runtime variables and last successful deployment/image. Confirm a usable rollback image. Resolve outstanding feature work without resetting the remote Hermes checkout.
2. After the Compose commit reaches `demo`, change this existing Coolify Git application's build strategy to **Docker Compose**, base directory `/`, Compose location `/docker-compose.yaml`. Load the repository definition, service `belot`. Do not use “Docker Compose Empty”, which loses Git/branch automation.
3. Assign the intended HTTPS domains to **service `belot` at port 3000**: `https://belot-demo.trakiyski.work:3000` and `https://www.belot-demo.trakiyski.work:3000`. Preserve the temporary domain only if still needed. Port 3000 is an internal proxy target; users open normal HTTPS. Use only the deployment manifest, not the local override. Keep raw Compose deployment off so Coolify manages routing.
4. Review runtime `ENABLE_PRACTICE_BOTS` explicitly. Default `0`; allow `1` only for clearly labelled DEMO practice. Keep it `0` on any future MAIN deployment. Configure secrets privately in Coolify, never in Git. Deploy the intended commit; wait for container health and successful deployment logs.
5. Verify normal TLS (no insecure bypass), both domain routes, React assets, `/healthz`, and an authorized disposable table flow. Match Coolify's deployed commit to the approved commit. Record evidence before saying “live”. `/healthz` is process health, not multiplayer readiness; its current `practiceTables` field is not a runtime bot-enable flag.

Rollback: select the recorded old build strategy/configuration and known-good deployment/image through Coolify; verify UI/API/TLS again. A rollback cannot restore tables already lost. Do not delete volumes, unrelated apps or DNS as a rollback step.

Official reference: [Coolify Git Docker Compose deployment](https://coolify.io/docs/applications/builds/docker-compose).

## Operator sync for the dedicated Hermes profile

After the reviewed docs reach the trusted VPS checkout, preserve remote untracked work, then use `ops/hermes/sync-profile.sh` to back up and replace **only** the Rado profile's SOUL with the updated persona and appendix. Restart only `hermes-gateway-belot-rado` when idle and with operator authorization. Other profiles must remain untouched. The script does not update config, credentials, skill trust or gateway routing.

Stop using `wrangler pages deploy` for Belot. Pages cannot host this Go runtime. Keep the agent focused on tested feature PRs and DEMO evidence; an operator or tightly scoped integration handles Coolify. Do not add a broad Coolify token to Hermes just to automate deployment.

## Remaining acceptance gates

- **Telegram:** `getMe` must identify the intended bot, then Rado verifies one text and one voice request through his authorized session. A separately controlled unauthorized account must be rejected. Inspect sanitized results, not private conversation histories. Do not consume Telegram updates with `getUpdates` while Hermes is polling.
- **Skills and workspace:** verify project skills are trusted and loaded, including `ask-niki`; confirm ordinary requests read the current repo, use the correct cwd and choose relevant skills without slash commands. Presence on disk alone does not establish discovery. Prove sandbox mount/permission boundaries without probing other agents' data.
- **GitHub delivery:** prove feature-branch push and a real PR into `demo` using repo-scoped API authentication. An SSH deploy key only proves Git access. No `gh` binary or profile `GH_TOKEN`/`GITHUB_TOKEN` was found; an alternative API integration remains unverified. Do not create broader credentials without approval.
- **Release gates:** require human review and named passing CI checks (`go-engine`, `react-vite`, `legacy-policy`, and the new `compose-smoke` after its first successful run) on MAIN; retain admin enforcement and no force-push/deletion. Ensure DEMO integration is CI-gated, since auto-deploy alone is not a test gate. Scope Hermes away from MAIN release authority. Changing protections requires operator approval.
- **Phone acceptance:** Rado sees a real deployed feature, verifies it, and explicitly names the tested version when approving MAIN. Current guest seats, polling and practice bots do not satisfy account identity, public matchmaking, reconnect handling or rule-inspector requirements.

A harmless Telegram test for Rado: “Read the Belot repo and tell me what is playable now and what is still missing. Make no changes.” Then send the same request as a voice note. This checks basic routing and transcription; it does not prove the feature-to-DEMO release loop.
