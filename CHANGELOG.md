# Changelog

All notable changes to the n8n CLI project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of n8n CLI tool
- Custom ASCII art with user-provided "n8n deploy" design
- Multilingual support (English and Spanish)
- Complete n8n API integration with CRUD operations
- GitLab CI/CD pipeline with multi-environment support
- Comprehensive development and deployment documentation
- Automated validation and testing scripts
- Workflow templates and examples
- SQLite database for local tracking
- Branch-based environment detection

### Features
- `sync` command - Sync workflows from n8n instances to local files
- `deploy` command - Deploy workflows from local files to n8n instances
- `validate` command - Validate workflow JSON files
- `welcome` command - Display futuristic welcome screen with Matrix effects
- `branch` command - Manage Git branch to environment mappings
- `init` command - Initialize new n8n workflow projects
- `rollback` command - Rollback to previous deployments

### Security
- API key protection with environment variables
- No secrets committed to repository
- Secure GitLab CI/CD variable handling
- Input validation for all user inputs

### Documentation
- Complete development setup guide
- Quick start guide for 5-minute setup
- Contributing guidelines for team collaboration
- Repository setup guide for GitLab configuration
- Issue and merge request templates

## [1.0.0] - 2025-01-22

### Added
- Initial stable release
- Core CLI functionality
- Multi-environment support (development, staging, production)
- GitLab CI/CD integration
- Workflow management capabilities
- Custom ASCII art integration
- Multilingual interface support