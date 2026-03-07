#!/bin/bash
# Build script for WASM plugins

set -e

echo "🔨 Building WASM Plugins..."

# Check if TinyGo is installed
if ! command -v tinygo &> /dev/null; then
    echo "❌ TinyGo is not installed"
    echo "Please install TinyGo: https://tinygo.org/getting-started/install/"
    exit 1
fi

# Create output directory
mkdir -p build/plugins

# Build math plugin
echo "Building math_plugin.wasm..."
tinygo build -o build/plugins/math_plugin.wasm \
    -target=wasi \
    -no-debug \
    -opt=2 \
    examples/plugins/math_plugin.go

echo "✅ math_plugin.wasm built successfully"

# Show file size
ls -lh build/plugins/math_plugin.wasm

echo ""
echo "🎉 All plugins built successfully!"
echo ""
echo "To test the plugin, run:"
echo "  go run examples/plugins/plugin_example.go build/plugins/math_plugin.wasm"
