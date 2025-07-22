# Detección de Cambios sin Sincronizar

## 🔍 Estrategias de Detección

### 1. **Detección Manual (Actual)**
```bash
# Verificar cambios manualmente
./n8n-ops sync --env development --dry-run
./n8n-ops status --env development

# Comparar con último commit
git status workflows/development/
git diff workflows/development/
```

### 2. **Comparación por Hash/Timestamp**
El CLI puede comparar:
- **Hash del workflow** en n8n vs local
- **Timestamp de modificación** (updatedAt)
- **Número de versión** del workflow

### 3. **Sync Automático Programado**
```bash
# Cron job cada 15 minutos
*/15 * * * * cd /path/to/project && ./n8n-ops sync --env development

# Script con notificaciones
./n8n-ops sync --env development && \
  if git status --porcelain workflows/development/ | grep -q .; then
    echo "⚠️  Workflows modificados sin sincronizar detectados"
    git status workflows/development/
  fi
```

## 🎯 Comando `check` - Detección de Cambios

### Implementación Propuesta
```bash
# Verificar si hay cambios sin sincronizar
./n8n-ops check --env development

# Output esperado:
✅ Workflows sincronizados: 3
⚠️  Workflows modificados: 2
   - customer-onboarding (modified 5 min ago)
   - payment-processing (modified 12 min ago)

💡 Ejecuta './n8n-ops sync --env development' para sincronizar
```

### Casos de Uso
```bash
# Antes de hacer commit
./n8n-ops check --env development
git add . && git commit -m "..."

# Verificación automática en CI/CD
./n8n-ops check --env staging --fail-if-changes
```

## 📊 Comparación Detallada

### Workflow Local vs n8n Instance
```json
// Estado actual en n8n (via API)
{
  "id": "1001",
  "name": "Customer Onboarding", 
  "updatedAt": "2025-07-22T16:30:00Z",
  "versionId": 15,
  "nodes": [...] // 8 nodos
}

// Estado local (archivo JSON)
{
  "id": "1001", 
  "name": "Customer Onboarding",
  "updatedAt": "2025-07-22T15:45:00Z", 
  "versionId": 14,
  "nodes": [...] // 7 nodos
}

// Resultado: Workflow modificado en n8n (45 min ago)
```

## 🔄 Flujo de Detección

### Proceso Completo
```
1. Developer edita workflow en n8n web interface
   ↓
2. n8n guarda cambios (updatedAt, versionId++)
   ↓  
3. ./n8n-ops check --env development
   ↓
4. CLI compara timestamps/versions/hashes
   ↓
5. CLI reporta workflows "out of sync"
   ↓
6. Developer ejecuta sync para actualizar local
   ↓
7. Git muestra cambios para commit
```

## ⚙️ Configuración de Monitoreo

### GitLab CI/CD con Detección
```yaml
# Job que verifica cambios sin sincronizar
check-sync-status:
  stage: validate
  script:
    - ./n8n-ops check --env development --json > sync-status.json
    - |
      if [ $(jq '.modified | length' sync-status.json) -gt 0 ]; then
        echo "⚠️ Workflows modified in n8n but not synced to Git"
        jq '.modified[] | "- \(.name) (modified \(.timeAgo))"' sync-status.json
        echo "💡 Consider running sync before deployment"
      fi
  rules:
    - if: '$CI_COMMIT_BRANCH == "develop"'
```

### Webhook n8n → GitLab (Avanzado)
```bash
# Configurar webhook en n8n que llame GitLab API
# Cuando se guarda workflow → Trigger pipeline sync

# En n8n workflow:
# Trigger: Webhook "On workflow save"
# Action: HTTP Request to GitLab API
# URL: https://gitlab.com/api/v4/projects/123/pipeline
```

## 📈 Estrategias por Ambiente

### Development (Detección Frecuente)
```bash
# Verificar cada 5 minutos
*/5 * * * * ./n8n-ops check --env development --notify-if-changes

# Auto-sync si hay cambios menores
./n8n-ops sync --env development --auto-commit
```

### Staging (Detección Manual)
```bash  
# Verificar antes de deploy
./n8n-ops check --env staging
./n8n-ops sync --env staging --dry-run
```

### Production (Solo Alertas)
```bash
# Solo alertar, nunca auto-sync
./n8n-ops check --env production --alert-only
```

## 🛠️ Implementación del Comando `check`

### Uso Básico
```bash
# Verificar estado de sincronización
./n8n-ops check --env development

# Modo silencioso (solo exit code)
./n8n-ops check --env development --quiet
echo $? # 0=sync, 1=changes detected

# JSON para scripts
./n8n-ops check --env development --json
```

### Output JSON
```json
{
  "environment": "development",
  "lastSync": "2025-07-22T15:45:00Z",
  "totalWorkflows": 5,
  "synchronized": 3,
  "modified": 2,
  "workflows": {
    "synchronized": [
      {"id": "1003", "name": "Email Notifications", "status": "sync"}
    ],
    "modified": [
      {
        "id": "1001", 
        "name": "Customer Onboarding",
        "localVersion": 14,
        "remoteVersion": 15, 
        "lastModified": "2025-07-22T16:30:00Z",
        "timeAgo": "5 minutes ago"
      }
    ]
  }
}
```

## 💡 Mejores Prácticas

### 1. **Verificación Pre-Commit**
```bash
# Git hook pre-commit
#!/bin/bash
./n8n-ops check --env development --fail-if-changes
if [ $? -eq 1 ]; then
  echo "⚠️ Workflows modified in n8n. Run sync first:"
  echo "   ./n8n-ops sync --env development"
  exit 1
fi
```

### 2. **Monitoreo Continuo**
```bash
# Script de monitoreo
while true; do
  ./n8n-ops check --env development --quiet
  if [ $? -eq 1 ]; then
    echo "$(date): Changes detected in development workflows"
    ./n8n-ops check --env development
  fi
  sleep 300 # 5 minutos
done
```

### 3. **Integración con Slack/Teams**
```bash
# Notificar a equipo si hay cambios
./n8n-ops check --env production --json | \
  jq -r 'if .modified > 0 then "⚠️ \(.modified) workflows modified in production without sync" else empty end' | \
  curl -X POST -H 'Content-Type: application/json' \
       -d "{\"text\": \"$(cat)\"}" \
       $SLACK_WEBHOOK_URL
```

El comando `check` te permitirá detectar proactivamente cuando hay workflows modificados en n8n que necesitan sincronización con Git.