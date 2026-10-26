#!/usr/bin/env bash
# Operator-only (not the agent): copy ops/hermes/SOUL.md + appendix into the live profile.
set -euo pipefail
P="${HERMES_PROFILE_DIR:-$HOME/.hermes/profiles/belot-rado}"
SRC="${1:-/srv/rado/projects/belot/ops/hermes}"
cp "$P/SOUL.md" "$P/SOUL.md.bak"
cat "$SRC/SOUL.md" "$SRC/SOUL-profile-appendix.md" > "$P/SOUL.md"
echo "SOUL.md synced. Then run: systemctl --user restart hermes-gateway-belot-rado"
