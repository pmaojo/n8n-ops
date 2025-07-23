# n8n-ops CLI Tool - Replit Documentation

## Overview

n8n-ops is a comprehensive command-line tool built in Go for managing n8n workflows across multiple environments. It provides enterprise-grade automation, version control, and multi-environment orchestration for n8n workflows with GitLab CI/CD integration.

## User Preferences

Preferred communication style: Simple, everyday language.

## System Architecture

### Core Technology Stack
- **Language**: Go 1.19+ with Cobra CLI framework
- **Configuration**: YAML-based configuration with environment variables
- **Database**: SQLite for local tracking and metadata
- **Version Control**: Git with GitLab integration
- **API Integration**: n8n REST API for workflow management

### CLI Architecture
The application follows a modular CLI structure using the Cobra framework:
- Command hierarchy with subcommands (`sync`, `deploy`, `validate`, `check`, `status`, etc.)
- Environment-based operations (development, staging, production)
- Configuration management through YAML files and environment variables
- Structured logging with different levels (info, debug, error)

## Key Components

### 1. Command Structure
- **Core Commands**: `sync`, `deploy`, `validate`, `check`, `status`, `branch`, `init`, `rollback`
- **Utility Commands**: `welcome`, `help`, `version`
- **Environment Support**: All commands support `--env` flag for multi-environment operations

### 2. Workflow Management
- **Sync Engine**: Bidirectional synchronization between n8n instances and local JSON files
- **Validation System**: JSON schema validation for workflow files
- **Template System**: Pre-built workflow templates (basic-webhook, data-processing, scheduled-task)
- **Metadata Tracking**: Git commit tracking and sync timestamps

### 3. Multi-Environment Support
- **Environment Isolation**: Separate directories for development, staging, and production workflows
- **Branch Mapping**: Git branch to environment mapping system
- **Configuration Separation**: Environment-specific API keys and URLs
- **Deployment Gates**: Manual approval gates for staging and production deployments

### 4. Security Layer
- **API Key Management**: Environment variable-based credential storage
- **No Secrets in Code**: All sensitive data externalized to environment variables
- **Input Validation**: Comprehensive validation for all user inputs
- **Secure CI/CD**: GitLab CI/CD variables for secure deployment

## Data Flow

### 1. Workflow Synchronization Flow
```
n8n Instance → n8n REST API → n8n-ops CLI → Local JSON Files → Git Repository
```

### 2. Deployment Flow
```
Git Repository → GitLab CI/CD → n8n-ops CLI → n8n REST API → n8n Instance
```

### 3. File Structure
```
workflows/
├── development/     # Development environment workflows
├── staging/         # Staging environment workflows
├── production/      # Production environment workflows
└── README.md        # Documentation
```

### 4. Metadata Management
- **Sync Metadata**: Tracks sync timestamps, Git commits, and environment information
- **Local Database**: SQLite database for deployment tracking and rollback capabilities
- **Git Integration**: Automatic detection of uncommitted workflow changes

## External Dependencies

### 1. n8n API Integration
- **REST API**: Full CRUD operations for workflows via n8n REST API
- **Authentication**: API key-based authentication for each environment
- **Real-time Updates**: Hot deployment without service interruption

### 2. GitLab Integration
- **CI/CD Pipeline**: 6-stage pipeline (build, validate, test, deploy)
- **Branch Protection**: Protected branches with manual approval gates
- **Issue Templates**: Bug report and feature request templates
- **Merge Request Templates**: Standardized PR process

### 3. Version Control
- **Git Flow**: Standard Git flow with feature branches
- **Branch Strategy**: Main (production), staging, develop branches
- **Semantic Versioning**: SemVer compliance for releases

### 4. Configuration Management
- **Environment Variables**: API keys, URLs, and credentials
- **YAML Configuration**: User preferences and default settings
- **Template System**: Workflow templates with variable substitution

## Deployment Strategy

### 1. Zero-Downtime Architecture
- **Hot Updates**: API-based workflow updates without service interruption
- **Non-invasive**: Only manages workflows, doesn't touch n8n infrastructure
- **Instant Rollback**: Quick reversion to previous workflow versions

### 2. Multi-Environment Pipeline
- **Development**: Automatic deployment on `develop` branch
- **Staging**: Manual deployment with approval on `staging` branch
- **Production**: Manual deployment with approval on `main` branch

### 3. CI/CD Pipeline Stages
1. **Build**: Compile Go binary and validate syntax
2. **Validate**: JSON schema validation and workflow checks
3. **Test**: Run automated tests and CLI functionality checks
4. **Deploy**: Environment-specific deployment with API calls
5. **Backup**: Automatic backup before production deployments
6. **Monitor**: Post-deployment validation and health checks

### 4. Security Considerations
- **API Key Rotation**: Support for rotating API keys across environments
- **Audit Trail**: Complete deployment history and change tracking
- **Access Control**: Environment-based access restrictions
- **Credential Isolation**: Separate credentials for each environment

### 5. Monitoring and Alerting
- **Status Monitoring**: Real-time workflow status across environments
- **Change Detection**: Automatic detection of uncommitted workflow changes
- **Web UI**: Dashboard for monitoring and management (planned feature)
- **Notification System**: Integration with Slack/Discord for alerts

The system is designed to be enterprise-ready with comprehensive error handling, logging, and operational tooling for managing n8n workflows at scale.

## Recent Changes (Latest)

### Project Merge Success - Production Ready (July 23, 2025)
- **Completed**: Successful merge with passing tests - system confirmed fully operational
- **Visual Interface**: Spectacular ASCII art and futuristic Matrix-style interface restored
- **Real-time Monitoring**: Active workflow failure detection with GitLab issue creation
- **Multi-mode Operations**: Daemon mode, monitoring mode, and CLI commands all functional
- **Mock Server Integration**: Complete mock n8n API server for development and testing
- **Enterprise Grade**: 47 Go files, comprehensive CLI tool with professional-grade architecture
- **User Confirmed**: System described as "impresionante" (impressive) with all tests passing

### n8n Client Refactor - SOLID Principles Implementation (July 23, 2025)
- **Completed**: Full refactor of n8n API client following SOLID principles
- **Interface Segregation**: Split into WorkflowReader, WorkflowWriter, WorkflowExecutor interfaces
- **Dependency Injection**: HTTP client injected from outside for better testability
- **Context Support**: All methods now accept context.Context for proper cancellation
- **Generic Request Handler**: doRequest[T]() function eliminates code duplication
- **Constructor Fix**: No side effects in constructor, separate Ping() method for connectivity testing
- **Performance**: Optimized HTTP transport with proper timeouts and connection pooling
- **Testing**: Comprehensive test suite with 23.6% code coverage, including benchmarks
- **Error Handling**: Structured error types with helper functions for error checking
- **Code Quality**: Following Go best practices with proper documentation and type safety

### Automated Issue Management System (July 23, 2025)
- **Completed**: Full workflow failure monitoring and GitLab issue management system
- **Automatic Issue Creation**: Creates detailed GitLab issues when workflows fail consecutively
- **Smart Detection**: Configurable failure thresholds and environment-specific settings
- **Recovery Tracking**: Automatically updates issues when workflows recover
- **Rich Context**: Issues include failure details, troubleshooting steps, and metadata
- **Flexible Configuration**: Command-line options and environment-specific overrides
- **Demo Mode**: Mock issue manager for testing without GitLab integration
- **Comprehensive Documentation**: Monitoring guide with setup and best practices
- **GitLab Backend Integration**: Complete GitLab ecosystem integration with CI/CD, issues, and collaboration features