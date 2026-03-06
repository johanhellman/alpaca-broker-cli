---
description: How to handle Dependabot or automated dependency PRs safely
---
When asked to review, modify, or merge multiple Dependabot (or automated dependency update) Pull Requests, follow these steps strictly to avoid resolving complex Git merge conflicts and desyncing branches:

1. **List open PRs**: Run `gh pr list` to view all available dependency updates.
2. **Review the PRs**: Understand the scope of the updates to ensure they are standard semantic version bumps and the CI tests are passing.
3. **Merge sequentially, allowing auto-rebase**:
   - Merge **ONLY ONE** PR at a time (e.g., using `gh pr merge <PR_NUMBER> --squash --delete-branch`). Follow this by `git pull origin master` locally.
   - Wait for Dependabot to automatically rebase its remaining open branches against the new `master` branch.
   - You can also forcefully trigger a rebase on the remaining PRs by commenting `@dependabot rebase` on them using `gh pr comment <PR_NUMBER> --body "@dependabot rebase"`.
4. **DO NOT** attempt to checkout all Dependabot branches locally and merge them manually at the same time. Since they often modify the exact same files (e.g. `go.mod` or `.github/workflows/ci.yml`), doing so will result in painful merge conflicts and require overriding GitHub branch protections to resolve.
