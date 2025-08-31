#!/bin/bash

set -e

echo "🔍 Running CI checks locally..."

echo "📋 Step 1: Go fmt check"
if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
    echo "❌ The following files are not formatted:"
    gofmt -s -l .
    echo "Run 'go fmt ./...' to fix formatting issues"
    exit 1
fi
echo "✅ Go fmt check passed"

echo "📋 Step 2: Go vet"
go vet ./...
echo "✅ Go vet passed"

echo "📋 Step 3: Go mod verify"
go mod verify
echo "✅ Go mod verify passed"

echo "📋 Step 4: Check if golangci-lint is available"
if command -v golangci-lint &> /dev/null; then
    echo "📋 Step 5: Running golangci-lint"
    golangci-lint run
    echo "✅ golangci-lint passed"
else
    echo "⚠️  golangci-lint not installed, skipping lint checks"
    echo "   Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
fi

echo "📋 Step 6: Check if gosec is available"
if command -v gosec &> /dev/null; then
    echo "Running gosec security scanner"
    gosec ./...
    echo "✅ gosec security scan passed"
else
    echo "⚠️  gosec not installed, skipping security scan"
    echo "   Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"
fi

echo "📋 Step 7: Running tests"
export CSRF_SECRET="test-csrf-secret"
export SESSION_SECRET="test-session-secret"
export DATABASE_URL="postgres://postgres:password@localhost:5432/app_db?sslmode=disable"

# Check if database is available
if pg_isready -h localhost -p 5432 -U postgres &> /dev/null; then
    echo "Running tests with database"
    go test -race -coverprofile=coverage.out -covermode=atomic ./...
    echo "✅ Tests passed"
    
    if [ -f coverage.out ]; then
        echo "📊 Coverage report:"
        go tool cover -func=coverage.out | tail -1
    fi
else
    echo "⚠️  PostgreSQL not available, running tests without database"
    go test ./internal/domain/... ./pkg/...
    echo "✅ Unit tests passed (integration tests skipped)"
fi

echo "🎉 All CI checks completed successfully!"