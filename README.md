# git-commit-template

[![CI](https://github.com/carhartl/git-commit-template/actions/workflows/ci.yml/badge.svg)](https://github.com/carhartl/git-commit-template/actions/workflows/ci.yml)

Command-line utility for managing a Git commit message template within the current project.

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

## Usage

Pass issue ref numbers with or without prefix:

```bash
git-commit-template set --issue 123 --dry-run
git-commit-template set --issue #123 --dry-run
```

Produces:

```
Subject (keep under 50 characters)

Context/description (what and why)

Addresses: #123
```

Pass co-author with fuzzy matching entries in author file:

```bash
echo 'John Doe <john@example.com>' >> $HOME/.git-commit-template-authors
echo 'Mary Tester <mary@example.com>' >> $HOME/.git-commit-template-authors
git-commit-template set --pair doe --dry-run
```

Produces:

```
Subject (keep under 50 characters)

Context/description (what and why)

Co-authored-by: John Doe <john@example.com>
```

Multiple co-authors:

```bash
git-commit-template set --pair doe --pair mary --dry-run
```

Produces:

```
Subject (keep under 50 characters)

Context/description (what and why)

Co-authored-by: John Doe <john@example.com>
Co-authored-by: Mary Tester <mary@example.com>
```

Tip:

Due to the executable's naming you can also call it like so:

```bash
git commit-template ...
```

## Configuration

| Env var                            | Default                              |
| ---------------------------------- | ------------------------------------ |
| `GIT_COMMIT_TEMPLATE_ISSUE_PREFIX` | `#`                                  |
| `GIT_COMMIT_TEMPLATE_AUTHOR_FILE`  | `$HOME/.git-commit-template-authors` |

## Releasing

```bash
op run --env-file .env -- ./release.sh 0.1.0
```

## Also see

https://cbea.ms/git-commit/
