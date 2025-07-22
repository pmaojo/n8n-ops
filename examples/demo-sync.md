# Demo: Sincronización n8n → Git

## 🎯 Proceso Completo de Sincronización

### 1. **Configurar Credenciales**
```bash
# Crear archivo .env con credenciales
cat > .env << EOF
# Development Environment  
N8N_DEV_URL=https://n8n-dev.tuempresa.com
N8N_DEV_API_KEY=n8n_api_xxxxxxxxxxxxxxx

# Staging Environment
N8N_STAGING_URL=https://n8n-staging.tuempresa.com  
N8N_STAGING_API_KEY=n8n_api_yyyyyyyyyyyyyyy

# Production Environment
N8N_PROD_URL=https://n8n-prod.tuempresa.com
N8N_PROD_API_KEY=n8n_api_zzzzzzzzzzzzzzz
EOF

# Cargar variables
source .env
```

### 2. **Sincronización Básica**
```bash
# Sincronizar desde producción
./n8n-ops sync --env production --verbose

# Output esperado:
INFO[2025-07-22 16:00:00] Starting workflow sync                       environment=production
INFO[2025-07-22 16:00:01] Fetching workflows from n8n instance        
INFO[2025-07-22 16:00:02] Found workflows                              count=5
INFO[2025-07-22 16:00:02] Created workflow                             workflow="Customer Onboarding" file=customer_onboarding_1001.json
INFO[2025-07-22 16:00:03] Updated workflow                             workflow="Payment Processing" file=payment_processing_1002.json
INFO[2025-07-22 16:00:03] Sync completed successfully                  created=3 updated=1 skipped=1
```

### 3. **Ver Cambios en Git**
```bash
# Ver archivos sincronizados
ls -la workflows/production/
# customer_onboarding_1001.json
# payment_processing_1002.json  
# email_notifications_1003.json
# _sync_metadata.json

# Ver cambios en Git
git status
# modified:   workflows/production/payment_processing_1002.json
# new file:   workflows/production/customer_onboarding_1001.json

# Ver diferencias específicas
git diff workflows/production/payment_processing_1002.json
```

### 4. **Commit a Git**
```bash
# Validar workflows antes de commit
./n8n-ops validate ./workflows/production/

# Commit cambios
git add workflows/production/
git commit -m "sync: update production workflows from n8n instance

- Updated payment processing timeout settings
- Added new customer onboarding workflow
- Synced from n8n-prod.tuempresa.com at $(date)"

# Push a repositorio
git push origin main
```

## 🚀 Flujo de Desarrollo Completo

### Escenario: Desarrollar nueva funcionalidad

```bash
# 1. Crear feature branch
git checkout -b feature/loyalty-program
git push -u origin feature/loyalty-program

# 2. Sincronizar development
./n8n-ops sync --env development --branch feature/loyalty-program

# 3. Desarrollar en n8n web interface
# - Ir a https://n8n-dev.tuempresa.com
# - Crear/editar "Loyalty Program Workflow"
# - Probar funcionalidad

# 4. Sincronizar cambios de vuelta  
./n8n-ops sync --env development --verbose

# 5. Ver cambios y commit
git diff workflows/development/
git add workflows/development/loyalty_program_1004.json
git commit -m "feat: add loyalty program workflow

- Calculates points based on purchase amount
- Integrates with email notifications  
- Triggers reward voucher generation"

# 6. Push y crear PR
git push origin feature/loyalty-program
# Crear Pull Request en GitLab
```

## 📁 Estructura de Archivos Generados

### Archivos de Workflow
```json
// workflows/production/customer_onboarding_1001.json
{
  "id": "1001",
  "name": "Customer Onboarding",
  "active": true,
  "nodes": [...],
  "connections": {...},
  "syncMetadata": {
    "syncDate": "2025-07-22T16:00:02Z",
    "environment": "production", 
    "gitCommit": "a1b2c3d4e5f6",
    "syncedBy": "developer@tuempresa.com"
  }
}
```

### Metadata de Sincronización
```json
// workflows/production/_sync_metadata.json
{
  "lastSync": "2025-07-22T16:00:03Z",
  "environment": "production",
  "totalWorkflows": 5,
  "created": 3,
  "updated": 1, 
  "skipped": 1,
  "gitCommit": "a1b2c3d4e5f6",
  "pipelineId": "12345",
  "syncedBy": "developer@tuempresa.com"
}
```

## 🛠️ Opciones Avanzadas

### Sincronización Forzada
```bash
# Sobrescribir cambios locales
./n8n-ops sync --env production --force

# Cambiar directorio de salida
./n8n-ops sync --env production --output ./backups/prod-workflows
```

### Script Automatizado
```bash
# Usar el script helper
./examples/sync-workflow.sh production

# Sincronizar todos los ambientes
./examples/sync-workflow.sh

# Output:
# 🔄 Syncing all environments...
# 
# 📥 Syncing development environment...
# ✅ Successfully synced development workflows
# 
# 📥 Syncing staging environment...  
# ✅ Successfully synced staging workflows
#
# 📥 Syncing production environment...
# ✅ Successfully synced production workflows
```

## 💡 Casos de Uso Prácticos

### 1. **Backup Antes de Deploy**
```bash
# Sincronizar estado actual antes de hacer deploy
./n8n-ops sync --env production --output ./backup/pre-deploy-$(date +%Y%m%d)
./n8n-ops deploy --env production
```

### 2. **Comparar Ambientes**
```bash
# Sincronizar ambos ambientes
./n8n-ops sync --env staging
./n8n-ops sync --env production

# Comparar diferencias
diff -r workflows/staging/ workflows/production/
```

### 3. **Migración de Ambiente**
```bash
# Copiar workflows de producción a staging
./n8n-ops sync --env production
cp -r workflows/production/* workflows/staging/
./n8n-ops deploy --env staging --dry-run
./n8n-ops deploy --env staging
```

El comando sync te permite mantener tu código de workflows siempre actualizado con los cambios que hagas en la interfaz web de n8n, garantizando que todo quede versionado en Git.