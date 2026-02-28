# Contributing Guidelines

This repository is strictly proprietary and commercial. Contributions are restricted to authorized personnel only.

## Getting Started

1. Ensure Go `1.26.0` or higher is installed.
2. Clone the repository and run `make install`.
3. Read the `ROADMAP.md` before starting work to understand the current iteration logic.

## Development Workflow

- Create a feature branch from `master`.
- Write or update tests in `scripts/test-trader-e2e.sh` and `scripts/test-broker-e2e.sh`.
- Run `make lint && make test` before committing. The CI will strictly block any code-smells or test failures.
- Submit a Pull Request for review.
