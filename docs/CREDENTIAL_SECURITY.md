# n8n Credential Security & API Integration

## 🔒 Cómo Funciona la API de Credenciales de n8n

> **Nota:** Las variables de entorno deben seguir el esquema `N8N_API_KEY_<ENV>` y `N8N_URL_<ENV>`. Por ejemplo `N8N_API_KEY_DEV` y `N8N_URL_DEV`.

### **Arquitectura de Seguridad de n8n**

n8n maneja las credenciales de forma segura usando:

1. **Almacenamiento Encriptado**: Las credenciales se encriptan en la base de datos de n8n
2. **Referencias por ID**: Los workflows solo contienen IDs de credenciales, no valores reales
3. **Separación por Ambiente**: Cada instancia de n8n tiene sus propias credenciales

```json
// Workflow JSON - Solo contiene referencias
{
  "nodes": [
    {
      "id": "smtp_node",
      "type": "n8n-nodes-base.emailSend",
      "credentials": {
        "smtp": {
          "id": "smtp_prod_001",    // ← Solo el ID
          "name": "SMTP Production"  // ← Nombre descriptivo
        }
      }
    }
  ]
}
```

### **API Endpoints de n8n para Credenciales**

```bash
# Listar credenciales
GET /api/v1/credentials

# Obtener credencial específica  
GET /api/v1/credentials/{id}

# Crear nueva credencial
POST /api/v1/credentials

# Actualizar credencial
PUT /api/v1/credentials/{id}

# Eliminar credencial
DELETE /api/v1/credentials/{id}
```

## 🛡️ Estrategia de Seguridad de n8n-ops

### **Principios de Seguridad**

1. **Nunca almacenar credenciales en archivos**
2. **Solo usar variables de entorno**
3. **Mapeo automático por ambiente**
4. **Validación antes de despliegue**

### **Mapeo de Credenciales por Ambiente**

```bash
# Development
export SMTP_HOST_DEVELOPMENT="smtp.mailtrap.io"
export SMTP_USER_DEVELOPMENT="dev_user"
export SMTP_PASSWORD_DEVELOPMENT="dev_password"

# Staging  
export SMTP_HOST_STAGING="smtp-staging.company.com"
export SMTP_USER_STAGING="staging_user"
export SMTP_PASSWORD_STAGING="staging_password"

# Production
export SMTP_HOST_PRODUCTION="smtp.sendgrid.net" 
export SMTP_USER_PRODUCTION="apikey"
export SMTP_PASSWORD_PRODUCTION="SG.real_api_key"
```

### **Convención de Naming**

```
{SERVICE}_{FIELD}_{ENVIRONMENT}

Ejemplos:
- SMTP_HOST_DEVELOPMENT
- POSTGRES_PASSWORD_PRODUCTION  
- STRIPE_SECRET_KEY_STAGING
- AWS_ACCESS_KEY_ID_PRODUCTION
```

## 🔄 Flujo de Trabajo Completo

### **1. Desarrollo Local**

```bash
# Configurar variables de entorno
export N8N_URL_DEV="http://localhost:5678"
export N8N_API_KEY_DEV="n8n_api_dev_key"
export SMTP_HOST_DEVELOPMENT="smtp.mailtrap.io"
export SMTP_USER_DEVELOPMENT="dev_user"
export SMTP_PASSWORD_DEVELOPMENT="dev_password"

# Crear workflow con credenciales de desarrollo
n8n-ops branch create email-notifications
# Workflow usa credencial "smtp_dev_001" automáticamente
```

### **2. Despliegue a Staging**

```bash
# Variables de staging (configuradas en CI/CD)
export N8N_URL_STAGING="https://n8n-staging.company.com"
export N8N_API_KEY_STAGING="n8n_api_staging_key" 
export SMTP_HOST_STAGING="smtp-staging.company.com"
export SMTP_USER_STAGING="staging_user"
export SMTP_PASSWORD_STAGING="staging_password"

# n8n-ops mapea automáticamente las credenciales
n8n-ops sync --to-n8n --env staging
# Workflow usa credencial "smtp_staging_001" automáticamente
```

### **3. Producción**

```bash
# Variables de producción (en vault seguro)
export N8N_URL_PROD="https://n8n.company.com"
export N8N_API_KEY_PROD="n8n_api_prod_key"
export SMTP_HOST_PRODUCTION="smtp.sendgrid.net"
export SMTP_USER_PRODUCTION="apikey"  
export SMTP_PASSWORD_PRODUCTION="SG.real_sendgrid_key"

# Despliegue seguro a producción
n8n-ops sync --to-n8n --env production
# Workflow usa credencial "smtp_prod_001" automáticamente
```

## 🔧 Comandos de n8n-ops para Credenciales

### **Listar Credenciales**

```bash
n8n-ops credentials list --env development
```

Output:
```
🔐 Credential Mappings - development Environment
=====================================

📋 N8N API Credentials:
  N8N_URL_DEV = http://localhost:5678
  N8N_API_KEY_DEV = n8n_***dev

🔧 Workflow Node Credentials:
  SMTP:
    host = smtp.mailtrap.io (SMTP_HOST_DEVELOPMENT)
    user = dev_user (SMTP_USER_DEVELOPMENT)
    password = *** (SMTP_PASSWORD_DEVELOPMENT)

  PostgreSQL:
    host = localhost (POSTGRES_HOST_DEVELOPMENT)
    database = myapp_dev (POSTGRES_DB_DEVELOPMENT)
    username = dev_user (POSTGRES_USER_DEVELOPMENT)
```

