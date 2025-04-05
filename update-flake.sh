#!/bin/sh

set -eu

OUT=$(mktemp -d -t nar-hash-XXXXXX)
rm -rf "$OUT"
go mod vendor -o "$OUT"
go tool nardump --sri "$OUT" >go.mod.sri
rm -rf "$OUT"
