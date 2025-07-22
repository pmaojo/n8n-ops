# GitLab CI/CD + n8n API Integration

## 🔄 Cómo Funciona la Integración

### GitLab CI/CD → n8n API
El pipeline de GitLab ejecuta el CLI `n8n-ops` que se conecta directamente a las APIs de n8n:

```
GitLab CI/CD Pipeline
        ↓
    n8n-ops CLI
        ↓
    n8n API Calls
        ↓
n8n Instance (Dev/Staging/Prod)
```

## 🚀 Pipeline Stages y API Calls

### 1. **Build Stage** (Sin API)
```yaml
build-cli:
  script:
    - go build -o n8n-ops main.go  # Solo compila, no llama APIs
```

### 2. **Validate Stage** (API de lectura)
```yaml
validate-workflows:
  script:
    - ./n8n-ops validate ./workflows/ --verbose
    # Internamente puede llamar GET /api/v1/workflows para validar IDs
```

### 3. **Deploy Stage** (API completa)
```yaml
deploy-production:
  before_script:
    - export N8N_API_KEY="${N8N_PROD_API_KEY}"
    - export N8N_URL="${N8N_PROD_URL}"
  script:
    - ./n8n-ops deploy --env production --verbose
```

## 📡 API Calls Específicos del CLI

### Durante `n8n-ops sync`
```http
# 1. Autenticación
GET https://n8n-prod.tuempresa.com/api/v1/me
Headers: X-N8N-API-KEY: n8n_api_xxxxxxxx

# 2. Listar workflows
GET https://n8n-prod.tuempresa.com/api/v1/workflows
Headers: X-N8N-API-KEY: n8n_api_xxxxxxxx

# 3. Descargar cada workflow
GET https://n8n-prod.tuempresa.com/api/v1/workflows/1001
Headers: X-N8N-API-KEY: n8n_api_xxxxxxxx
```

### Durante `n8n-ops deploy`
```http
# 1. Verificar workflow existe
GET https://n8n-prod.tuempresa.com/api/v1/workflows/1001
Headers: X-N8N-API-KEY: n8n_api_xxxxxxxx

# 2. Actualizar workflow (si existe)
PUT https://n8n-prod.tuempresa.com/api/v1/workflows/1001
Headers: X-N8N-API-KEY: n8n_api_xxxxxxxx
Content-Type: application/json
Body: {workflow JSON completo}

# 3. Crear workflow (si no existe)
POST https://n8n-prod.tuempresa.com/api/v1/workflows
Headers: X-N8N-API-KEY: n8n_api_xxxxxxxx
Content-Type: application/json
Body: {workflow JSON completo}

# 4. Activar/desactivar workflow
POST https://n8n-prod.tuempresa.com/api/v1/workflows/1001/activate
Headers: X-N8N-API-KEY: n8n_api_xxxxxxxx
```

### Durante `n8n-ops status`
```http
# 1. Obtener lista de workflows
GET https://n8n-prod.tuempresa.com/api/v1/workflows
Headers: X-N8N-API-KEY: n8n_api_xxxxxxxx

# 2. Obtener detalles de cada workflow
GET https://n8n-prod.tuempresa.com/api/v1/workflows/1001
Headers: X-N8N-API-KEY: n8n_api_xxxxxxxx
```

## 🔑 Configuración de Credenciales en GitLab

### Variables de CI/CD
En GitLab → Project → Settings → CI/CD → Variables:

```bash
# Development
N8N_DEV_URL = https://n8n-dev.tuempresa.com
N8N_DEV_API_KEY = n8n_api_dev_xxxxxxxxxxxxxxxxx

# Staging  
N8N_STAGING_URL = https://n8n-staging.tuempresa.com
N8N_STAGING_API_KEY = n8n_api_staging_yyyyyyyyyyy

# Production (Protected)
N8N_PROD_URL = https://n8n-prod.tuempresa.com  
N8N_PROD_API_KEY = n8n_api_prod_zzzzzzzzzzzzz
```