### **Validar Credenciales**

```bash
n8n-ops credentials validate --env production
```

Output:
```
✅ Validating Credentials - production Environment
=========================================

📊 Validation Results:
  Required: 12
  Present:  10
  Missing:  2

❌ Missing Credentials:
  • STRIPE_SECRET_KEY_PRODUCTION
  • AWS_SECRET_ACCESS_KEY_PRODUCTION

💡 Set missing credentials:
  export STRIPE_SECRET_KEY_PRODUCTION="your_value_here"
  export AWS_SECRET_ACCESS_KEY_PRODUCTION="your_value_here"
```

### **Generar Template**

```bash
n8n-ops credentials template --env production
```

Output:
```bash
# n8n-ops Environment Variables - PRODUCTION
# Copy these variables to your environment

# N8N API Configuration
export N8N_URL_PROD="https://n8n.company.com"
export N8N_API_KEY_PROD="your_production_api_key_here"

# Workflow Node Credentials
# SMTP
export SMTP_HOST_PRODUCTION="smtp.sendgrid.net"
export SMTP_USER_PRODUCTION="apikey"
export SMTP_PASSWORD_PRODUCTION="SG.your_sendgrid_key"

# Database
export POSTGRES_HOST_PRODUCTION="prod-db.company.com"
export POSTGRES_DB_PRODUCTION="myapp_production"
export POSTGRES_USER_PRODUCTION="prod_user"
export POSTGRES_PASSWORD_PRODUCTION="your_db_password"
```

### **Crear, Actualizar y Eliminar Credenciales**

```bash
# Crear
n8n-ops credentials create --file cred.json --env development

# Actualizar
n8n-ops credentials update <id> --file cred.json --env development

# Eliminar
n8n-ops credentials delete <id> --env development
```

## 🎯 Mapeo Automático de Credenciales

### **Algoritmo de Mapeo**

1. **Analizar Workflow**: Identificar nodos que requieren credenciales
2. **Mapear por Ambiente**: Determinar ID de credencial correcto
3. **Validar Existencia**: Verificar que la credencial existe en n8n
4. **Actualizar Referencias**: Cambiar IDs en el workflow JSON

```go
// Ejemplo de mapeo automático
func MapCredentialsForEnvironment(workflowJSON []byte, targetEnv string) {
    // Desarrollo: smtp_dev_001
    // Staging:    smtp_staging_001  
    // Producción: smtp_prod_001
    
    for cada nodo con credenciales {
        if (ambiente == "production") {
            credentialID = "smtp_prod_001"
        } else if (ambiente == "staging") {
            credentialID = "smtp_staging_001"
        } else {
            credentialID = "smtp_dev_001"  
        }
    }
}
```

## 🔐 Seguridad en CI/CD

### **GitLab CI/CD Variables**

```yaml
# .gitlab-ci.yml
variables:
  # Development (masked variables)
  N8N_URL_DEV: "http://localhost:5678"
  
deploy_staging:
  stage: deploy
  script:
    - n8n-ops credentials validate --env staging
    - n8n-ops sync --to-n8n --env staging
  variables:
    # Staging credentials (protected variables)
    N8N_URL_STAGING: "https://n8n-staging.company.com"
  only:
    - staging

deploy_production:
  stage: deploy  
  script:
    - n8n-ops credentials validate --env production
    - n8n-ops sync --to-n8n --env production
  variables:
    # Production credentials (protected + masked)
    N8N_URL_PROD: "https://n8n.company.com"
  only:
    - production
  when: manual
```

### **Vault Integration (Avanzado)**

```bash
# Integración con HashiCorp Vault
n8n-ops credentials sync --vault-path secret/n8n/production
# Automáticamente obtiene credenciales del vault y las configura
```

## 🚨 Mejores Prácticas de Seguridad

### **✅ Hacer**

1. **Usar variables de entorno** para todas las credenciales
2. **Rotar credenciales** regularmente
3. **Validar antes** de cada despliegue
4. **Usar diferentes credenciales** por ambiente
5. **Logs sin credenciales** (enmascarar valores)

### **❌ No Hacer**

1. **Nunca** commit credenciales en Git
2. **Nunca** usar mismas credenciales en dev/prod
3. **Nunca** loggear valores de credenciales
4. **Nunca** compartir credenciales por email/chat
5. **Nunca** hardcodear credenciales en workflows

## 🔄 Migración de Credenciales

### **Script de Migración**

```bash
# Migrar credenciales de desarrollo a staging
n8n-ops credentials migrate \
  --from development \
  --to staging \
  --dry-run

# Ejecutar migración real
n8n-ops credentials migrate \
  --from development \
  --to staging \
  --confirm
```

Este sistema garantiza que las credenciales se manejen de forma segura, usando solo variables de entorno y mapeo automático por ambientes, sin nunca exponer valores sensibles en archivos o logs.