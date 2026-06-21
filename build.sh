#!/bin/bash

# Build the Go server
echo "Building the server..."
GOOS=linux GOARCH=amd64 go build -o bin/server-linux-amd64 . && echo "Server built successfully for Linux AMD64" || echo "Failed to build server for Linux AMD64"
GOOS=darwin GOARCH=amd64 go build -o bin/server-darwin-amd64 . && echo "Server built successfully for macOS AMD64" || echo "Failed to build server for macOS AMD64"
GOOS=windows GOARCH=amd64 go build -o bin/server-windows-amd64.exe . && echo "Server built successfully for Windows AMD64" || echo "Failed to build server for Windows AMD64"

GOOS=linux GOARCH=arm64 go build -o bin/server-linux-arm64 . && echo "Server built successfully for Linux ARM64" || echo "Failed to build server for Linux ARM64"
GOOS=darwin GOARCH=arm64 go build -o bin/server-darwin-arm64 . && echo "Server built successfully for macOS ARM64" || echo "Failed to build server for macOS ARM64"
GOOS=windows GOARCH=arm64 go build -o bin/server-windows-arm64.exe . && echo "Server built successfully for Windows ARM64" || echo "Failed to build server for Windows ARM64"
