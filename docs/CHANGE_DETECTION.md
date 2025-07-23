# n8n-ops Change Detection System

## 🎯 Problema Resuelto

**Pregunta**: "¿Cómo sabe n8n-ops que he editado un flujo en mi n8n local en Docker?"

**Respuesta**: n8n-ops implementa un sistema inteligente de **detección automática de cambios** que monitorea tanto tu instancia de n8n como los archivos JSON locales, comparando timestamps y contenido para determinar qué ha cambiado y en qué dirección sincronizar.

## 🔄 Estrategia de Sincronización Bidireccional

### 1. **Detección Automática** 
n8n-ops puede detectar cambios de 3 formas:

```bash
# Modo Watch - Monitoreo continuo (RECOMENDADO)
n8n-ops watch --env development --interval 30s

# Sync Manual - Verificación bajo demanda  
n8n-ops sync --env development

# Check Status - Solo verificar sin sincronizar
n8n-ops check --env development
```

### 2. **Direcciones de Sincronización**

#### 🔽 FROM n8n TO Git
```bash
# Cuando editas workflows en la UI de n8n:
n8n-ops sync --from-n8n --env development
```
- Detecta workflows modificados en n8n UI
- Descarga JSON actualizado a archivos locales
- Actualiza timestamps para seguimiento

#### 🔼 FROM Git TO n8n  
```bash
# Cuando editas archivos JSON localmente:
n8n-ops sync --to-n8n --env development
```
- Detecta archivos JSON modificados localmente
- Sube cambios a la instancia de n8n vía API
- Actualiza workflows en n8n UI

#### ⚡ BIDIRECCIONAL (Automático)
```bash
# Modo inteligente (por defecto):
n8n-ops sync --env development
```
- Compara timestamps de modificación
- Sincroniza en ambas direcciones según necesidad
- Detecta conflictos y pide resolución

## 📊 Lógica de Detección de Cambios

### Algoritmo de Comparación

```
1. OBTENER workflows desde n8n API
2. OBTENER archivos JSON locales 
3. PARA cada workflow:
   - Comparar hash MD5 del contenido
   - Si difieren → Comparar timestamps:
     * Remote > Local + 30s → FROM n8n TO Git
     * Local > Remote + 30s → FROM Git TO n8n  
     * Diferencia < 30s → CONFLICTO (pedir resolución)
4. DETECTAR workflows nuevos/eliminados
5. GENERAR reporte de cambios
```

### Ejemplo Práctico

**Escenario**: Editas un workflow en n8n UI

```bash
# 1. n8n-ops detecta el cambio
$ n8n-ops check --env development
🔍 Workflow Sync Status - development Environment
📝 Changes detected:
   • "Customer Onboarding" modified in n8n (remote newer)
   • Last modified: 2 minutes ago
   
⚡ Action required: Run sync to download changes

# 2. Sincronizar cambios automáticamente  
$ n8n-ops sync --env development
🔄 Syncing workflows from development environment...
📝 Workflow "Customer Onboarding" updated from n8n
📁 Saved to: workflows/development/customer-onboarding.json
✅ 1 workflow synchronized
```

## 🕵️ Modo Watch - Monitoreo Continuo

### Configuración de Monitoring
```bash
# Monitoreo básico cada 10 segundos
n8n-ops watch --env development

# Monitoreo con auto-commit a Git
n8n-ops watch --env development --auto-commit --interval 30s

# Monitoreo con sync automático
n8n-ops watch --env development --auto-sync --interval 1m
```

### Output del Watch Mode
```
👁️ Watching n8n workflows - development environment
🔄 Check interval: 30s
✅ Connected to n8n API. Monitoring for changes...

[14:32:15] 🆕 New workflow detected: Email Automation (wf_abc123)
[14:32:15] 🔄 Auto-syncing changes...
[14:32:16] ✅ Changes synced successfully
[14:32:16] 📝 Changes committed to Git

[14:35:20] 📝 Workflow updated: Customer Onboarding (wf_def456)  
[14:35:20] 🔄 Auto-syncing changes...
[14:35:21] ✅ Changes synced successfully
```

