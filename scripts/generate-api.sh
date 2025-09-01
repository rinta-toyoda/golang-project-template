#!/bin/bash

set -e

echo "Generating API code from OpenAPI specifications..."

# Generate auth API
openapi-generator generate \
  -i api/auth/openapi.yaml \
  -g go-gin-server \
  -o ./gen/openapi/auth \
  --global-property=apis,models,apiTests=false,apiDocs=false,modelTests=false,modelDocs=false \
  --additional-properties=packageName=authapi,enumClassPrefix=true,outputAsLibrary=true

# Generate v1 API (when ready)
if [ -f "api/v1/openapi.yaml" ]; then
  openapi-generator generate \
    -i api/v1/openapi.yaml \
    -g go-gin-server \
    -o ./gen/openapi/v1 \
    --global-property=apis,models,apiTests=false,apiDocs=false,modelTests=false,modelDocs=false \
    --additional-properties=packageName=v1api,enumClassPrefix=true,outputAsLibrary=true
fi

# Format generated files
echo "Formatting generated Go files..."
gofmt -w ./gen/openapi/auth/go/*.go 2>/dev/null || true
gofmt -w ./gen/openapi/v1/go/*.go 2>/dev/null || true

echo "API generation completed successfully"