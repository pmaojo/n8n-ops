#!/bin/bash

echo "🔐 n8n CREDENTIAL SECURITY SYSTEM DEMO"
echo "======================================"
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔒 CÓMO FUNCIONA LA API DE CREDENCIALES DE n8n${NC}"
echo "=================================================="
echo ""

echo "1. ARQUITECTURA DE SEGURIDAD DE n8n:"
echo "   • Las credenciales se encriptan en la base de datos de n8n"
echo "   • Los workflows solo contienen IDs de credenciales, no valores"
echo "   • Cada ambiente tiene sus propias credenciales"
echo ""

echo "2. ESTRUCTURA DE WORKFLOW JSON:"
echo '   {
     "nodes": [
       {
         "id": "smtp_node",
         "type": "n8n-nodes-base.emailSend",
         "credentials": {
           "smtp": {
             "id": "smtp_prod_001",      ← Solo el ID
             "name": "SMTP Production"   ← Nombre descriptivo
           }
         }
       }
     ]
   }'
echo ""

echo "3. API ENDPOINTS DE n8n:"
echo "   GET  /api/v1/credentials        # Listar credenciales"
echo "   GET  /api/v1/credentials/{id}   # Obtener credencial específica"
echo "   POST /api/v1/credentials        # Crear nueva credencial"
echo "   PUT  /api/v1/credentials/{id}   # Actualizar credencial"
echo ""

echo -e "${YELLOW}🛡️ ESTRATEGIA DE SEGURIDAD DE n8n-ops${NC}"
echo "=========================================="
echo ""

echo "PRINCIPIOS DE SEGURIDAD:"
echo "✅ Nunca almacenar credenciales en archivos"
echo "✅ Solo usar variables de entorno"
echo "✅ Mapeo automático por ambiente"
echo "✅ Validación antes de despliegue"
echo ""

echo "CONVENCIÓN DE NAMING:"
echo "{SERVICE}_{FIELD}_{ENVIRONMENT}"
echo ""
echo "Ejemplos:"
echo "• SMTP_HOST_DEVELOPMENT"
echo "• POSTGRES_PASSWORD_PRODUCTION"
echo "• STRIPE_SECRET_KEY_STAGING"
echo "• AWS_ACCESS_KEY_ID_PRODUCTION"
echo ""

echo -e "${YELLOW}🔄 FLUJO DE TRABAJO COMPLETO${NC}"
echo "============================="
echo ""

echo "PASO 1: DESARROLLO LOCAL"
echo "------------------------"
echo "export N8N_URL_DEVELOPMENT=\"http://localhost:5678\""
echo "export N8N_API_KEY_DEVELOPMENT=\"n8n_api_dev_key\""
echo "export SMTP_HOST_DEVELOPMENT=\"smtp.mailtrap.io\""
echo "export SMTP_USER_DEVELOPMENT=\"dev_user\""
echo "export SMTP_PASSWORD_DEVELOPMENT=\"dev_password\""
echo ""
echo "# Crear workflow con credenciales de desarrollo"
echo "n8n-ops branch create email-notifications"
echo "# Workflow usa credencial \"smtp_dev_001\" automáticamente"
echo ""

echo "PASO 2: DESPLIEGUE A STAGING"
echo "----------------------------"
echo "export N8N_URL_STAGING=\"https://n8n-staging.company.com\""
echo "export N8N_API_KEY_STAGING=\"n8n_api_staging_key\""
echo "export SMTP_HOST_STAGING=\"smtp-staging.company.com\""
echo "export SMTP_USER_STAGING=\"staging_user\""
echo "export SMTP_PASSWORD_STAGING=\"staging_password\""
echo ""
echo "# n8n-ops mapea automáticamente las credenciales"
echo "n8n-ops sync --to-n8n --env staging"
echo "# Workflow usa credencial \"smtp_staging_001\" automáticamente"
echo ""

echo "PASO 3: PRODUCCIÓN"
echo "------------------"
echo "export N8N_URL_PRODUCTION=\"https://n8n.company.com\""
echo "export N8N_API_KEY_PRODUCTION=\"n8n_api_prod_key\""
echo "export SMTP_HOST_PRODUCTION=\"smtp.sendgrid.net\""
echo "export SMTP_USER_PRODUCTION=\"apikey\""
echo "export SMTP_PASSWORD_PRODUCTION=\"SG.real_sendgrid_key\""
echo ""
echo "# Despliegue seguro a producción"
echo "n8n-ops sync --to-n8n --env production"
echo "# Workflow usa credencial \"smtp_prod_001\" automáticamente"
echo ""

echo -e "${YELLOW}🔧 COMANDOS DE n8n-ops PARA CREDENCIALES${NC}"
echo "==========================================="
echo ""

echo "LISTAR CREDENCIALES:"
echo "n8n-ops credentials list --env development"
echo ""
echo "Salida:"
echo "🔐 Credential Mappings - development Environment"
echo "====================================="
echo ""
echo "📋 N8N API Credentials:"
echo "  N8N_URL_DEVELOPMENT = http://localhost:5678"
echo "  N8N_API_KEY_DEVELOPMENT = n8n_***dev"
echo ""
echo "🔧 Workflow Node Credentials:"
echo "  SMTP:"
echo "    host = smtp.mailtrap.io (SMTP_HOST_DEVELOPMENT)"
echo "    user = dev_user (SMTP_USER_DEVELOPMENT)"
echo "    password = *** (SMTP_PASSWORD_DEVELOPMENT)"
echo ""

