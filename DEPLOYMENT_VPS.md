# 🚀 Despliegue de n8n-ops en VPS - Guía Completa

## 💡 ¿Por qué usar n8n-ops en un VPS?

### Beneficios Principales

**1. Backup Automático de Workflows**
- Todos los workflows se guardan como archivos JSON en el VPS
- Historial completo de cambios en Git
- Backup geográfico separado de tu instancia n8n principal
- Recuperación fácil ante pérdida de datos

**2. Monitoreo Centralizado**
- Web UI accesible desde cualquier lugar
- Monitoreo 24/7 de cambios en workflows
- Alertas automáticas de workflows no commiteados
- Dashboard central para múltiples ambientes

**3. Colaboración Empresarial**
- Equipo completo puede acceder al dashboard
- Control de versiones centralizado
- Flujo de trabajo DevOps profesional
- Integración con GitLab/GitHub

**4. Automatización Robusta**
- Sync programado entre n8n y Git
- Validación automática de credenciales
- Detección proactiva de problemas
- CI/CD integrado para workflows

## 🏗️ Arquitectura de Despliegue VPS

```
┌─────────────────────────────────────────────────────────┐
│                        VPS Server                       │
├─────────────────────────────────────────────────────────┤
│  🌐 Web UI (Port 5000)                                 │
│  ├── Dashboard de monitoreo                            │
│  ├── Gestión de workflows                              │
│  ├── Estado de credenciales                            │
│  └── Operaciones Git                                   │
├─────────────────────────────────────────────────────────┤
│  📁 Filesystem Backup                                  │
│  ├── /workflows/development/                           │
│  ├── /workflows/staging/                               │
│  ├── /workflows/production/                            │
│  └── Git repository con historial                      │
├─────────────────────────────────────────────────────────┤
│  🔄 n8n-ops CLI                                        │
│  ├── Sync automático programado                        │
│  ├── Validación de credenciales                        │
│  ├── Detección de cambios                              │
│  └── Operaciones de branch                             │
├─────────────────────────────────────────────────────────┤
│  🔗 Conexiones Externas                                │
│  ├── n8n Development (puerto 5678)                     │
│  ├── n8n Staging (https://n8n-staging.empresa.com)     │
│  ├── n8n Production (https://n8n.empresa.com)          │
│  └── Git remoto (GitLab/GitHub)                        │
└─────────────────────────────────────────────────────────┘
```

## 🛠️ Instalación en VPS

### Paso 1: Preparar el VPS

```bash
# Actualizar sistema
sudo apt update && sudo apt upgrade -y

# Instalar dependencias
sudo apt install -y git curl build-essential

# Instalar Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### Paso 2: Clonar y Compilar n8n-ops

```bash
# Clonar repositorio
git clone https://gitlab.com/tu-empresa/n8n-ops.git
cd n8n-ops

# Compilar para producción
go build -o n8n-ops-server main.go

# Hacer ejecutable
chmod +x n8n-ops-server
```

### Paso 3: Configurar Variables de Ambiente

```bash
# Crear archivo de configuración
sudo nano /etc/environment

# Agregar variables de n8n
N8N_URL_DEVELOPMENT=http://localhost:5678
N8N_API_KEY_DEVELOPMENT=tu_api_key_dev

N8N_URL_STAGING=https://n8n-staging.empresa.com
N8N_API_KEY_STAGING=tu_api_key_staging

N8N_URL_PRODUCTION=https://n8n.empresa.com
N8N_API_KEY_PRODUCTION=tu_api_key_prod

# Variables de credenciales por ambiente
SMTP_HOST_DEVELOPMENT=smtp.mailtrap.io
SMTP_USER_DEVELOPMENT=user_dev
SMTP_PASSWORD_DEVELOPMENT=pass_dev

SMTP_HOST_PRODUCTION=smtp.sendgrid.net
SMTP_USER_PRODUCTION=apikey
SMTP_PASSWORD_PRODUCTION=SG.real_key

# Variables Git
GITLAB_TOKEN=tu_gitlab_token
GIT_AUTHOR_NAME="n8n-ops Backup"
GIT_AUTHOR_EMAIL="ops@empresa.com"
```

### Paso 4: Configurar Servicio Systemd

```bash
# Crear archivo de servicio
sudo nano /etc/systemd/system/n8n-ops.service
```

```ini
[Unit]
Description=n8n-ops Workflow Management Service
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu/n8n-ops
ExecStart=/home/ubuntu/n8n-ops/n8n-ops-server ui --port 5000
Restart=always
RestartSec=10
Environment=PATH=/usr/local/go/bin:/usr/bin:/bin
EnvironmentFile=/etc/environment

[Install]
WantedBy=multi-user.target
```

```bash
# Habilitar y iniciar servicio
sudo systemctl enable n8n-ops
sudo systemctl start n8n-ops
sudo systemctl status n8n-ops
```

### Paso 5: Configurar Nginx Reverse Proxy

```bash
# Instalar Nginx
sudo apt install -y nginx

