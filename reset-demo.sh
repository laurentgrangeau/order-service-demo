#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

echo "==> Stopping and deleting Scion demo agents (if running)..."
if command -v scion >/dev/null 2>&1; then
  for agent in lead-mgr sec-patch grpc-mig qa-chaos; do
    scion delete "$agent" 2>/dev/null || true
  done
fi

echo "==> Cleaning up local .scion/agents/ state..."
rm -rf .scion/agents/

echo "==> Removing any leftover scion/* git branches..."
for br in $(git branch --list 'scion/*' | sed 's/^[* ]*//'); do
  git branch -D "$br" || true
done

echo "==> Verifying baseline Go build & tests..."
go test -race ./...

echo "==> Demo reset complete! Ready for live presentation."
