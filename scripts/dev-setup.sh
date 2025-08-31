#!/bin/bash

set -e

echo "Setting up development environment..."

# Create .env if it doesn't exist
if [ ! -f .env ]; then
    echo "Creating .env file from template..."
    cp .env.example .env
    echo "Please edit .env file with your configuration"
fi

# Create necessary directories
mkdir -p build/{bin,tmp}
mkdir -p docs

# Install lefthook if npm is available
if command -v npm &> /dev/null; then
    echo "Installing lefthook..."
    npx lefthook install
else
    echo "npm not found, skipping lefthook installation"
fi

# Download Go modules
echo "Downloading Go modules..."
go mod download

echo "Development environment setup complete!"
echo "Run 'task up' to start the services"