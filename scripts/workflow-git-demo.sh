#!/bin/bash

echo "🚨 DEMO: Workflows Editados No Commiteados"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}📋 CREANDO SIMULACIÓN DE WORKFLOWS MODIFICADOS${NC}"
echo "============================================="
echo ""

# Create workflows directory structure
mkdir -p workflows/development
mkdir -p workflows/staging  
mkdir -p workflows/production

echo "Creando archivos de workflow de ejemplo..."

# Create some workflow files
cat > workflows/development/customer-onboarding.json << 'EOF'
{
  "id": "customer_onboarding_001",
  "name": "Customer Onboarding Process",
  "nodes": [
    {
      "id": "trigger_node",
      "type": "n8n-nodes-base.manualTrigger",
      "position": [250, 300]
    },
    {
      "id": "email_node",
      "type": "n8n-nodes-base.emailSend",
      "credentials": {
        "smtp": {
          "id": "smtp_dev_001",
          "name": "SMTP Development"
        }
      },
      "position": [450, 300]
    }
  ],
  "connections": {},
  "settings": {
    "executionOrder": "v1"
  },
  "updatedAt": "2025-01-22T10:30:00.000Z"
}
EOF

cat > workflows/development/email-notifications.json << 'EOF'
{
  "id": "email_notifications_002", 
  "name": "Email Notification System",
  "nodes": [
    {
      "id": "webhook_trigger",
      "type": "n8n-nodes-base.webhook",
      "position": [250, 300]
    },
    {
      "id": "postgres_node", 
      "type": "n8n-nodes-base.postgres",
      "credentials": {
        "postgres": {
          "id": "postgres_dev_001",
          "name": "PostgreSQL Development"
        }
      },
      "position": [450, 300]
    }
  ],
  "connections": {},
  "settings": {
    "executionOrder": "v1"
  },
  "updatedAt": "2025-01-22T11:45:00.000Z"
}
EOF

cat > workflows/production/payment-processing.json << 'EOF'
{
  "id": "payment_processing_003",
  "name": "Payment Processing Workflow", 
  "nodes": [
    {
      "id": "api_trigger",
      "type": "n8n-nodes-base.httpRequest",
      "position": [250, 300]
    },
    {
      "id": "stripe_node",
      "type": "n8n-nodes-base.stripe", 
      "credentials": {
        "stripe": {
          "id": "stripe_prod_001",
          "name": "Stripe Production"
        }
      },
      "position": [450, 300]
    }
  ],
  "connections": {},
  "settings": {
    "executionOrder": "v1"
  },
  "updatedAt": "2025-01-22T09:15:00.000Z"
}
EOF

echo "✅ Workflows creados"
echo ""

echo -e "${YELLOW}🔧 SIMULANDO ESTADO DE GIT${NC}"
echo "============================="
echo ""

# Initialize git if not already
if [ ! -d ".git" ]; then
    echo "Inicializando repositorio Git..."
    git init
    git config user.email "test@n8n-ops.com"
    git config user.name "n8n-ops Demo"
fi

# Add some files to git but leave others uncommitted
echo "Agregando algunos workflows al staging..."
git add workflows/production/payment-processing.json

echo "Modificando archivos para simular cambios sin commitear..."

# Modify a file that's already tracked
echo '// Modified at '$(date) >> workflows/development/customer-onboarding.json

# Create an untracked file
cat > workflows/development/new-workflow.json << 'EOF'
{
  "id": "new_workflow_004",
  "name": "New Workflow Created",
  "nodes": [],
  "connections": {},
  "updatedAt": "2025-01-22T12:00:00.000Z"
}
EOF

echo ""
echo -e "${BLUE}📊 ESTADO ACTUAL DE GIT:${NC}"
echo "========================"
echo ""

git status --porcelain

echo ""
echo -e "${RED}🚨 WORKFLOWS NO COMMITEADOS DETECTADOS:${NC}"
echo "==========================================="
echo ""

