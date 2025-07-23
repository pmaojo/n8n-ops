# n8n CLI Makefile
# Simplifies common development tasks

.PHONY: build test clean install setup validate deploy help

# Default target
help:
	@echo "n8n CLI Development Commands:"
	@echo ""
	@echo "  build      Build the CLI binary"
	@echo "  test       Run tests and validation"
	@echo "  clean      Clean build artifacts"
	@echo "  install    Install to local system"
	@echo "  setup      Set up development environment"
	@echo "  validate   Validate all workflows"
	@echo "  deploy     Deploy to development environment"
	@echo "  help       Show this help message"
	@echo ""

# Build the CLI binary
build:
	@echo "Building n8n CLI..."
	go build -o n8n-ops main.go
	chmod +x n8n-ops
	@echo "Build completed: n8n-ops"

# Run tests and validation
test:
	@echo "Running tests..."
	go test ./...
	go vet ./...
	@echo "Running CLI validation..."
	./n8n-ops --help > /dev/null
	./n8n-ops welcome > /dev/null
	@echo "Tests completed successfully"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f n8n-ops
	rm -f *.exe *.dll *.so *.dylib
	go clean
	@echo "Clean completed"

# Install CLI to system PATH
install: build
	@echo "Installing n8n CLI to /usr/local/bin..."
	sudo cp n8n-ops /usr/local/bin/
	sudo chmod +x /usr/local/bin/n8n-ops
	@echo "Installation completed"

# Set up development environment
setup:
	@echo "Setting up development environment..."
	./scripts/setup-dev.sh
	@echo "Development setup completed"

# Validate all workflow files
validate: build
	@echo "Validating workflow files..."
	./n8n-ops validate ./workflows/development/ || true
	./n8n-ops validate ./workflows/staging/ || true
	./n8n-ops validate ./workflows/production/ || true
	@echo "Validation completed"

# Deploy to development environment
deploy: build validate
	@echo "Deploying to development environment..."
	./scripts/pre-deploy.sh development
	./n8n-ops deploy --env development --dry-run
	./n8n-ops deploy --env development
	./scripts/post-deploy.sh development
	@echo "Deployment completed"

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build -o dist/n8n-ops-linux-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -o dist/n8n-ops-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o dist/n8n-ops-darwin-arm64 main.go
	GOOS=windows GOARCH=amd64 go build -o dist/n8n-ops-windows-amd64.exe main.go
	@echo "Multi-platform build completed in dist/"

# Create distribution package
dist: clean build-all
	@echo "Creating distribution packages..."
	mkdir -p dist
	cp README.md QUICK_START.md config.example.yaml dist/
	tar -czf dist/n8n-ops-linux-amd64.tar.gz -C dist n8n-ops-linux-amd64 README.md QUICK_START.md config.example.yaml
	tar -czf dist/n8n-ops-darwin-amd64.tar.gz -C dist n8n-ops-darwin-amd64 README.md QUICK_START.md config.example.yaml
	tar -czf dist/n8n-ops-darwin-arm64.tar.gz -C dist n8n-ops-darwin-arm64 README.md QUICK_START.md config.example.yaml
	zip -j dist/n8n-ops-windows-amd64.zip dist/n8n-ops-windows-amd64.exe dist/README.md dist/QUICK_START.md dist/config.example.yaml
	@echo "Distribution packages created in dist/"

# Development workflow
dev: build
	@echo "Starting development workflow..."
	./n8n-ops welcome
	@echo "CLI ready for development"

# GitLab CI simulation (local testing)
ci-test:
	@echo "Simulating GitLab CI pipeline..."
	go build -o n8n-ops main.go
	./n8n-ops validate ./workflows/ --verbose
	./n8n-ops --help
	./n8n-ops welcome
	@echo "CI simulation completed successfully"