# Schemactor Makefile
# Provides build, release, and installation targets

.PHONY: help build build-all release install clean version

# Default target
help:
	@echo "Schemactor Build System"
	@echo ""
	@echo "Targets:"
	@echo "  build       Build for current platform (development version)"
	@echo "  build-all   Build for all supported platforms (development version)"
	@echo "  release     Build release binaries for all platforms"
	@echo "  install     Install current build to GOPATH/bin"
	@echo "  clean       Remove build artifacts"
	@echo "  version     Show version information"
	@echo "  help        Show this help"
	@echo ""
	@echo "Examples:"
	@echo "  make build          # Build for current platform"
	@echo "  make release        # Build release (requires git tag)"
	@echo "  make install        # Install to GOPATH/bin"

# Build for current platform (development)
build:
	@./scripts/build.sh --platforms "$$(go env GOOS)/$$(go env GOARCH)"

# Build for all supported platforms (development)
build-all:
	@./scripts/build.sh --platforms "linux/amd64,darwin/amd64"

# Build release binaries for all platforms
release:
	@echo "Building release binaries..."
	@./scripts/build.sh --platforms "linux/amd64,darwin/amd64" --release
	@echo ""
	@echo "Release binaries created in dist/"
	@echo "Upload these files to your release:"

# Install to GOPATH/bin
install:
	@echo "Installing schemactor to GOPATH/bin..."
	@./scripts/build.sh --platforms "$$(go env GOOS)/$$(go env GOARCH)" --output-file "schemactor"
	@cp dist/schemactor "$$(go env GOPATH)/bin/schemactor"
	@echo "Installed to $$(go env GOPATH)/bin/schemactor"
	@echo "Run 'schemactor --version' to verify installation"

# Clean build artifacts
clean:
	@echo "Removing build artifacts..."
	@rm -rf dist/
	@rm -f schemactor schemator-*
	@echo "Clean complete"

# Show version information
version:
	@if [ -f dist/schemactor ]; then \
		echo "Built version:"; \
		./dist/schemactor --version; \
	else \
		echo "Development version:"; \
		go run ./cmd/schemactor --version; \
	fi

# Additional convenience targets

# Quick development build (alias for build)
dev: build

# Local installation without building (assumes dist/schemactor exists)
install-local:
	@if [ ! -f dist/schemactor ]; then \
		echo "Error: No built binary found. Run 'make build' first."; \
		exit 1; \
	fi
	@cp dist/schemactor "$$(go env GOPATH)/bin/schemactor"
	@echo "Installed to $$(go env GOPATH)/bin/schemactor"