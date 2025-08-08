#!/bin/bash
# BMAD Contract Verification Gate
# This script must pass for any agent work to be accepted

set -euo pipefail

echo "🔧 LUCIEN CLI - VERIFICATION GATE"
echo "================================="

# Change to project root
cd "$(dirname "$0")/.."

# Ensure build directory exists
mkdir -p build

echo "📦 Step 1: Go Module Verification"
go mod tidy
go mod verify

echo "🔍 Step 2: Static Analysis"
if ! go vet ./...; then
    echo "❌ go vet failed - code has issues"
    exit 1
fi

echo "🧪 Step 3: Unit Tests"
if [[ -n "${CI:-}" ]]; then
    # CI gets full coverage and race detection
    if ! go test -v -race -coverprofile=coverage.out ./...; then
        echo "❌ Unit tests or race detector failed"
        exit 1
    fi
    echo "📊 Coverage report generated: coverage.out"
else
    # Local dev gets fast feedback
    if ! go test -v ./...; then
        echo "❌ Unit tests failed"
        exit 1
    fi
fi

echo "🏗️ Step 4: Build Verification"
if ! go build -o build/lucien-verify cmd/lucien/main.go; then
    echo "❌ Build failed"
    exit 1
fi

echo "🚀 Step 5: Smoke Test"
SMOKE_OUTPUT=$(timeout 10s ./build/lucien-verify --help 2>&1)
SMOKE_EXIT=$?
if [ $SMOKE_EXIT -eq 0 ] && echo "$SMOKE_OUTPUT" | grep -q "Usage of"; then
    echo "✅ Smoke test passed - binary runs correctly"
else
    echo "❌ Smoke test failed - binary doesn't run properly"
    echo "Exit code: $SMOKE_EXIT"
    echo "Output: $SMOKE_OUTPUT"
    exit 1
fi

# Clean up temp binary
rm -f build/lucien-verify

echo "✅ VERIFICATION COMPLETE - All gates passed"
echo "🎯 Ready for next BMAD phase"