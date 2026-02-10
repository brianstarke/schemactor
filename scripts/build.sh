#!/bin/bash

# Schemactor Build Script
# Handles version detection and builds for different platforms

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
PLATFORMS="linux/amd64,darwin/amd64"
OUTPUT_DIR="./dist"
VERSION="dev"

# Parse arguments
PLATFORMS="linux/amd64,darwin/amd64"
VERSION="dev"
RELEASE=false
OUTPUT_FILE="schemactor"

while [[ $# -gt 0 ]]; do
    case $1 in
        --platforms)
            PLATFORMS="$2"
            shift 2
            ;;
        --output-dir)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        --version)
            VERSION="$2"
            shift 2
            ;;
        --release)
            RELEASE=true
            shift
            ;;
        --output-file)
            OUTPUT_FILE="$2"
            shift 2
            ;;
        --help)
            echo "Usage: $0 [options]"
            echo "  --platforms <list>  Comma-separated list of platforms (default: linux/amd64,darwin/amd64)"
            echo "  --output-dir <dir>   Output directory (default: ./dist)"
            echo "  --version <ver>      Override version (auto-detected if not provided)"
            echo "  --release            Build for release (requires clean git state)"
            echo "  --output-file <name> Base output filename (default: schemactor)"
            echo "  --help               Show this help"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Detect version if not provided
if [[ "$VERSION" == "dev" ]]; then
    detect_version() {
        local version
        
        # Check if we're in a git repo
        if git rev-parse --git-dir > /dev/null 2>&1; then
            # Check if we're on an exact tag
            if exact_tag=$(git describe --tags --exact-match HEAD 2>/dev/null); then
                version="$exact_tag"
            elif git tag > /dev/null 2>&1; then
                # Has tags - use describe to get latest tag + commits
                version=$(git describe --tags --always --abbrev=7 2>/dev/null)
                # If the result doesn't start with a proper version number, prefix it
                if [[ ! "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+.*$ ]]; then
                    local commit=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
                    version="v0.1.0-dev-$commit"
                fi
            else
                # No tags - use development version
                local commit=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
                version="v0.1.0-dev-$commit"
            fi
            
            # Check for uncommitted changes
            if [[ -n $(git status --porcelain 2>/dev/null) ]]; then
                version="${version}-dirty"
            fi
        else
            # Not a git repo
            version="v0.1.0-dev-unknown"
        fi
        
        echo "$version"
    }
    
    VERSION=$(detect_version)
fi

# Check release requirements
if [[ "$RELEASE" == true ]]; then
    if [[ "$VERSION" == *"dirty"* ]]; then
        echo -e "${RED}Error: Cannot create release build with dirty working tree${NC}"
        echo "Please commit or stash your changes first."
        exit 1
    fi
    
    if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo -e "${RED}Error: Release version must be a clean semantic version (e.g., v0.1.0)${NC}"
        echo "Current detected version: $VERSION"
        echo "Create a git tag for proper release versioning."
        exit 1
    fi
fi

echo -e "${BLUE}Building schemactor${NC}"
echo -e "${BLUE}Version: ${GREEN}$VERSION${NC}"
echo -e "${BLUE}Platforms: ${GREEN}$PLATFORMS${NC}"
echo -e "${BLUE}Release build: ${GREEN}$RELEASE${NC}"
echo

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Build function
build_platform() {
    local platform="$1"
    local goos="${platform%/*}"
    local goarch="${platform#*/}"
    local output_name="$OUTPUT_FILE"
    
    # Add platform suffix for cross-platform builds
    if [[ "$PLATFORMS" == *","* ]]; then
        output_name="${OUTPUT_FILE}-${goos}-${goarch}"
    fi
    
    # Add version prefix for release builds
    if [[ "$RELEASE" == true ]]; then
        output_name="${OUTPUT_FILE}-${VERSION}-${goos}-${goarch}"
    fi
    
    echo -e "${YELLOW}Building for ${goos}/${goarch} -> ${output_name}${NC}"
    
    GOOS="$goos" GOARCH="$goarch" go build \
        -ldflags="-X main.Version=${VERSION}" \
        -o "${OUTPUT_DIR}/${output_name}" \
        ./cmd/schemactor
    
    # Verify the binary was created
    if [[ ! -f "${OUTPUT_DIR}/${output_name}" ]]; then
        echo -e "${RED}Error: Failed to build ${output_name}${NC}"
        return 1
    fi
    
    echo -e "${GREEN}✓ Built ${output_name}${NC}"
}

# Split platforms and build for each
IFS=',' read -ra PLATFORM_ARRAY <<< "$PLATFORMS"
for platform in "${PLATFORM_ARRAY[@]}"; do
    build_platform "$platform"
done

echo
echo -e "${GREEN}Build complete!${NC}"
echo -e "${BLUE}Output directory: ${GREEN}$OUTPUT_DIR${NC}"
echo -e "${BLUE}Version: ${GREEN}$VERSION${NC}"

# Show built files
echo
echo -e "${BLUE}Built files:${NC}"
ls -la "$OUTPUT_DIR"