# Check for workflow files in git status
git status --porcelain | grep workflows/ | while read line; do
    status=${line:0:2}
    file=${line:3}
    
    if [[ $file == *".json" ]]; then
        case $status in
            "M ") echo "📝 MODIFICADO: $file" ;;
            " M") echo "📝 MODIFICADO: $file" ;;
            "A ") echo "➕ AGREGADO: $file" ;;
            "??") echo "❓ SIN RASTREAR: $file" ;;
            "D ") echo "🗑️ ELIMINADO: $file" ;;
        esac
    fi
done

echo ""
echo -e "${YELLOW}💡 COMANDOS DE n8n-ops PARA DETECCIÓN:${NC}"
echo "========================================"
echo ""

echo "VERIFICAR CAMBIOS SIN COMMITEAR:"
echo "n8n-ops status --check-uncommitted"
echo ""

echo "VER ESTADO COMPLETO:"
echo "n8n-ops status"
echo ""

echo "BLOQUEAR SYNC SI HAY CAMBIOS:"
echo "n8n-ops sync --env development  # Bloqueará si hay cambios"
echo ""

echo "FORZAR SYNC IGNORANDO CAMBIOS:"
echo "n8n-ops sync --env development --force"
echo ""

echo -e "${YELLOW}🎯 FUNCIONALIDAD EN WEB UI:${NC}"
echo "==========================="
echo ""

echo "LA WEB UI EN http://localhost:5000 MOSTRARÁ:"
echo ""
echo "📊 Dashboard Principal:"
echo "  • ⚠️ Banner rojo: '2 Uncommitted Workflow Changes'"
echo "  • 🚨 Card especial: Lista de workflows modificados"
echo "  • 📋 Overview: Estado de salud general"
echo ""

echo "🔧 Funciones Interactivas:"
echo "  • Botón 'Commit Changes' - commits automático"
echo "  • Warning antes de sync operations"
echo "  • Status visual de cada workflow"
echo ""

echo -e "${GREEN}✅ RESULTADOS ESPERADOS:${NC}"
echo "=========================="
echo ""

echo "COMANDO: n8n-ops status --check-uncommitted"
echo "SALIDA:"
echo ""
echo "🔍 Checking for Uncommitted Workflow Changes"
echo "==========================================="
echo ""
echo "🚨 WARNING: Uncommitted Workflow Changes Detected!"
echo "====================================="
echo ""
echo "📝 • Customer Onboarding Process (modified)"
echo "   Environment: development | File: workflows/development/customer-onboarding.json"
echo ""
echo "❓ • New Workflow Created (untracked)"
echo "   Environment: development | File: workflows/development/new-workflow.json"
echo ""
echo "💡 Recommendation:"
echo "   git add ."
echo "   git commit -m \"Update 2 workflow(s)\""
echo "   git push origin main"
echo ""

echo -e "${GREEN}🎯 BENEFICIOS DEL SISTEMA:${NC}"
echo "=========================="
echo ""

echo "PREVENCIÓN DE PÉRDIDA DE DATOS:"
echo "• Detecta workflows editados en n8n UI"
echo "• Avisa antes de operaciones de sync" 
echo "• Evita sobrescritura accidental"
echo ""

echo "COLABORACIÓN SEGURA:"
echo "• Todos los cambios deben estar en Git"
echo "• Historial completo de modificaciones"
echo "• Resolución de conflictos controlada"
echo ""

echo "FLUJO DE TRABAJO MEJORADO:"
echo "• Visual warnings en Web UI"
echo "• Comandos CLI informativos"
echo "• Integración con DevOps pipelines"
echo ""

echo -e "${BLUE}🚀 SIGUIENTES PASOS:${NC}"
echo "=================="
echo ""

echo "1. Abrir Web UI: http://localhost:5000"
echo "2. Ver banner de warning en dashboard"
echo "3. Probar: n8n-ops status --check-uncommitted"
echo "4. Commitear cambios: git add . && git commit -m \"Save workflows\""
echo "5. Verificar: Web UI ya no muestra warnings"
echo ""

echo "¡Sistema de detección de workflows no commiteados funcionando!"