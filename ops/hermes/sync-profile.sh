#!/usr/bin/env bash
# Operator-only (not the agent): copy ops/hermes/SOUL.md + appendix into the live profile.
set -euo pipefail
P="${HERMES_PROFILE_DIR:-$HOME/.hermes/profiles/belot-rado}"
SRC="${1:-/srv/rado/projects/belot/ops/hermes}"
# Validate and assemble before replacing the live persona.
test -r "$SRC/SOUL.md" && test -r "$SRC/SOUL-profile-appendix.md"
# Check every managed skill before changing any live file.
for skill in "$SRC"/profile-skills/*; do
  [[ -d "$skill" ]] || continue
  test -r "$skill/SKILL.md"
done
TMP=$(mktemp "$P/SOUL.md.tmp.XXXXXX")
SKILL_TMP=""
trap 'rm -f "$TMP"; [[ -z "$SKILL_TMP" ]] || rm -f "$SKILL_TMP"' EXIT
cat "$SRC/SOUL.md" "$SRC/SOUL-profile-appendix.md" > "$TMP"
cp "$P/SOUL.md" "$P/SOUL.md.bak"
mv "$TMP" "$P/SOUL.md"
for skill in "$SRC"/profile-skills/*; do
  [[ -d "$skill" ]] || continue
  target="$P/skills/$(basename "$skill")"
  mkdir -p "$target"
  [[ ! -f "$target/SKILL.md" ]] || cp "$target/SKILL.md" "$target/SKILL.md.bak"
  SKILL_TMP=$(mktemp "$target/SKILL.md.tmp.XXXXXX")
  cp "$skill/SKILL.md" "$SKILL_TMP"
  mv "$SKILL_TMP" "$target/SKILL.md"
  SKILL_TMP=""
done
echo "SOUL.md and managed profile skill bridges synced. Then run: systemctl --user restart hermes-gateway-belot-rado"
