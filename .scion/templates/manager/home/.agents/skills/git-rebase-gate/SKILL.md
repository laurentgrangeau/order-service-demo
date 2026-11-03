---
name: git-rebase-gate
description: Non-interactive Git branch verification, rebasing, and release integration gate for Manager agents.
---

# Git Rebase & Integration Gate

1. Fetch and inspect candidate worker branches (`scion/sec-patch`, `scion/grpc-mig`, `scion/qa-chaos`).
2. Create or checkout the unified release branch:
   ```bash
   git checkout -B scion/release-v2.4
   ```
3. Merge or rebase each worker branch non-interactively (`GIT_EDITOR=true`).
4. Run the verification suite before reporting completion:
   ```bash
   go test -race -v ./...
   ```