echo "VALIDAR CREDENCIALES:"
echo "n8n-ops credentials validate --env production"
echo ""
echo "Salida:"
echo "✅ Validating Credentials - production Environment"
echo "==========================================="
echo ""
echo "📊 Validation Results:"
echo "  Required: 12"
echo "  Present:  10"
echo "  Missing:  2"
echo ""
echo "❌ Missing Credentials:"
echo "  • STRIPE_SECRET_KEY_PRODUCTION"
echo "  • AWS_SECRET_ACCESS_KEY_PRODUCTION"
echo ""

echo "GENERAR TEMPLATE:"
echo "n8n-ops credentials template --env production"
echo ""
echo "Salida:"
echo "# n8n-ops Environment Variables - PRODUCTION"
echo "# Copy these variables to your environment"
echo ""
echo "# N8N API Configuration"
echo "export N8N_URL_PRODUCTION=\"https://n8n.company.com\""
echo "export N8N_API_KEY_PRODUCTION=\"your_production_api_key_here\""
echo ""

echo -e "${YELLOW}🎯 MAPEO AUTOMÁTICO DE CREDENCIALES${NC}"
echo "==================================="
echo ""

echo "ALGORITMO DE MAPEO:"
echo "1. Analizar Workflow: Identificar nodos que requieren credenciales"
echo "2. Mapear por Ambiente: Determinar ID de credencial correcto"
echo "3. Validar Existencia: Verificar que la credencial existe en n8n"
echo "4. Actualizar Referencias: Cambiar IDs en el workflow JSON"
echo ""

echo "EJEMPLO DE MAPEO:"
echo "Desarrollo:  smtp_dev_001"
echo "Staging:     smtp_staging_001"
echo "Producción:  smtp_prod_001"
echo ""

echo -e "${YELLOW}🔐 SEGURIDAD EN CI/CD${NC}"
echo "===================="
echo ""

echo "GITLAB CI/CD VARIABLES:"
echo "# .gitlab-ci.yml"
echo "deploy_staging:"
echo "  script:"
echo "    - n8n-ops credentials validate --env staging"
echo "    - n8n-ops sync --to-n8n --env staging"
echo "  variables:"
echo "    N8N_URL_STAGING: \"https://n8n-staging.company.com\""
echo "  only:"
echo "    - staging"
echo ""

echo "deploy_production:"
echo "  script:"
echo "    - n8n-ops credentials validate --env production"
echo "    - n8n-ops sync --to-n8n --env production"
echo "  variables:"
echo "    N8N_URL_PRODUCTION: \"https://n8n.company.com\""
echo "  only:"
echo "    - production"
echo "  when: manual"
echo ""

echo -e "${YELLOW}🚨 MEJORES PRÁCTICAS DE SEGURIDAD${NC}"
echo "=================================="
echo ""

echo -e "${GREEN}✅ HACER:${NC}"
echo "• Usar variables de entorno para todas las credenciales"
echo "• Rotar credenciales regularmente"
echo "• Validar antes de cada despliegue"
echo "• Usar diferentes credenciales por ambiente"
echo "• Logs sin credenciales (enmascarar valores)"
echo ""

echo -e "${RED}❌ NO HACER:${NC}"
echo "• Nunca commit credenciales en Git"
echo "• Nunca usar mismas credenciales en dev/prod"
echo "• Nunca loggear valores de credenciales"
echo "• Nunca compartir credenciales por email/chat"
echo "• Nunca hardcodear credenciales en workflows"
echo ""

echo -e "${BLUE}💡 CASOS DE USO REALES${NC}"
echo "====================="
echo ""

echo "CASO 1: SMTP PARA EMAILS"
echo "• Development: Mailtrap (testing)"
echo "• Staging: SMTP interno (testing)"
echo "• Production: SendGrid (real emails)"
echo ""

echo "CASO 2: BASE DE DATOS"
echo "• Development: SQLite local"
echo "• Staging: PostgreSQL staging"
echo "• Production: PostgreSQL clustered"
echo ""

echo "CASO 3: PAYMENTS (STRIPE)"
echo "• Development: Test keys (sk_test_...)"
echo "• Staging: Test keys (sk_test_...)"
echo "• Production: Live keys (sk_live_...)"
echo ""

echo "CASO 4: CLOUD STORAGE (AWS S3)"
echo "• Development: Bucket dev"
echo "• Staging: Bucket staging"
echo "• Production: Bucket prod"
echo ""

echo -e "${GREEN}🏆 RESUMEN EJECUTIVO${NC}"
echo "===================="
echo ""

echo "El sistema de credenciales de n8n-ops:"
echo ""
echo "• SEGURIDAD TOTAL: Solo variables de entorno, nunca archivos"
echo "• MAPEO AUTOMÁTICO: Credenciales correctas por ambiente"
echo "• VALIDACIÓN PREVIA: Verificación antes de despliegue"
echo "• TRAZABILIDAD: Logs seguros sin exponer credenciales"
echo "• ESCALABILIDAD: Soporte para múltiples ambientes y equipos"
echo ""
echo "¡Perfecto para uso empresarial con máxima seguridad!"
echo ""

echo -e "${BLUE}🚀 PRÓXIMOS PASOS${NC}"
echo "=================="
echo ""
echo "1. Configurar variables de entorno por ambiente"
echo "2. Validar credenciales existentes: n8n-ops credentials validate"
echo "3. Generar template: n8n-ops credentials template"
echo "4. Configurar CI/CD con variables seguras"
echo "5. Capacitar al equipo en mejores prácticas"
echo ""
echo "🎯 ¡Listo para manejo seguro de credenciales!"