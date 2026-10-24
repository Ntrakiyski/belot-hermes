#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
printf 'Belot workspace checks (no live deployment performed)\n'
command -v git >/dev/null || { echo 'Git missing'; exit 1; }
command -v go >/dev/null || { echo 'Go missing'; exit 1; }
command -v node >/dev/null || { echo 'Node missing'; exit 1; }
(cd services/game && go test ./...)
node --test tests/*.test.mjs
if [ -d apps/web/node_modules ]; then (cd apps/web && npm run build); else echo 'React build SKIPPED: npm packages not installed'; fi
printf 'Workspace checks passed (skipped tasks above do not count as verified).\n'
