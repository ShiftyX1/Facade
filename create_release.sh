#!/bin/bash

# Script to create a local release of the Facade project
# Usage: ./create_release.sh [version]

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

if [ -n "$1" ]; then
    VERSION="$1"
else
    VERSION=$(git describe --tags --always --dirty)
    warn "No version specified, using: $VERSION"
fi

info "Creating release for version: $VERSION"

if [[ "$VERSION" == *"dirty"* ]]; then
    warn "Working directory has uncommitted changes"
fi

RELEASE_DIR="release-$VERSION"
info "Creating release directory: $RELEASE_DIR"
rm -rf "$RELEASE_DIR"
mkdir -p "$RELEASE_DIR"

PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

build_platform() {
    local platform=$1
    local goos=$(echo $platform | cut -d'/' -f1)
    local goarch=$(echo $platform | cut -d'/' -f2)
    
    info "Building for $goos/$goarch..."
    
    local suffix="$goos-$goarch"
    local ext=""
    if [ "$goos" = "windows" ]; then
        ext=".exe"
        suffix="$suffix$ext"
    fi
    
    GOOS=$goos GOARCH=$goarch CGO_ENABLED=0 go build \
        -ldflags="-s -w -X 'main.Version=$VERSION' -X 'main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)'" \
        -o "$RELEASE_DIR/facade-$suffix" \
        cmd/facade/main.go
    
    cd "$RELEASE_DIR"
    if [ "$goos" = "windows" ]; then
        zip "facade-$VERSION-$suffix.zip" "facade-$suffix"
        rm "facade-$suffix"
    else
        tar czf "facade-$VERSION-$suffix.tar.gz" "facade-$suffix"
        rm "facade-$suffix"
    fi
    cd ..
    
    info "✓ Created facade-$VERSION-$suffix archive"
}

info "Running tests..."
make test || error "Tests failed"

info "Building for all platforms..."
for platform in "${PLATFORMS[@]}"; do
    build_platform "$platform"
done

info "Copying configuration files..."
cp -r configs "$RELEASE_DIR/"
cp README.md "$RELEASE_DIR/"
cp LICENSE "$RELEASE_DIR/" 2>/dev/null || warn "LICENSE file not found"

cd "$RELEASE_DIR"
echo "# Facade $VERSION Release" > RELEASE_NOTES.md
echo "" >> RELEASE_NOTES.md
echo "## Files in this release:" >> RELEASE_NOTES.md
echo "" >> RELEASE_NOTES.md
ls -la *.tar.gz *.zip 2>/dev/null | awk '{print "- " $9 " (" $5 " bytes)"}' >> RELEASE_NOTES.md
echo "" >> RELEASE_NOTES.md
echo "## Installation:" >> RELEASE_NOTES.md
echo "" >> RELEASE_NOTES.md
echo "1. Download the appropriate archive for your platform" >> RELEASE_NOTES.md
echo "2. Extract the archive" >> RELEASE_NOTES.md
echo "3. Run \`./facade-* --help\` to get started" >> RELEASE_NOTES.md
echo "" >> RELEASE_NOTES.md
echo "## Quick Start:" >> RELEASE_NOTES.md
echo "" >> RELEASE_NOTES.md
echo "\`\`\`bash" >> RELEASE_NOTES.md
echo "# Extract (example for Linux)" >> RELEASE_NOTES.md
echo "tar xzf facade-$VERSION-linux-amd64.tar.gz" >> RELEASE_NOTES.md
echo "" >> RELEASE_NOTES.md
echo "# Run with example config" >> RELEASE_NOTES.md
echo "./facade-linux-amd64 -c configs/example.yaml" >> RELEASE_NOTES.md
echo "\`\`\`" >> RELEASE_NOTES.md

cd ..

info "🎉 Release created successfully!"
echo ""
echo "📦 Release directory: $RELEASE_DIR"
echo "📋 Files created:"
ls -la "$RELEASE_DIR"/*.tar.gz "$RELEASE_DIR"/*.zip 2>/dev/null
echo ""
echo "🚀 To test a build:"
echo "  cd $RELEASE_DIR"
echo "  tar xzf facade-$VERSION-linux-amd64.tar.gz  # or appropriate archive"
echo "  ./facade-linux-amd64 --version"
echo ""
echo "📖 Release notes: $RELEASE_DIR/RELEASE_NOTES.md"
