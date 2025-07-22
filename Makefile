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
	go build -o n8n-cli main.go
	chmod +x n8n-cli
	@echo "Build completed: n8n-cli"

# Run tests and validation
test:
	@echo "Running tests..."
	go test ./...
	go vet ./...
	@echo "Running CLI validation..."
	./n8n-cli --help > /dev/null
	./n8n-cli welcome > /dev/null
	@echo "Tests completed successfully"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f n8n-cli
	rm -f *.exe *.dll *.so *.dylib
	go clean
	@echo "Clean completed"

# Install CLI to system PATH
install: build
	@echo "Installing n8n CLI to /usr/local/bin..."
	sudo cp n8n-cli /usr/local/bin/
	sudo chmod +x /usr/local/bin/n8n-cli
	@echo "Installation completed"

# Set up development environment
setup:
	@echo "Setting up development environment..."
	./scripts/setup-dev.sh
	@echo "Development setup completed"

# Validate all workflow files
validate: build
	@echo "Validating workflow files..."
	./n8n-cli validate ./workflows/development/ || true
	./n8n-cli validate ./workflows/staging/ || true
	./n8n-cli validate ./workflows/production/ || true
	@echo "Validation completed"

# Deploy to development environment
deploy: build validate
	@echo "Deploying to development environment..."
	./scripts/pre-deploy.sh development
	./n8n-cli deploy --env development --dry-run
	./n8n-cli deploy --env development
	./scripts/post-deploy.sh development
	@echo "Deployment completed"

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build -o dist/n8n-cli-linux-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -o dist/n8n-cli-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o dist/n8n-cli-darwin-arm64 main.go
	GOOS=windows GOARCH=amd64 go build -o dist/n8n-cli-windows-amd64.exe main.go
	@echo "Multi-platform build completed in dist/"

# Create distribution package
dist: clean build-all
	@echo "Creating distribution packages..."
	mkdir -p dist
	cp README.md QUICK_START.md config.example.yaml dist/
	tar -czf dist/n8n-cli-linux-amd64.tar.gz -C dist n8n-cli-linux-amd64 README.md QUICK_START.md config.example.yaml
	tar -czf dist/n8n-cli-darwin-amd64.tar.gz -C dist n8n-cli-darwin-amd64 README.md QUICK_START.md config.example.yaml
	tar -czf dist/n8n-cli-darwin-arm64.tar.gz -C dist n8n-cli-darwin-arm64 README.md QUICK_START.md config.example.yaml
	zip -j dist/n8n-cli-windows-amd64.zip dist/n8n-cli-windows-amd64.exe dist/README.md dist/QUICK_START.md dist/config.example.yaml
	@echo "Distribution packages created in dist/"

# Development workflow
dev: build
	@echo "Starting development workflow..."
	./n8n-cli welcome
	@echo "CLI ready for development"

# GitLab CI simulation (local testing)
ci-test:
	@echo "Simulating GitLab CI pipeline..."
	go build -o n8n-cli main.go
	./n8n-cli validate ./workflows/ --verbose
	./n8n-cli --help
	./n8n-cli welcome
	@echo "CI simulation completed successfully"