# git-commit-template

[![CI](https://github.com/carhartl/git-commit-template/actions/workflows/ci.yml/badge.svg)](https://github.com/carhartl/git-commit-template/actions/workflows/ci.yml)

CLI utility for managing a Git commit message template within the current project.

- Optional issue reference
- Optional Co-authored-by line

## Installation

With [Go](https://golang.org/):

```bash
go install github.com/carhartl/git-commit-template
```

Via [Homebrew](https://brew.sh/):

```bash
brew install carhartl/tap/git-commit-template
```

## Configuration

| Env var                            | default |
| ---------------------------------- | ------- |
| `GIT_COMMIT_TEMPLATE_ISSUE_PREFIX` | `#`     |

## Releasing

```bash
./release.sh 0.0.2
```

## Also see

https://cbea.ms/git-commit/