### Configuración por Ambiente
```yaml
# Deploy Development (automático)
deploy-development:
  before_script:
    - export N8N_API_KEY="${N8N_DEV_API_KEY}"
    - export N8N_URL="${N8N_DEV_URL}"
  script:
    - ./n8n-ops deploy --env development --verbose
    # API calls a https://n8n-dev.tuempresa.com
  rules:
    - if: '$CI_COMMIT_BRANCH == "develop"'

# Deploy Production (manual)
deploy-production:
  before_script:
    - export N8N_API_KEY="${N8N_PROD_API_KEY}"  
    - export N8N_URL="${N8N_PROD_URL}"
  script:
    - ./n8n-ops deploy --env production --verbose
    # API calls a https://n8n-prod.tuempresa.com  
  when: manual
  rules:
    - if: '$CI_COMMIT_BRANCH == "main"'
```

## 📊 Flow Completo con APIs

### Ejemplo: Push a Branch `main`

```
1. Git Push to main
   ↓
2. GitLab CI/CD triggers pipeline
   ↓
3. Build Stage: Compila n8n-ops
   ↓
4. Validate Stage: 
   - Valida archivos JSON locales
   - Opcionalmente: GET /api/v1/workflows (verificar IDs)
   ↓
5. Deploy Production (Manual):
   - Usuario hace clic "Deploy" en GitLab
   - Pipeline ejecuta: ./n8n-ops deploy --env production
   - CLI hace API calls:
     * GET /api/v1/workflows (listar existentes)
     * PUT /api/v1/workflows/1001 (actualizar)
     * POST /api/v1/workflows (crear nuevos)
     * POST /api/v1/workflows/1001/activate
   ↓
6. Verificación:
   - ./n8n-ops status --env production
   - GET /api/v1/workflows (verificar estado final)
```

## 🛡️ Seguridad y Autenticación

### API Keys n8n
```bash
# Generar API key en n8n:
# 1. ir a Settings → API Keys
# 2. Create API Key → Copy token
# 3. Usar en GitLab CI/CD variables
```

### Variables Protegidas
```yaml
# En GitLab CI/CD Variables:
N8N_PROD_API_KEY:
  - Protected: ✓ (solo branches protegidas) 
  - Masked: ✓ (ocultar en logs)
  - Environment scope: production
```

## 📈 Monitoreo de API Calls

### Logs en GitLab CI/CD
```bash
INFO[2025-07-22 16:30:00] Connecting to n8n API                        url=https://n8n-prod.tuempresa.com
INFO[2025-07-22 16:30:01] Found 5 workflows in n8n instance           
INFO[2025-07-22 16:30:02] Updating workflow via API                   workflow="Payment Processing" id=1002
INFO[2025-07-22 16:30:03] Creating new workflow via API               workflow="Loyalty Program" 
INFO[2025-07-22 16:30:04] Activating workflow                          workflow="Loyalty Program" id=1005
INFO[2025-07-22 16:30:05] Deploy completed successfully               updated=2 created=1
```

### Logs en n8n
```bash
# n8n logs mostrarán:
API call: GET /api/v1/workflows - 200 OK
API call: PUT /api/v1/workflows/1002 - 200 OK  
API call: POST /api/v1/workflows - 201 Created
API call: POST /api/v1/workflows/1005/activate - 200 OK
```

## 🎯 Ventajas de esta Arquitectura

### 1. **Automatización Completa**
- Push a Git → Deploy automático a n8n
- Sin intervención manual en producción

### 2. **Multi-ambiente**  
- Mismas APIs, diferentes instancias
- Credenciales separadas por ambiente

### 3. **Auditabilidad**
- Todos los API calls registrados en GitLab
- Trazabilidad completa Git → n8n

### 4. **Rollback Rápido**
- API calls para revertir a versión anterior
- Backup automático antes de deploy

La integración es **completamente transparente**: GitLab ejecuta tu CLI, y el CLI maneja todas las llamadas API a n8n. Tú solo haces `git push` y el resto es automático.