#!/bin/bash
# Release script for n8n CLI
# Creates GitHub/GitLab releases with cross-platform binaries

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

# Check if version is provided
if [ -z "$1" ]; then
    echo -e "${RED}❌ Version required. Usage: $0 <version>${NC}"
    echo -e "Example: $0 1.0.0"
    exit 1
fi

VERSION=$1
echo -e "${BLUE}Creating release v${VERSION}...${NC}"

# Validate version format (semantic versioning)
if ! echo "$VERSION" | grep -qE '^[0-9]+\.[0-9]+\.[0-9]+$'; then
    echo -e "${RED}❌ Invalid version format. Use semantic versioning (e.g., 1.0.0)${NC}"
    exit 1
fi

# Check if we're on main/master branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" != "main" ] && [ "$CURRENT_BRANCH" != "master" ]; then
    echo -e "${YELLOW}⚠️  Warning: Not on main/master branch (current: ${CURRENT_BRANCH})${NC}"
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Check for uncommitted changes
if ! git diff --quiet HEAD; then
    echo -e "${RED}❌ Uncommitted changes found. Commit or stash them first.${NC}"
    exit 1
fi

# Update VERSION file
echo "${VERSION}" > VERSION
echo -e "${YELLOW}Updated VERSION file to ${VERSION}${NC}"

# Update CHANGELOG.md
echo -e "${YELLOW}Updating CHANGELOG.md...${NC}"
TODAY=$(date '+%Y-%m-%d')
sed -i.bak "s/## \[Unreleased\]/## [Unreleased]\n\n## [${VERSION}] - ${TODAY}/" CHANGELOG.md
rm CHANGELOG.md.bak 2>/dev/null || true

# Build cross-platform binaries
echo -e "${YELLOW}Building cross-platform binaries...${NC}"
./scripts/build.sh all

# Create release packages
echo -e "${YELLOW}Creating release packages...${NC}"
cd dist

# Create tar.gz for Unix systems
tar -czf n8n-cli-v${VERSION}-linux-amd64.tar.gz n8n-cli-linux-amd64 ../README.md ../QUICK_START.md ../config.example.yaml
tar -czf n8n-cli-v${VERSION}-linux-arm64.tar.gz n8n-cli-linux-arm64 ../README.md ../QUICK_START.md ../config.example.yaml
tar -czf n8n-cli-v${VERSION}-darwin-amd64.tar.gz n8n-cli-darwin-amd64 ../README.md ../QUICK_START.md ../config.example.yaml
tar -czf n8n-cli-v${VERSION}-darwin-arm64.tar.gz n8n-cli-darwin-arm64 ../README.md ../QUICK_START.md ../config.example.yaml

# Create zip for Windows
zip n8n-cli-v${VERSION}-windows-amd64.zip n8n-cli-windows-amd64.exe ../README.md ../QUICK_START.md ../config.example.yaml

# Update checksums
sha256sum n8n-cli-v${VERSION}-*.tar.gz n8n-cli-v${VERSION}-*.zip > n8n-cli-v${VERSION}-checksums.txt

cd ..

# Create git tag
echo -e "${YELLOW}Creating git tag v${VERSION}...${NC}"
git add VERSION CHANGELOG.md
git commit -m "Release v${VERSION}"
git tag -a "v${VERSION}" -m "Release v${VERSION}"

# Push to remote
echo -e "${YELLOW}Pushing to remote...${NC}"
git push origin main
git push origin "v${VERSION}"

# Generate release notes
echo -e "${YELLOW}Generating release notes...${NC}"
cat > release-notes.md << EOF
# n8n CLI v${VERSION}

## What's New

$(sed -n "/## \[${VERSION}\]/,/## \[/p" CHANGELOG.md | head -n -1 | tail -n +2)

## Installation

### Download Pre-built Binaries

Choose the appropriate binary for your system:

- **Linux AMD64**: \`n8n-cli-v${VERSION}-linux-amd64.tar.gz\`
- **Linux ARM64**: \`n8n-cli-v${VERSION}-linux-arm64.tar.gz\`
- **macOS AMD64**: \`n8n-cli-v${VERSION}-darwin-amd64.tar.gz\`
- **macOS ARM64**: \`n8n-cli-v${VERSION}-darwin-arm64.tar.gz\`
- **Windows AMD64**: \`n8n-cli-v${VERSION}-windows-amd64.zip\`

### Quick Install

\`\`\`bash
# Download and extract (Linux example)
wget https://github.com/your-org/n8n-cli/releases/download/v${VERSION}/n8n-cli-v${VERSION}-linux-amd64.tar.gz
tar -xzf n8n-cli-v${VERSION}-linux-amd64.tar.gz
chmod +x n8n-cli-linux-amd64
sudo mv n8n-cli-linux-amd64 /usr/local/bin/n8n-cli
\`\`\`

### Verify Installation

\`\`\`bash
n8n-cli welcome
\`\`\`

## Checksums

Verify your download with SHA256 checksums:

\`\`\`
$(cat dist/n8n-cli-v${VERSION}-checksums.txt)
\`\`\`

---

**Full Changelog**: https://github.com/your-org/n8n-cli/compare/v$(git describe --tags --abbrev=0 HEAD^)...v${VERSION}
EOF

echo -e "${GREEN}✅ Release v${VERSION} created successfully!${NC}"
echo -e "${BLUE}Next steps:${NC}"
echo -e "1. Go to GitLab/GitHub releases page"
echo -e "2. Create new release with tag v${VERSION}"
echo -e "3. Upload files from dist/ directory:"
ls -la dist/n8n-cli-v${VERSION}-*
echo -e "4. Use release-notes.md content as release description"
echo -e ""
echo -e "${GREEN}🎉 Release completed!${NC}"