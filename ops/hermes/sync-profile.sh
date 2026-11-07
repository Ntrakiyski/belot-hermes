#!/usr/bin/env bash
# Operator-only (not the agent): copy ops/hermes/SOUL.md + appendix into the live profile.
set -euo pipefail
P="${HERMES_PROFILE_DIR:-$HOME/.hermes/profiles/belot-rado}"
SRC="${1:-/srv/rado/projects/belot/ops/hermes}"
# Validate and assemble before replacing the live persona.
test -r "$SRC/SOUL.md" && test -r "$SRC/SOUL-profile-appendix.md"
TMP=$(mktemp "$P/SOUL.md.tmp.XXXXXX")
trap 'rm -f "$TMP"' EXIT
cat "$SRC/SOUL.md" "$SRC/SOUL-profile-appendix.md" > "$TMP"
cp "$P/SOUL.md" "$P/SOUL.md.bak"
mv "$TMP" "$P/SOUL.md"
echo "SOUL.md synced. Then run: systemctl --user restart hermes-gateway-belot-rado"
