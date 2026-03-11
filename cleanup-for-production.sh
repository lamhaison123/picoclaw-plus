#!/bin/bash
# Cleanup script for production release
# Moves temporary integration docs to archive

set -e

echo "🧹 Cleaning up temporary files for production..."

# Create archive directory
mkdir -p .archive/integration-docs

# Move temporary integration documentation
echo "📦 Archiving integration documentation..."
mv -f *COMPLETE*.md .archive/integration-docs/ 2>/dev/null || true
mv -f *INTEGRATION*.md .archive/integration-docs/ 2>/dev/null || true
mv -f *SUMMARY*.md .archive/integration-docs/ 2>/dev/null || true
mv -f *ANALYSIS*.md .archive/integration-docs/ 2>/dev/null || true
mv -f *STATUS*.md .archive/integration-docs/ 2>/dev/null || true
mv -f *CONFLICT*.md .archive/integration-docs/ 2>/dev/null || true
mv -f DOCS_REORGANIZATION.md .archive/integration-docs/ 2>/dev/null || true
mv -f REMAINING_INTEGRATION_OPPORTUNITIES.md .archive/integration-docs/ 2>/dev/null || true
mv -f QUICK_WINS_COMPLETE.md .archive/integration-docs/ 2>/dev/null || true

# Keep important docs in root
echo "✅ Keeping essential documentation in root..."
# These stay in root:
# - CHANGELOG.md
# - README.md
# - CONTRIBUTING.md
# - LICENSE
# - .env.example
# - .gitignore

# Keep test reports but move to docs
echo "📊 Moving test reports to docs..."
mkdir -p docs/testing
mv -f TEST_RESULTS_V0.2.1.md docs/testing/ 2>/dev/null || true
mv -f COMPREHENSIVE_TEST_REPORT.md docs/testing/ 2>/dev/null || true
mv -f VERIFICATION_CHECKLIST.md docs/testing/ 2>/dev/null || true

# Keep quick reference but move to docs
echo "📚 Moving reference docs..."
mv -f QUICK_REFERENCE_v0.2.1.md docs/reference/ 2>/dev/null || true
mv -f SEARCH_PROVIDERS_GUIDE.md docs/reference/ 2>/dev/null || true

# Keep release notes but move to docs
echo "📝 Moving release notes..."
mkdir -p docs/releases
mv -f RELEASE_NOTES_v0.2.1_INTEGRATION.md docs/releases/RELEASE_NOTES_v0.2.1.md 2>/dev/null || true

# Clean up build artifacts
echo "🗑️  Cleaning build artifacts..."
rm -rf build/*.test 2>/dev/null || true
rm -rf *.test 2>/dev/null || true

# Update .gitignore
echo "📋 Updating .gitignore..."
cat >> .gitignore << 'EOF'

# Archive directory (integration docs)
.archive/

# Test artifacts
*.test
coverage.out
*.prof

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Local config
config.json
auth.json
.env
.env.local

# Runtime
*.pid
*.sock
*.log

# Build
build/
dist/
EOF

echo "✨ Cleanup complete!"
echo ""
echo "📁 File organization:"
echo "  - Integration docs → .archive/integration-docs/"
echo "  - Test reports → docs/testing/"
echo "  - Reference docs → docs/reference/"
echo "  - Release notes → docs/releases/"
echo ""
echo "✅ Repository is ready for production!"
