# Guía de Uso - n8n-ops

## Instalación y Configuración Inicial

### 1. Compilar el CLI
```bash
go build -o n8n-ops main.go
chmod +x n8n-ops
```

### 2. Configurar Variables de Entorno
```bash
# API de n8n
export N8N_API_KEY="tu-api-key-de-n8n"
export N8N_URL="https://tu-n8n-instance.com"

# GitLab (para issues automáticas)
export GITLAB_TOKEN="tu-gitlab-token"
export GITLAB_PROJECT_ID="12345"
export GITLAB_URL="https://gitlab.com"
```

## Comandos Principales

### 📊 Monitoreo Automático (Nuevo)

#### Modo Producción
```bash
# Monitorear workflows en producción
./n8n-ops monitor --env production

# Con configuración personalizada
./n8n-ops monitor \
  --env production \
  --failure-threshold 3 \
  --check-interval 30s \
  --gitlab-project 12345
```

#### Modo Demo (Sin GitLab)
```bash
# Probar el sistema sin configuración real
./n8n-ops monitor --demo --env development
```

**¿Qué hace?**
- Consulta la API de n8n cada minuto (configurable)
- Detecta workflows que fallan consecutivamente
- Crea issues automáticamente en GitLab
- Actualiza issues cuando workflows se recuperan

### 🔄 Sincronización de Workflows

#### Bajar workflows de n8n a archivos JSON
```bash
# Sincronizar desde ambiente de desarrollo
./n8n-ops sync --env development

# Sincronizar desde producción
./n8n-ops sync --env production --backup
```

#### Subir workflows locales a n8n
```bash
# Desplegar cambios a desarrollo
./n8n-ops sync --env development --upload

# Desplegar a producción (con validación)
./n8n-ops sync --env production --upload --validate
```

### 👁️ Daemon Mode (File Watching)

```bash
# Observar cambios en archivos JSON y sincronizar automáticamente
./n8n-ops --daemon --env development

# Con backup automático
./n8n-ops --daemon --env development --backup
```

### ✅ Validación de Workflows

```bash
# Validar todos los workflows
./n8n-ops validate

# Validar ambiente específico
./n8n-ops validate --env staging

# Solo validación de sintaxis
./n8n-ops validate --syntax-only
```

### 📈 Estado del Sistema

```bash
# Ver estado general
./n8n-ops status

# Estado de ambiente específico
./n8n-ops status --env production

# Con detalles de últimas executions
./n8n-ops status --env production --show-executions
```

### 🔧 Gestión de Credenciales

```bash
# Listar credenciales disponibles
./n8n-ops credentials list --env production

# Validar credenciales
./n8n-ops credentials validate --env production

# Rotar credenciales (si implementado)
./n8n-ops credentials rotate --credential-id "db-conn"
```

## Casos de Uso Comunes

### 1. Setup Inicial de Proyecto
```bash
# Inicializar estructura de proyecto
./n8n-ops init --project-name "mi-proyecto"

# Bajar workflows existentes
./n8n-ops sync --env development
./n8n-ops sync --env production

# Configurar monitoreo
./n8n-ops monitor --env production --setup
```

### 2. Desarrollo Diario
```bash
# Modo desarrollo con watching automático
./n8n-ops --daemon --env development &

# Hacer cambios en archivos JSON
# Los cambios se sincronizan automáticamente

# Validar antes de commit
./n8n-ops validate --all
```

### 3. Deploy a Producción
```bash
# Validar workflows
./n8n-ops validate --env production

# Backup actual
./n8n-ops sync --env production --backup-only

# Deploy con verificación
./n8n-ops sync --env production --upload --verify

# Iniciar monitoreo post-deploy
./n8n-ops monitor --env production --duration 10m
```

### 4. Troubleshooting
```bash
# Ver workflows con problemas
./n8n-ops check --env production --show-failed

# Ver executions recientes
./n8n-ops status --env production --show-executions

# Validar conectividad
./n8n-ops status --health-check
```

## Estructura de Directorios

Después del `init`, tendrás:
```
mi-proyecto/
├── workflows/
│   ├── development/
│   ├── staging/
│   └── production/
├── .gitlab/
│   ├── issue_templates/
│   └── merge_request_templates/
├── docs/
├── config.yaml
└── README.md
```

## Configuración Avanzada

### Archivo config.yaml
```yaml
# Configuración global
default_environment: development
backup_enabled: true
validation_strict: true

# Por ambiente
environments:
  development:
    n8n_url: "http://localhost:5678"
    api_key: "${DEV_N8N_API_KEY}"
    create_issues: false
    
  production:
    n8n_url: "https://n8n-prod.company.com"
    api_key: "${PROD_N8N_API_KEY}"
    create_issues: true
    failure_threshold: 2

# GitLab integration
gitlab:
  url: "https://gitlab.company.com"
  project_id: "12345"
  token: "${GITLAB_TOKEN}"
  
# Monitoreo
monitoring:
  check_interval: "1m"
  failure_threshold: 3
  auto_recovery: true
```

## Integración con CI/CD

### GitLab CI Pipeline
```yaml
stages:
  - validate
  - deploy
  - monitor

validate:
  script:
    - n8n-ops validate --all-environments

deploy:production:
  script:
    - n8n-ops sync --env production --upload --verify
  when: manual
  only:
    - main

monitor:
  script:
    - n8n-ops monitor --env production --duration 5m
  after_script:
    - echo "Check GitLab issues for any failures"
```

## Ejemplos Prácticos

### Monitoreo 24/7 en Servidor
```bash
# Crear servicio systemd
sudo tee /etc/systemd/system/n8n-ops-monitor.service << EOF
[Unit]
Description=n8n-ops Workflow Monitor
After=network.target

[Service]
Type=simple
User=n8n-ops
WorkingDirectory=/opt/n8n-ops
ExecStart=/opt/n8n-ops/n8n-ops monitor --env production
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl enable n8n-ops-monitor
sudo systemctl start n8n-ops-monitor
```

### Script de Backup Diario
```bash
#!/bin/bash
# backup-workflows.sh
DATE=$(date +%Y%m%d)
./n8n-ops sync --env production --backup-only
tar -czf "backups/workflows-$DATE.tar.gz" workflows/production/
echo "Backup created: workflows-$DATE.tar.gz"
```

### Notificaciones Slack
```bash
# En script personalizado
if ./n8n-ops status --env production --check-failed; then
    curl -X POST -H 'Content-type: application/json' \
    --data '{"text":"⚠️ Workflows failing in production"}' \
    $SLACK_WEBHOOK_URL
fi
```

## Comandos de Ayuda

```bash
# Ayuda general
./n8n-ops --help

# Ayuda de comando específico
./n8n-ops monitor --help
./n8n-ops sync --help

# Ver ejemplos
./n8n-ops welcome

# Ver versión
./n8n-ops version
```

## Tips y Mejores Prácticas

### 1. Ambientes Separados
- Usa `--env` siempre para evitar accidentes
- Configuraciones diferentes para dev/staging/prod

### 2. Backups
- Siempre haz backup antes de deploy a producción
- Automatiza backups diarios

### 3. Monitoreo
- Usa umbrales diferentes por ambiente
- Monitoring continuo en producción
- Demo mode para testing

### 4. Colaboración
- Usa templates GitLab para issues consistentes
- CI/CD para deploys automáticos
- Code review para cambios críticos

El sistema está diseñado para ser simple pero potente, proporcionando monitoreo automático y gestión completa de workflows n8n.