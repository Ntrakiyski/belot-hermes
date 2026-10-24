#!/usr/bin/env bash
# Requires a blank GitHub repo created by its owner, and local Git credentials.
set -euo pipefail
if [[ $# -ne 1 ]]; then
  echo 'Usage: bash scripts/push-to-github.sh https://github.com/OWNER/belot-hermes.git' >&2
  exit 2
fi
if ! git show-ref --verify --quiet refs/heads/main || ! git show-ref --verify --quiet refs/heads/demo; then
  echo 'Expected main and demo git branches. Use the git bundle or initialize both branches first.' >&2
  exit 1
fi
remote="$1"
if git remote get-url origin >/dev/null 2>&1; then git remote set-url origin "$remote"; else git remote add origin "$remote"; fi
git push -u origin main demo
