# Contributing to n8n CLI

We welcome contributions! This guide explains how to set up your development environment and contribute to the project.

## Development Setup

### Prerequisites
- Go 1.19+
- Git
- Access to n8n instances

### Quick Setup
```bash
# Clone and setup
git clone https://gitlab.com/your-org/n8n-ops.git
cd n8n-ops
./scripts/setup-dev.sh

# Build and test
go build -o n8n-ops main.go
./n8n-ops welcome
```

## Branch Strategy

### Main Branches
- `main` → Production environment
- `staging` → Staging environment  
- `develop` → Development environment

### Feature Development
```bash
# Create feature branch from develop
git checkout develop
git pull origin develop
git checkout -b feature/your-feature-name

# Make changes, test, commit
git add .
git commit -m "feat: add new feature"
git push origin feature/your-feature-name

# Create merge request to develop branch
```

## Code Standards

### Go Code Style
- Follow Go conventions (`gofmt`, `go vet`)
- Add comments for exported functions
- Write tests for new functionality
- Use meaningful variable names

### Commit Messages
```
feat: add new payment workflow template
fix: resolve API timeout in production sync  
docs: update deployment guide
style: format code with gofmt
test: add validation tests for workflows
refactor: simplify n8n client interface
```

## Testing

### Validate Changes
```bash
# Validate workflow files
./n8n-ops validate ./workflows/

# Test CLI compilation
go build -o n8n-ops main.go

# Test basic functionality  
./n8n-ops --help
./n8n-ops welcome --lang es

# Run Go tests
go test ./...
```

## Workflow Development

### Adding New Workflows
1. Create workflow in development n8n instance
2. Sync to local: `./n8n-ops sync --env development`
3. Validate: `./n8n-ops validate ./workflows/development/`
4. Test deployment: `./n8n-ops deploy --env development --dry-run`
5. Commit changes to git

### Workflow Templates
- Add templates to `config/templates/`
- Follow JSON structure in existing examples
- Include placeholder variables like `{{WEBHOOK_PATH}}`

## Merge Request Process

### Before Creating MR
- [ ] Code builds successfully
- [ ] All workflow files validate  
- [ ] Tests pass
- [ ] Documentation updated if needed
- [ ] Commit messages follow convention

### MR Requirements
- **Target Branch**: Usually `develop` for features
- **Description**: Explain what changes and why
- **Testing**: How you tested the changes
- **Screenshots**: For UI/visual changes

### Review Process
1. Automated CI/CD checks must pass
2. At least one team member review required
3. All conversations resolved
4. Squash merge preferred

## Release Process

### Development → Staging
1. Create MR from `develop` to `staging`
2. Review and approval required
3. Manual deployment to staging via GitLab CI/CD
4. QA testing in staging environment

### Staging → Production  
1. Create MR from `staging` to `main`
2. Senior developer approval required
3. Manual deployment to production
4. Monitor deployment success

## Adding Features

### New CLI Commands
1. Create command file in `cmd/` directory
2. Follow existing patterns (cobra commands)
3. Add help text and examples
4. Update main command registration
5. Add to documentation

### New n8n API Methods
1. Add method to `internal/client/n8n.go`
2. Include error handling and logging
3. Add configuration options if needed
4. Update client interface

### Internationalization
1. Add translations to `internal/i18n/messages.go`
2. Use `i18n.T("key")` in code
3. Test with `--lang es` flag
4. Update language support documentation

## Environment Configuration

### Local Development
```yaml
# ~/.n8n-ops.yaml
environments:
  development:
    url: "https://n8n-dev.yourcompany.com"
    api_key: "${N8N_API_KEY_DEV}"
```

### GitLab CI/CD Variables
| Variable | Description | Masked |
|----------|-------------|---------|
| `N8N_API_KEY_DEV` | Development n8n API key | ✅ |
| `N8N_API_KEY_STAGING` | Staging n8n API key | ✅ |
| `N8N_API_KEY_PROD` | Production n8n API key | ✅ |

## Troubleshooting

### Common Issues
- **Build fails**: Run `go mod tidy` and `go clean`
- **API connection fails**: Check API keys and URLs
- **Workflow validation fails**: Check JSON syntax
- **GitLab CI fails**: Check environment variables

### Getting Help
- Check existing issues in GitLab
- Review documentation files
- Ask team members
- Create detailed issue with:
  - Steps to reproduce
  - Expected vs actual behavior
  - Environment details
  - Error logs

## Security

### API Keys
- Never commit API keys to git
- Use environment variables only
- Rotate keys regularly
- Limit key permissions to minimum required

### Code Security
- No hardcoded secrets in source code
- Validate all user inputs
- Use HTTPS for all API calls
- Log security events appropriately

---

Thank you for contributing to n8n CLI! 🚀