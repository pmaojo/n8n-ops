# 🚀 n8n-ops Quick Start Guide

**Get up and running in 5 minutes!**

## The Easiest Way: Interactive Onboarding

```bash
# Run the interactive onboarding wizard
./n8n-ops onboard
```

This wizard will:
- Create your project structure
- Configure your environments
- Set up API keys
- Test connections
- Show you next steps

## Manual Setup (Alternative)

### Step 1: Install n8n-ops

```bash
# Clone the repository
git clone https://gitlab.com/your-org/n8n-ops.git
cd n8n-ops

# Build the CLI
go build -o n8n-ops main.go

# Make it executable (Linux/macOS)
chmod +x n8n-ops

# Test it works
./n8n-ops welcome
```

### Step 2: Initialize Your Project

```bash
# Create a new project in the current directory
./n8n-ops init .

# Or create a new project in a specific directory
./n8n-ops init my-workflows
```

### Step 3: Get n8n API Keys

1. **Log in to your n8n instance** (e.g., http://localhost:5678)
2. **Go to Settings → API Keys**
3. **Click "Create API Key"**
4. **Name it** `n8n-ops-development`
5. **Copy the key** (starts with `n8n_api_`)

![API Key Creation](docs/images/api-key-creation.png)

### Step 4: Configure Your Environment

```bash
# Copy the example .env file
cp .env.example .env

# Edit the .env file with your API keys
nano .env
```

Add your API keys to the `.env` file:

```
N8N_DEV_API_KEY=n8n_api_your_dev_key_here
N8N_DEV_URL=http://localhost:5678
```

### Step 5: Test Your Connection

```bash
# Test the connection to your development environment
./n8n-ops sync --env development --dry-run
```

You should see:
```
✅ Connected to n8n instance
✅ Found X workflows
```

## Common Workflows

### Sync Workflows from n8n

```bash
# Download workflows from n8n to your local files
./n8n-ops sync --env development
```

### Make Changes

Either:
- Edit the JSON files directly in `workflows/development/`
- Make changes in the n8n UI, then sync again

### Deploy Changes

```bash
# Test your changes first
./n8n-ops validate ./workflows/development/
./n8n-ops deploy --env development --dry-run

# Deploy for real
./n8n-ops deploy --env development
```

### Commit to Git

```bash
git add workflows/development/
git commit -m "feat: update payment workflow"
git push
```

## Troubleshooting

### Can't Connect to n8n?

1. Check your API key in `.env`
2. Verify the n8n URL is correct
3. Make sure n8n is running
4. Try a manual API test:
   ```bash
   curl -H "X-N8N-API-KEY: your_api_key" http://localhost:5678/api/v1/workflows
   ```

### Build Issues?

```bash
go clean
go mod tidy
go build -o n8n-ops main.go
```

## Visual Cheat Sheet

```
┌─────────────────┐         ┌─────────────────┐         ┌─────────────────┐
│                 │  sync   │                 │ deploy  │                 │
│  n8n Instance   │ ◄─────► │  Local Files    │ ◄─────► │  n8n Instance   │
│  (development)  │         │  (Git-managed)  │         │  (production)   │
└─────────────────┘         └─────────────────┘         └─────────────────┘
```

## Next Steps

- Try `./n8n-ops --help` to see all commands
- Check out the [Development Guide](DEVELOPMENT.md) for advanced usage
- Set up GitLab CI/CD for automated deployments

---

**Happy workflow automation! 🚀**