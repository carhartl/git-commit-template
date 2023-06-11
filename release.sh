#!/bin/sh

set -eu

versionregex='v[[:digit:]]+\.[[:digit:]]+\.[[:digit:]]+'

git stash --include-untracked
grep -rl -E "$versionregex" --exclude-dir=.git --exclude-dir=.github --exclude-dir=.vscode --exclude-dir=dist --exclude=go.mod --exclude=go.sum . | xargs sed -i '' -E 's/'"$versionregex"'/'"v$1"'/g'
git add .
git commit -m "Update for v$1"
git tag "v$1" -m "Release v$1"
git push
goreleaser release --clean
git stash pop --index
