# Guía de Sincronización n8n → Git

## 🔄 Comando Sync - Importar desde n8n

### Propósito
El comando `sync` descarga workflows desde tu instancia n8n y los guarda como archivos JSON en tu proyecto local, listos para commit a Git.

### Uso Básico

```bash
# Sincronizar desde producción
./n8n-ops sync --env production

# Sincronizar desde staging
./n8n-ops sync --env staging

# Sincronizar desde desarrollo  
./n8n-ops sync --env development

# Modo verbose para ver detalles
./n8n-ops sync --env production --verbose
```

## 📁 Estructura de Directorios

Después del sync, tendrás:

```
workflows/
├── development/
│   ├── customer-onboarding.json
│   ├── payment-processing.json
│   └── email-notifications.json
├── staging/
│   ├── customer-onboarding.json
│   └── payment-processing.json
└── production/
    ├── customer-onboarding.json
    └── payment-processing.json
```

## 🎯 Flujo de Trabajo con Sync

### 1. **Desarrollar en n8n Web Interface**
- Crea/edita workflows en la interfaz web de n8n
- Prueba y valida funcionalidad
- Cuando esté listo, sincroniza a local

### 2. **Sincronizar a Local**
```bash
# Descargar workflows modificados
./n8n-ops sync --env development --verbose

# Ver qué cambió
git status
git diff
```

### 3. **Commit a Git**
```bash
# Revisar cambios
./n8n-ops validate ./workflows/development/

# Commit cambios
git add workflows/development/
git commit -m "feat: update customer onboarding workflow"
git push origin develop
```

## 📋 Configuración Requerida

### Variables de Entorno
Necesitas configurar las credenciales para cada ambiente:

```bash
# Development
export N8N_URL_DEV="https://n8n-dev.tuempresa.com"
export N8N_API_KEY_DEV="tu_api_key_development"

# Staging  
export N8N_URL_STAGING="https://n8n-staging.tuempresa.com"
export N8N_API_KEY_STAGING="tu_api_key_staging"

# Production
export N8N_URL_PROD="https://n8n-prod.tuempresa.com" 
export N8N_API_KEY_PROD="tu_api_key_production"
```

### En GitLab CI/CD
Las mismas variables se configuran en GitLab → Settings → CI/CD → Variables

## 🔍 Ejemplo Práctico

### Escenario: Modificar Workflow en Producción

```bash
# 1. Sincronizar estado actual
./n8n-ops sync --env production
./n8n-ops status --env production

# 2. Crear branch para cambios  
git checkout -b hotfix/update-payment-timeout
git add workflows/production/
git commit -m "sync: current production workflows"

# 3. Hacer cambios en n8n web interface
# (editar workflow en https://n8n-prod.tuempresa.com)

# 4. Sincronizar cambios
./n8n-ops sync --env production --verbose

# 5. Ver diferencias
git diff workflows/production/payment-processing.json

# 6. Validar y commit
./n8n-ops validate ./workflows/production/
git add workflows/production/payment-processing.json
git commit -m "fix: increase payment timeout to 30 seconds"
git push origin hotfix/update-payment-timeout
```

## 📊 Output Esperado del Sync

```bash
$ ./n8n-ops sync --env production --verbose

INFO[2025-07-22 15:30:00] Connecting to n8n instance                   env=production url=https://n8n-prod.tuempresa.com
INFO[2025-07-22 15:30:01] Found 3 workflows in production              
INFO[2025-07-22 15:30:01] Downloading: Customer Onboarding Workflow    id=1001
INFO[2025-07-22 15:30:02] Downloading: Payment Processing System       id=1002  
INFO[2025-07-22 15:30:03] Downloading: Email Notifications             id=1003
INFO[2025-07-22 15:30:04] Saved: ./workflows/production/customer-onboarding.json
INFO[2025-07-22 15:30:04] Saved: ./workflows/production/payment-processing.json
INFO[2025-07-22 15:30:04] Saved: ./workflows/production/email-notifications.json

✅ Successfully synced 3 workflows from production
📁 Files saved to: ./workflows/production/
💡 Tip: Run 'git status' to see changes, then commit to Git
```

## 🛠️ Opciones Avanzadas

### Sync Selectivo
```bash
# Solo workflows activos
./n8n-ops sync --env production --active-only

# Solo workflows específicos (por ID)
./n8n-ops sync --env production --workflow-ids 1001,1002

# Incluir workflows inactivos  
./n8n-ops sync --env production --include-inactive
```

### Formato de Archivos
```bash
# Con metadata extendida
./n8n-ops sync --env production --include-metadata

# Formato compacto (sin espacios)
./n8n-ops sync --env production --compact

# Backup antes de overwrite
./n8n-ops sync --env production --backup
```

## 🔄 Sincronización Bidireccional

### n8n → Local (Sync)
```bash
# Descargar cambios desde n8n
./n8n-ops sync --env production
```

### Local → n8n (Deploy)  
```bash
# Subir cambios locales a n8n
./n8n-ops deploy --env production
```

### Verificar Sincronización
```bash
# Ver estado después del sync
./n8n-ops status --env production
```

## ⚠️ Consideraciones Importantes

### 1. **Conflictos de Versión**
- Si alguien más editó workflows en n8n, el sync puede sobrescribir cambios locales
- Siempre hacer backup antes de sync importante
- Usar `git stash` para guardar cambios locales temporalmente

### 2. **Credenciales Sensibles**
- Los workflows pueden contener API keys o passwords
- Revisar archivos JSON antes de commit
- Usar variables de entorno en lugar de hardcodear secretos

### 3. **Workflows Grandes**
- Workflows con muchos nodos pueden generar archivos JSON grandes
- Considerar `.gitignore` para workflows de testing temporal

## 🎯 Mejores Prácticas

### 1. **Sync Regular**
```bash
# Al inicio del día
./n8n-ops sync --env development

# Antes de hacer cambios locales
./n8n-ops sync --env production --backup
```

### 2. **Verificación Post-Sync**
```bash
# Verificar integridad
./n8n-ops validate ./workflows/production/

# Revisar cambios
git diff --name-only
git diff workflows/production/
```

### 3. **Automatización en CI/CD**
```yaml
# En .gitlab-ci.yml
sync-production:
  script:
    - ./n8n-ops sync --env production
    - git add workflows/production/
    - git commit -m "sync: latest production workflows [skip ci]"
    - git push origin sync-$(date +%Y%m%d)
```

El comando sync es tu puente entre n8n y Git, permitiendo versionar y colaborar en workflows de forma profesional.