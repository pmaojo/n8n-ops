#!/bin/bash

# Script para crear repositorio público completamente limpio
# n8n-ops CLI Tool

echo "🚀 Creando repositorio público n8n-ops..."

# 1. Crear directorio para repositorio público
mkdir -p ../n8n-ops-public
cd ../n8n-ops-public

echo "📁 Directorio creado: n8n-ops-public"

# 2. Copiar archivos esenciales (excluir .git y archivos temporales)
echo "📋 Copiando archivos del proyecto..."

# Código principal
cp ../workspace/main.go .
cp ../workspace/go.mod .
cp ../workspace/go.sum .
cp ../workspace/Makefile .
cp ../workspace/LICENSE .
cp ../workspace/README.md .

# Directorios principales
cp -r ../workspace/cmd ./
cp -r ../workspace/internal ./
cp -r ../workspace/config ./
cp -r ../workspace/workflows ./
cp -r ../workspace/tests ./
cp -r ../workspace/mock-n8n-server ./
cp -r ../workspace/scripts ./
cp -r ../workspace/docker ./

# Configuración y CI/CD
cp ../workspace/.gitlab-ci.yml .
cp ../workspace/config.example.yaml .
cp ../workspace/.env.example .
cp ../workspace/.gitignore .

# Documentación profesional
cp ../workspace/USER_STORIES.md .
cp ../workspace/CONTRIBUTING.md .
cp ../workspace/DEVELOPMENT.md .
cp ../workspace/SECURITY.md .
cp ../workspace/QUICK_START.md .
cp ../workspace/PRODUCTION_SETUP.md .
cp ../workspace/DEPLOYMENT_VPS.md .
cp ../workspace/WORKFLOW_GUIDE.md .
cp ../workspace/WORKFLOW_EXAMPLES.md .
cp ../workspace/GITLAB_API_INTEGRATION.md .
cp ../workspace/CREDENTIAL_SECURITY.md .
cp ../workspace/CHANGELOG.md .
cp ../workspace/VERSIONING_GUIDE.md .
cp ../workspace/ZERO_DOWNTIME_ARCHITECTURE.md .

echo "✅ Archivos copiados correctamente"

# 3. Inicializar git limpio
echo "🔧 Inicializando repositorio Git..."
git init
git branch -M main

# 4. Crear primer commit profesional
echo "📝 Creando commit inicial..."
git add .

git commit -m "feat: n8n-ops CLI tool for enterprise workflow management

Complete CLI solution for managing n8n workflows across multiple environments
with GitLab integration, automated synchronization, and comprehensive tooling.

Core Features:
• Multi-environment sync (development, staging, production)  
• GitLab CI/CD integration with branch tracking
• Daemon mode with real-time file watching
• Comprehensive credential management across environments
• Automated backups with version control
• Professional CLI with 10+ commands
• JSON validation with business rules
• Cross-platform deployment support

Technical Architecture:
• Go 1.19+ with clean modular design
• Cobra CLI framework for professional UX
• SQLite for local tracking and metadata
• n8n REST API integration (100% compatibility)
• GitOps-level visibility and control
• Comprehensive logging system (5 levels)
• Mock server for development testing
• Enterprise-grade error handling

This tool enables organizations to manage n8n workflows with the same
rigor and automation as modern DevOps practices, providing GitLab
integration, multi-environment support, and comprehensive audit trails."

echo "🎉 Repositorio público creado exitosamente!"
echo ""
echo "📍 Directorio: ../n8n-ops-public"
echo "🔗 Para conectar a GitLab:"
echo "   cd ../n8n-ops-public"
echo "   git remote add origin https://gitlab.com/tu-usuario/n8n-ops.git"
echo "   git push -u origin main"
echo ""
echo "✨ ¡Listo para ser público!"