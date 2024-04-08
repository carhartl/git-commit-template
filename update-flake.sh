#!/bin/sh

set -eu

grep _ ./tools.go | awk -F'"' '{print $2}' | xargs -tI % go install %
OUT=$(mktemp -d -t nar-hash-XXXXXX)
rm -rf "$OUT"
go mod vendor -o "$OUT"
go run tailscale.com/cmd/nardump --sri "$OUT" >go.mod.sri
rm -rf "$OUT"
