# Security Policy

## Supported Versions

We actively support the following versions of n8n CLI:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

The n8n CLI team takes security vulnerabilities seriously. We appreciate your efforts to responsibly disclose your findings.

### How to Report

**DO NOT** create a public GitHub/GitLab issue for security vulnerabilities.

Instead, please report security vulnerabilities by:

1. **Email**: Send details to `security@yourcompany.com`
2. **Private Issue**: Create a confidential issue in GitLab (if available)
3. **GPG Encrypted**: For sensitive issues, use our public GPG key

### What to Include

Please include the following information in your report:

- **Description**: Clear description of the vulnerability
- **Impact**: What an attacker could achieve
- **Steps to Reproduce**: Detailed steps to reproduce the issue
- **Proof of Concept**: Code or commands demonstrating the vulnerability
- **Environment**: OS, Go version, n8n CLI version
- **Suggested Fix**: If you have ideas for a fix

### Example Report

```
Subject: [SECURITY] API Key Exposure in n8n CLI

Description:
The n8n CLI stores API keys in plaintext in configuration files, 
making them accessible to other processes on the system.

Impact:
- API keys could be read by malicious software
- Credentials could be accidentally committed to version control
- Unauthorized access to n8n instances

Steps to Reproduce:
1. Configure n8n CLI with API key
2. Check file permissions on ~/.n8n-ops.yaml
3. API key is readable by all processes

Environment:
- OS: Ubuntu 22.04
- Go: 1.19
- n8n CLI: 1.0.0

Suggested Fix:
- Store API keys in OS keyring
- Encrypt configuration file
- Restrict file permissions to 600
```

### Response Timeline

We aim to respond to security reports according to the following timeline:

- **Initial Response**: Within 48 hours
- **Assessment**: Within 5 business days  
- **Fix Development**: Within 30 days for critical issues
- **Public Disclosure**: After fix is released

### Security Measures

The n8n CLI implements the following security measures:

#### API Key Protection
- API keys stored in environment variables
- Configuration files exclude sensitive data by default
- Warning messages when API keys detected in config files

#### Network Security
- HTTPS-only connections to n8n instances
- Certificate validation enabled by default
- Request timeout limits to prevent DoS

#### Input Validation
- Workflow JSON validation before deployment
- Path traversal protection for file operations
- Command injection prevention

#### Access Control
- No privileged operations required
- Principle of least privilege
- Clear separation between environments

### Vulnerability Categories

We consider the following types of vulnerabilities:

#### Critical Severity
- Remote code execution
- Unauthorized access to n8n instances
- API key exposure or theft
- Data corruption or loss

#### High Severity
- Privilege escalation
- Authentication bypass
- Information disclosure
- Denial of service attacks

#### Medium Severity
- Input validation issues
- Configuration vulnerabilities
- Logging sensitive data

#### Low Severity
- Error message information disclosure
- Minor configuration issues

### Security Best Practices

When using n8n CLI, follow these security practices:

#### API Key Management
- Use environment variables for API keys
- Rotate API keys regularly
- Use minimum required permissions
- Monitor API key usage

#### Configuration Security
- Restrict configuration file permissions: `chmod 600 ~/.n8n-ops.yaml`
- Don't commit configuration files to version control
- Use separate API keys for each environment

#### Network Security
- Use VPN when accessing n8n instances remotely
- Enable SSL/TLS certificate validation
- Monitor network traffic for anomalies

#### Environment Isolation
- Use separate API keys for dev/staging/production
- Implement proper access controls in GitLab CI/CD
- Audit deployment logs regularly

### Acknowledgments

We appreciate the security research community and will acknowledge researchers who report vulnerabilities:

- Public acknowledgment on our security page (optional)
- Credit in release notes for fixed vulnerabilities
- Hall of Fame for significant contributions

### Contact Information

- **Security Email**: security@yourcompany.com
- **General Contact**: support@yourcompany.com
- **Website**: https://yourcompany.com/security

### GPG Key

```
-----BEGIN PGP PUBLIC KEY BLOCK-----
[Your GPG public key here]
-----END PGP PUBLIC KEY BLOCK-----
```

---

Thank you for helping keep n8n CLI and our users safe!