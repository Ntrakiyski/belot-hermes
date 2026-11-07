# Belot deployment and Hermes operations

## Scope and current evidence (2026-11-07)

This is an engineering DEMO, not the authenticated four-human product. The application keeps tables in memory: **every container replacement/restart loses active tables**. Neither this runbook nor a healthy process certifies Belot rule correctness.

## Current deployment and operating path

- DEMO uses Git Docker Compose, manifest `/docker-compose.yaml`, service `belot` on internal port3000. v0.3.0 was verified at merge `2998b7a4433b1151bad72e5a4a4a9477d404f879`; PR#6 and its deployment evidence are in GitHub.
- `https://belot-demo.trakiyski.work` serves UI/API over HTTPS. A tested PR merged into `demo` automatically triggers Coolify. Wait for all four named CI checks before authorized integration, then verify live version and affected gameplay. No new hosting credential or manual Pages deployment is needed.
- Hermes has existing configured GitHub access, confirmed by the owner. Use that integration; an absent CLI is not evidence of missing access. Keep the Docker socket out of its sandbox. See `.hermes/skills/belot-maintenance/SKILL.md` for the current no-Docker test path.
- Dedicated profile/service: `belot-rado`, `hermes-gateway-belot-rado`, Telegram `@Rado_agent_jp_bot`. Its sandbox mounts `/srv/rado` as `/workspace`, repo `/workspace/projects/belot`; skills must remain trusted and discoverable. Operator-owned profile bridges point at the maintained checkout instead of storing historical game/hosting facts.
- MAIN still requires explicit version approval and the documented human-review gate. Earlier protection inspection found PR enforcement/no force-push/delete but no required reviewer/CI count; changing protection is a separate authorized action. Do not infer that MAIN is automatically safe to ship.

Historical findings (Dockerfile mode,503routing and stale Pages persona) motivated the migration; they are not current deployment instructions. The optional `www` alias had no DNS record in the prior audit; use the verified canonical domain above.

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

The initial Compose migration is complete. The following is retained for an operator recreating the resource or rolling back its configuration; normal changes use the automatic Git merge path above. Do not promote to MAIN as part of infrastructure maintenance.

1. Wait for testers to finish. Record the old Dockerfile configuration, domains, runtime variables and last successful deployment/image. Confirm a usable rollback image. Resolve outstanding feature work without resetting the remote Hermes checkout.
2. After the Compose commit reaches `demo`, change this existing Coolify Git application's build strategy to **Docker Compose**, base directory `/`, Compose location `/docker-compose.yaml`. Load the repository definition, service `belot`. Do not use “Docker Compose Empty”, which loses Git/branch automation.
3. Assign the intended HTTPS domains to **service `belot` at port 3000**: `https://belot-demo.trakiyski.work:3000` and `https://www.belot-demo.trakiyski.work:3000`. Preserve the temporary domain only if still needed. Port 3000 is an internal proxy target; users open normal HTTPS. Use only the deployment manifest, not the local override. Keep raw Compose deployment off so Coolify manages routing.
4. Review runtime `ENABLE_PRACTICE_BOTS` explicitly. Default `0`; allow `1` only for clearly labelled DEMO practice. Keep it `0` on any future MAIN deployment. Configure secrets privately in Coolify, never in Git. Deploy the intended commit; wait for container health and successful deployment logs.
5. Verify normal TLS (no insecure bypass), both domain routes, React assets, `/healthz`, and an authorized disposable table flow. Match Coolify's deployed commit to the approved commit. Record evidence before saying “live”. `/healthz` is process health, not multiplayer readiness; `practiceEnabled` reports the runtime practice flag and `version` reports the app release; neither proves a specific commit by itself.

Rollback: select the recorded old build strategy/configuration and known-good deployment/image through Coolify; verify UI/API/TLS again. A rollback cannot restore tables already lost. Do not delete volumes, unrelated apps or DNS as a rollback step.

Official reference: [Coolify Git Docker Compose deployment](https://coolify.io/docs/applications/builds/docker-compose).

## Operator sync for the dedicated Hermes profile

After the reviewed docs reach the trusted VPS checkout, preserve remote untracked work, then use `ops/hermes/sync-profile.sh` to back up and replace **only** the Rado profile's SOUL and managed skill bridges with the updated sources. Restart only `hermes-gateway-belot-rado` when idle and with operator authorization. Other profiles must remain untouched. The script does not update config, credentials, skill trust or gateway routing. Preserve custom work and archive duplicate historical bridge skills outside discovery rather than keeping conflicting copies.

Stop using `wrangler pages deploy` for Belot. Pages cannot host this Go runtime. Keep the agent focused on tested feature PRs and DEMO evidence; an operator or tightly scoped integration handles Coolify. Do not add a broad Coolify token to Hermes just to automate deployment.

## Remaining acceptance gates

- **Telegram:** `getMe` must identify the intended bot, then Rado verifies an ordinary text request through his authorized session. Voice support is explicitly outside the requested scope. A separately controlled unauthorized account must be rejected. Inspect sanitized results, not private conversation histories. Do not consume Telegram updates with `getUpdates` while Hermes is polling.
- **Skills and workspace:** verify project skills are trusted and loaded, including `ask-niki`; confirm ordinary requests read the current repo, use the correct cwd and choose relevant skills without slash commands. Presence on disk alone does not establish discovery. Prove sandbox mount/permission boundaries without probing other agents' data.
- **GitHub delivery:** use the existing configured integration to inspect the repo, create the feature PR and observe CI. Access is owner-confirmed; retain its scope instead of creating replacements. Record the actual delivery evidence for each change.
- **Release gates:** require human review and named passing CI checks (`go-engine`, `react-vite`, `legacy-policy`, and `compose-smoke`) on MAIN; retain admin enforcement and no force-push/deletion. Ensure DEMO integration is CI-gated, since auto-deploy alone is not a test gate. Scope Hermes away from MAIN release authority. Changing protections requires operator approval.
- **Phone acceptance:** Rado sees a real deployed feature, verifies it, and explicitly names the tested version when approving MAIN. Current guest seats, polling and practice bots do not satisfy account identity, public matchmaking, reconnect handling or rule-inspector requirements.

A harmless Telegram test for Rado: “Read the Belot repo and tell me what is playable now and what is still missing. Make no changes.” This checks basic text routing; it does not prove the feature-to-DEMO release loop.
