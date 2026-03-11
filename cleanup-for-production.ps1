# Cleanup script for production release
# Moves temporary integration docs to archive

Write-Host "🧹 Cleaning up temporary files for production..." -ForegroundColor Cyan

# Create archive directory
New-Item -ItemType Directory -Force -Path ".archive/integration-docs" | Out-Null

# Move temporary integration documentation
Write-Host "📦 Archiving integration documentation..." -ForegroundColor Yellow

$tempDocs = @(
    "*COMPLETE*.md",
    "*INTEGRATION*.md", 
    "*SUMMARY*.md",
    "*ANALYSIS*.md",
    "*STATUS*.md",
    "*CONFLICT*.md",
    "DOCS_REORGANIZATION.md",
    "REMAINING_INTEGRATION_OPPORTUNITIES.md",
    "QUICK_WINS_COMPLETE.md"
)

foreach ($pattern in $tempDocs) {
    Get-ChildItem -Path . -Filter $pattern -File -ErrorAction SilentlyContinue | 
        Move-Item -Destination ".archive/integration-docs/" -Force -ErrorAction SilentlyContinue
}

# Keep test reports but move to docs
Write-Host "📊 Moving test reports to docs..." -ForegroundColor Yellow
New-Item -ItemType Directory -Force -Path "docs/testing" | Out-Null

$testReports = @(
    "TEST_RESULTS_V0.2.1.md",
    "COMPREHENSIVE_TEST_REPORT.md",
    "VERIFICATION_CHECKLIST.md"
)

foreach ($file in $testReports) {
    if (Test-Path $file) {
        Move-Item -Path $file -Destination "docs/testing/" -Force
    }
}

# Keep quick reference but move to docs
Write-Host "📚 Moving reference docs..." -ForegroundColor Yellow

$refDocs = @(
    "QUICK_REFERENCE_v0.2.1.md",
    "SEARCH_PROVIDERS_GUIDE.md"
)

foreach ($file in $refDocs) {
    if (Test-Path $file) {
        Move-Item -Path $file -Destination "docs/reference/" -Force
    }
}

# Keep release notes but move to docs
Write-Host "📝 Moving release notes..." -ForegroundColor Yellow
New-Item -ItemType Directory -Force -Path "docs/releases" | Out-Null

if (Test-Path "RELEASE_NOTES_v0.2.1_INTEGRATION.md") {
    Move-Item -Path "RELEASE_NOTES_v0.2.1_INTEGRATION.md" -Destination "docs/releases/RELEASE_NOTES_v0.2.1.md" -Force
}

# Clean up build artifacts
Write-Host "🗑️  Cleaning build artifacts..." -ForegroundColor Yellow
Remove-Item -Path "build/*.test" -Force -ErrorAction SilentlyContinue
Remove-Item -Path "*.test" -Force -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "✨ Cleanup complete!" -ForegroundColor Green
Write-Host ""
Write-Host "📁 File organization:" -ForegroundColor Cyan
Write-Host "  - Integration docs → .archive/integration-docs/"
Write-Host "  - Test reports → docs/testing/"
Write-Host "  - Reference docs → docs/reference/"
Write-Host "  - Release notes → docs/releases/"
Write-Host ""
Write-Host "✅ Repository is ready for production!" -ForegroundColor Green
