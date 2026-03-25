#!/bin/sh

set -eu

git tag "v$1" -m "Release v$1"
git push --tags
