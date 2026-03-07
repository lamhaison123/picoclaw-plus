#!/bin/bash
# Diagnostic script to check team loading environment

echo "=== Environment Check ==="
echo "HOME: $HOME"
echo "USER: $USER"
echo "PWD: $PWD"
echo ""

echo "=== Workspace Path ==="
WORKSPACE="${HOME}/.picoclaw/workspace"
echo "Workspace: $WORKSPACE"
echo ""

echo "=== Team State Directory ==="
TEAM_DIR="${WORKSPACE}/teams/active"
echo "Team directory: $TEAM_DIR"
echo ""

if [ -d "$TEAM_DIR" ]; then
    echo "✓ Team directory exists"
    echo ""
    echo "=== Team State Files ==="
    ls -lh "$TEAM_DIR"/*.json 2>/dev/null || echo "No .json files found"
    echo ""
    
    if [ -f "$TEAM_DIR/dev-team.json" ]; then
        echo "✓ dev-team.json exists"
        echo ""
        echo "=== dev-team.json Content (first 20 lines) ==="
        head -20 "$TEAM_DIR/dev-team.json"
    else
        echo "✗ dev-team.json NOT found"
    fi
else
    echo "✗ Team directory does NOT exist"
fi

echo ""
echo "=== Permissions ==="
ls -ld "$HOME/.picoclaw" 2>/dev/null || echo ".picoclaw directory not found"
ls -ld "$WORKSPACE" 2>/dev/null || echo "workspace directory not found"
ls -ld "$TEAM_DIR" 2>/dev/null || echo "teams/active directory not found"