## 🎛️ Configuración de Docker para n8n Local

### docker-compose.yml Recomendado
```yaml
version: '3.8'
services:
  n8n:
    image: n8nio/n8n
    ports:
      - "5678:5678"
    environment:
      - N8N_BASIC_AUTH_ACTIVE=false
      - N8N_API_KEY_ENABLE=true  # ← IMPORTANTE para n8n-ops
    volumes:
      - n8n_data:/home/node/.n8n
    restart: unless-stopped

volumes:
  n8n_data:
```

### Variables de Entorno para n8n-ops
```bash
# Configurar conexión a tu n8n local
export N8N_URL="http://localhost:5678"
export N8N_API_KEY="n8n_api_xxxxxxxxxxxxxxxx"

# Ejecutar n8n-ops
n8n-ops sync --env development
```

## 🔧 Resolución de Conflictos

### Conflictos Automáticos vs Manuales

#### Conflicto Detectado
```bash
$ n8n-ops sync --env development
⚠️ Conflict detected for workflow "Data Processing":
   • Local file:  Modified 14:30 (2 min ago)
   • n8n remote: Modified 14:32 (30 sec ago)
   
❓ Resolution required:
   [1] Use remote version (from n8n)
   [2] Use local version (upload to n8n)  
   [3] Show diff and merge manually
   [4] Skip this workflow
   
Choose option [1-4]: 
```

#### Resolución Forzada
```bash  
# Siempre usar versión remota (de n8n)
n8n-ops sync --from-n8n --force --env development

# Siempre usar versión local (subir a n8n)
n8n-ops sync --to-n8n --force --env development
```

## 📁 Estructura de Archivos Resultante

```
project/
├── workflows/
│   ├── development/
│   │   ├── customer-onboarding.json
│   │   ├── email-automation.json
│   │   └── data-processing.json
│   ├── staging/
│   │   ├── customer-onboarding.json
│   │   └── email-automation.json
│   └── production/
│       └── customer-onboarding.json
├── .n8n-ops/
│   ├── sync-state.json          # Estado de sincronización
│   └── change-reports/          # Reportes de cambios para CI/CD
│       ├── 2025-01-15-sync.json
│       └── 2025-01-16-sync.json
└── .git/
    └── # Control de versiones normal
```

## 🚀 Workflow Completo de Desarrollo

### 1. Setup Inicial
```bash
# Inicializar proyecto
n8n-ops init --name mi-proyecto

# Configurar n8n local  
export N8N_URL="http://localhost:5678"
export N8N_API_KEY="tu-api-key"

# Primera sincronización
n8n-ops sync --env development
```

### 2. Desarrollo Colaborativo
```bash
# Desarrollador A - trabaja en n8n UI
n8n-ops watch --env development --auto-commit

# Desarrollador B - edita archivos JSON
git pull
# editar workflows/development/mi-workflow.json
n8n-ops sync --to-n8n --env development
git add . && git commit -m "Update workflow logic"
```

### 3. Despliegue a Staging/Production
```bash
# Git maneja el despliegue (como sugieres)
git checkout staging
git merge development
n8n-ops sync --to-n8n --env staging

# Production via Git Flow
git checkout production  
git merge staging
n8n-ops sync --to-n8n --env production
```

## 💡 Beneficios del Enfoque

1. **Sin rollback/deploy complejos**: Git maneja versioning
2. **Detección inteligente**: Sabe automáticamente qué cambió
3. **Sincronización bidireccional**: Funciona editando en UI o archivos
4. **Colaboración real**: Múltiples desarrolladores sin conflictos
5. **CI/CD Ready**: Reportes JSON para pipelines automáticos
6. **Zero-downtime**: Solo llamadas API, sin tocar infraestructura

Este sistema permite que **trabajes como quieras** - editando en la UI de n8n o modificando JSONs directamente - mientras n8n-ops se encarga inteligentemente de mantener todo sincronizado.