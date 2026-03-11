#!/usr/bin/env bash
set -euo pipefail

# 1) Toolchain
export PATH=/root/.go/bin:$PATH

# 2) Writable workspace (avoid read-only /root/.cache/go-build)
export GOPATH=/tmp/go
export GOCACHE=/tmp/go-cache
export GOMODCACHE=/tmp/go-mod
export GOFLAGS="-modcacherw"

# 3) Stable module resolution
export GOPROXY="https://proxy.golang.org,direct"
export GOSUMDB="sum.golang.org"

mkdir -p "$GOPATH" "$GOCACHE" "$GOMODCACHE"

echo "== go env =="
go version
go env GOPATH GOCACHE GOMODCACHE GOPROXY GOSUMDB GOFLAGS

echo ""
echo "== download deps =="
go mod download
go mod verify

echo ""
echo "== runtime tests =="
go test -v ./pkg/collaborative -run TestEnhancedMetrics
go test -race ./pkg/collaborative