# Configurar virtual host
sudo nano /etc/nginx/sites-available/n8n-ops
```

```nginx
server {
    listen 80;
    server_name n8n-ops.empresa.com;

    location / {
        proxy_pass http://localhost:5000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

```bash
# Habilitar sitio
sudo ln -s /etc/nginx/sites-available/n8n-ops /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### Paso 6: SSL con Let's Encrypt

```bash
# Instalar Certbot
sudo apt install -y certbot python3-certbot-nginx

# Obtener certificado SSL
sudo certbot --nginx -d n8n-ops.empresa.com
```

## 🔄 Configurar Sync Automático

### Cron Job para Backup Continuo

```bash
# Editar crontab
crontab -e

# Agregar jobs de sync
# Sync cada 15 minutos desde development
*/15 * * * * cd /home/ubuntu/n8n-ops && ./n8n-ops-server sync --env development --to-filesystem

# Sync diario desde staging y production  
0 2 * * * cd /home/ubuntu/n8n-ops && ./n8n-ops-server sync --env staging --to-filesystem
0 3 * * * cd /home/ubuntu/n8n-ops && ./n8n-ops-server sync --env production --to-filesystem

# Commit automático diario
0 4 * * * cd /home/ubuntu/n8n-ops && git add . && git commit -m "Daily workflow backup $(date)" && git push origin main

# Validación de credenciales cada hora
0 * * * * cd /home/ubuntu/n8n-ops && ./n8n-ops-server credentials validate --env production
```

### Script de Backup Inteligente

```bash
# Crear script de backup
nano /home/ubuntu/n8n-ops/scripts/smart-backup.sh
```

```bash
#!/bin/bash
set -e

BACKUP_DIR="/home/ubuntu/n8n-ops"
LOG_FILE="/var/log/n8n-ops-backup.log"

echo "$(date): Starting smart backup" >> $LOG_FILE

cd $BACKUP_DIR

# Verificar cambios no commiteados
if ./n8n-ops-server status --check-uncommitted --json | jq -r '.workflows.uncommitted' | grep -q '0'; then
    echo "$(date): No uncommitted changes, proceeding with sync" >> $LOG_FILE
else
    echo "$(date): Uncommitted changes detected, committing first" >> $LOG_FILE
    git add .
    git commit -m "Auto-backup: $(date)"
fi

# Sync desde todos los ambientes
for env in development staging production; do
    echo "$(date): Syncing from $env" >> $LOG_FILE
    ./n8n-ops-server sync --env $env --to-filesystem || echo "$(date): Sync failed for $env" >> $LOG_FILE
done

# Push a Git remoto
git push origin main

echo "$(date): Backup completed successfully" >> $LOG_FILE
```

## 📊 Monitoreo y Alertas

### Script de Health Check

```bash
nano /home/ubuntu/n8n-ops/scripts/health-check.sh
```

```bash
#!/bin/bash

# Verificar estado del servicio
if ! systemctl is-active --quiet n8n-ops; then
    echo "ALERT: n8n-ops service is down" | mail -s "n8n-ops Alert" admin@empresa.com
fi

# Verificar conexión a n8n instances
for env in development staging production; do
    if ! ./n8n-ops-server status --env $env --json | jq -r '.health' | grep -q 'ok'; then
        echo "ALERT: n8n $env is unreachable" | mail -s "n8n-ops $env Alert" admin@empresa.com
    fi
done

# Verificar workflows no commiteados por más de 1 hora
uncommitted=$(./n8n-ops-server status --check-uncommitted --json | jq -r '.workflows.uncommitted')
if [ "$uncommitted" -gt 0 ]; then
    echo "WARNING: $uncommitted workflows uncommitted for over 1 hour" | mail -s "n8n-ops Uncommitted Alert" admin@empresa.com
fi
```

## 💰 Casos de Uso Empresariales

### 1. Empresa SaaS con Multiple Clientes
```
VPS Central -> Backup de workflows de todos los clientes
├── Cliente A: n8n-cliente-a.com
├── Cliente B: n8n-cliente-b.com  
└── Cliente C: n8n-cliente-c.com
```

### 2. Agencia de Marketing Digital
```
VPS Backup -> Workflows de campañas de clientes
├── Cliente 1: Email marketing workflows
├── Cliente 2: Social media automation
└── Cliente 3: Lead generation processes
```

### 3. Empresa de E-commerce
```
VPS Production -> Backup crítico de procesos de negocio
├── Order processing workflows
├── Inventory management
├── Customer service automation
└── Payment processing flows
```

## 🎯 Beneficios Específicos del VPS

### Backup y Recuperación
- **Backup geográfico**: VPS en ubicación diferente a n8n principal
- **Versionado completo**: Git mantiene historial de todos los cambios
- **Recuperación rápida**: Restore workflows en minutos
- **Backup programado**: Automático sin intervención manual

### Colaboración
- **Acceso centralizado**: Equipo accede desde cualquier lugar
- **Control de versiones**: Workflow review antes de producción
- **Audit trail**: Quién cambió qué y cuándo
- **Branching strategy**: Feature branches para desarrollos

### Monitoreo
- **Dashboard 24/7**: Vista de estado de todos los ambientes
- **Alertas proactivas**: Email/Slack cuando hay problemas
- **Métricas históricas**: Trends de cambios y deployments
- **Health checks**: Validación continua de conectividad

### Escalabilidad
- **Multi-tenant**: Múltiples clientes en un VPS
- **Load balancing**: Múltiples VPS si es necesario
- **Resource scaling**: CPU/RAM según necesidades
- **Geographic distribution**: VPS en múltiples regiones

## 🚀 Implementación Recomendada

### Para Empresas Pequeñas
- **VPS**: $20-50/mes (2-4 GB RAM)
- **Features**: Web UI + Backup automático + Alertas básicas
- **Ambientes**: Development + Production

### Para Empresas Medianas  
- **VPS**: $50-100/mes (4-8 GB RAM)
- **Features**: Todo lo anterior + CI/CD + Multi-ambiente
- **Ambientes**: Development + Staging + Production

### Para Empresas Grandes
- **VPS**: $100+/mes (8+ GB RAM)
- **Features**: Todo lo anterior + Multi-tenant + HA
- **Ambientes**: Múltiples proyectos/clientes

El ROI es excelente considerando que previene pérdida de workflows críticos y mejora la productividad del equipo.