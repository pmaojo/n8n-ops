# n8n OPS - Collaborative Workflow Management Tool

A Go-based operations tool for **compiling, managing, and deploying** n8n workflows across multiple environments with GitLab CI/CD integration. Features custom ASCII art, multilingual support (English/Spanish), and comprehensive workflow lifecycle management.

## Features

- 🔄 **Multi-Environment Sync** - Sync workflows between development, staging, and production
- 🚀 **Automated Deployment** - Deploy workflows with validation and rollback support
- ✅ **Workflow Validation** - Validate workflow JSON files for structure and compatibility
- 🔒 **GitLab CI/CD Integration** - Complete CI/CD pipeline with approval gates
- 📊 **Local Tracking** - SQLite database for tracking workflow states and deployments
- 🔄 **Rollback Support** - Rollback to previous deployments with one command
- 📝 **Comprehensive Logging** - Structured logging with configurable levels
- 🏗️ **Project Scaffolding** - Initialize new workflow projects with standard structure
- 🎨 **Custom ASCII Art** - Beautiful futuristic interface with Matrix effects
- 🌍 **Multilingual Support** - English and Spanish language support
- 📡 **Complete n8n API** - Full CRUD operations for workflows and executions

## Installation

### Download Pre-built Binaries

Download the latest release from the [releases page](https://github.com/n8n-workflows/cli/releases).

### Build from Source

```bash
# Clone the repository
git clone https://gitlab.com/your-org/n8n-ops.git
cd n8n-ops

# Build the operations tool
go build -o n8n-ops main.go

# Make executable (Linux/macOS)  
chmod +x n8n-ops

# Run setup script for development
./scripts/setup-dev.sh
