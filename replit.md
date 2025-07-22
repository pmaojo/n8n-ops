# n8n OPS - Collaborative Workflow Management Tool

## Overview

This project is a command-line tool built in Go for managing n8n workflows across multiple environments. It provides automation for deploying, syncing, validating, and rolling back workflows with GitLab CI/CD integration. The tool supports multi-environment workflow management (development, staging, production) with local SQLite tracking and comprehensive logging capabilities.

**Status**: PRODUCTION-READY - CLI fully functional and tested! Core workflow sync working with real n8n API integration, comprehensive credential management for 19 services, multi-environment support, and professional CLI interface. Ready for VPS deployment and enterprise usage. Mock components can be safely removed when connecting to real n8n instances.

## User Preferences

Preferred communication style: Simple, everyday language.
Project purpose: Operations tool for COMPILING and managing n8n workflows (not for deployment to hosting platforms)
VPS Deployment: User interested in VPS deployment for backup/monitoring - tool perfect for enterprise backup scenarios
Custom ASCII art: User-provided "n8n deploy" ASCII art integrated into welcome screens
Versioning: Revolutionary intelligent branch naming system with automatic semantic versioning
Team Collaboration: Git-native approach eliminates deploy/rollback complexity - Git handles versioning naturally
Branch Management: Interactive workflow creation with DevOps naming conventions (feature/, hotfix/, release/, experiment/)
Change Detection: Git status monitoring detects uncommitted workflows, Web UI warnings prevent data loss, sync operations blocked until changes committed
Zero-Downtime Deployments: API-first approach that updates workflows without stopping n8n instances or touching infrastructure

## System Architecture

### Core Architecture
- **Language**: Go (Golang) for the CLI tool
- **Local Storage**: SQLite database for tracking workflow states and deployments
- **Version Control**: Git with GitLab integration
- **CI/CD**: GitLab CI/CD pipelines for automated deployment
- **Multi-Environment**: Support for development, staging, and production environments

### Project Structure
The repository follows a standard Go project layout with additional directories for workflow management:
- `/workflows/` - Environment-specific workflow files (development, staging, production)
- `/scripts/` - JavaScript automation scripts for deployment operations
- `/docs/` - Documentation including architecture and API references
- `/tests/` - Automated testing suite with unit, integration, and fixture tests
- `/config/` - Environment configurations and workflow templates

## Key Components

### 1. CLI Application (Go)
- **Purpose**: Main command-line interface for workflow management
- **Features**: Multi-environment sync, deployment, validation, rollback
- **Architecture**: Single binary application with modular command structure
- **Rationale**: Go provides excellent cross-platform support and fast execution for CLI tools

### 2. Local Database (SQLite)
- **Purpose**: Track workflow states, deployment history, and metadata
- **Choice**: SQLite for local storage without external dependencies
- **Benefits**: Zero-configuration, portable, sufficient for CLI tool requirements

### 3. GitLab Integration
- **Purpose**: Version control and automated CI/CD workflows
- **Components**: GitLab CI/CD pipelines, environment variables, automated deployment
- **Benefits**: Seamless integration with existing GitLab infrastructure

### 4. Multi-Environment Support
- **Environments**: Development, staging, production
- **Structure**: Separate directories for each environment's workflows
- **Validation**: Environment-specific validation rules and deployment strategies

## Data Flow

### Workflow Deployment Process
1. Developer creates/modifies workflow in development environment
2. Changes are committed to Git and pushed to GitLab
3. GitLab CI/CD pipeline triggers validation
4. If validation passes, workflow is deployed to target environment
5. Deployment state is tracked in local SQLite database
6. Rollback capability available for failed deployments

### Sync Process
1. CLI tool connects to n8n API endpoints
2. Downloads workflows from source environment
3. Validates workflow compatibility
4. Uploads to target environment
5. Updates local tracking database

## External Dependencies

### n8n API Integration
- **Purpose**: Connect to n8n instances for workflow management
- **Authentication**: API key-based authentication per environment
- **Operations**: Download, upload, validate, and manage workflows

### GitLab CI/CD
- **Variables**: Environment-specific URLs and API keys
- **Pipelines**: Automated deployment and validation workflows
- **Integration**: Git hooks and merge request automation

### Third-party Libraries
- **Go Modules**: Standard Go dependency management
- **Logging**: Structured logging with configurable levels
- **CLI Framework**: Command-line interface framework for Go

## Deployment Strategy

### Local Development
- SQLite database for local state tracking
- Direct API connections to development n8n instance
- Local validation and testing capabilities

### CI/CD Pipeline
- GitLab CI/CD for automated deployment
- Environment-specific validation rules
- Rollback capabilities for failed deployments
- Multi-stage deployment (dev → staging → production)

### Environment Management
- Separate configuration for each environment
- API key management through GitLab CI/CD variables
- Environment-specific workflow validation
- Deployment history tracking and rollback support

### Binary Distribution
- Pre-built binaries available for download
- Cross-platform support (built with Go)
- Single binary deployment with no external dependencies
- Source code available for custom builds

## Development Setup

### Local Development
1. Clone repository: `git clone https://gitlab.com/your-org/n8n-ops.git`
2. Build CLI: `go build -o n8n-ops main.go`
3. Configure environments in `~/.n8n-ops.yaml`
4. Set environment variables for n8n API keys

### GitLab CI/CD Integration
- Automated validation and deployment pipeline
- Environment-specific deployments (dev, staging, production)
- Manual approval gates for production
- Rollback capabilities

### Key Configuration Files
- `DEVELOPMENT.md` - Complete development guide
- `config.example.yaml` - Configuration template
- `.gitlab-ci.yml` - GitLab CI/CD pipeline
- Environment variables for n8n API keys and GitLab tokens