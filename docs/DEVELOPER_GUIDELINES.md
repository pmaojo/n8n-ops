# n8n-ops Developer Guidelines

## Filosofía del Sistema

n8n-ops está diseñado para facilitar el desarrollo colaborativo de workflows n8n usando metodologías GitOps. El sistema permite que equipos trabajen en workflows de manera organizada, con control de versiones, testing automatizado, y deployment seguro.

## Arquitectura de Proyectos

### Modelo de Repositorio = Proyecto

**Concepto Clave**: Cada repositorio Git representa un **proyecto** que puede contener uno o múltiples workflows relacionados.

```
mi-proyecto-ecommerce/
├── workflows/
│   ├── development/
│   │   ├── Customer_Onboarding_1001.json
│   │   ├── Payment_Processing_1002.json
│   │   └── Order_Fulfillment_1003.json
│   ├── staging/
│   └── production/
├── config/
├── docs/
└── README.md
```

### ¿Cuándo Crear un Nuevo Repositorio?

**SÍ crear nuevo repositorio cuando:**
- Los workflows pertenecen a diferentes dominios de negocio
- Diferentes equipos son responsables del mantenimiento
- Ciclos de release independientes
- Diferentes niveles de seguridad/acceso

**NO crear nuevo repositorio cuando:**
- Los workflows comparten datos o se ejecutan en secuencia
- Mismo equipo de desarrollo
- Misma cadencia de deployment
- Workflows interdependientes

### Ejemplo de Organización

```
Repositorio: "ecommerce-workflows"
├── Customer journey workflows
├── Payment processing workflows  
├── Inventory management workflows

Repositorio: "data-analytics-workflows"  
├── ETL workflows
├── Reporting workflows
├── Data validation workflows

Repositorio: "notifications-workflows"
├── Email campaigns
├── SMS notifications
├── Push notifications
```

## Estrategia de Branching

### Branching Flexible (Recomendado)

**El sistema NO fuerza correspondencia 1:1 entre branches y ambientes**. Esta flexibilidad permite diferentes estrategias según las necesidades del equipo.

#### Opción 1: Git Flow Tradicional
```
main → production
├── staging → staging environment  
├── develop → development
└── feature/payment-fix → development (para testing)
```

#### Opción 2: Trunk-based Development
```
main → todos los ambientes
├── feature/new-workflow → development (para testing)
└── hotfix/critical-bug → staging → production
```

#### Opción 3: Environment Branches
```
production → production environment
staging → staging environment  
develop → development environment
feature/* → development environment
```

### Configuración de Branching

```yaml
# config/branching.yaml
branching:
  strategy: "flexible"
  mappings:
    main: ["production"]
    staging: ["staging"] 
    develop: ["development"]
    "feature/*": ["development"]
    "hotfix/*": ["staging", "production"]
```

## Mejores Prácticas para Workflows

### 1. Creación de Workflows

**✅ RECOMENDADO: Crear workflows desde el CLI**

```bash
# Inicializar proyecto
n8n-ops init --project ecommerce-workflows

# Crear workflow desde template
n8n-ops create --template payment-processing --name "Stripe_Payment_Flow"

# Sincronizar desde n8n existente
n8n-ops sync --from-n8n --env development
```

**Ventajas del CLI:**
- Estructura consistente de archivos
- Metadatos automáticos (git commit, timestamp, environment)
- Validación automática de JSON
- Backup automático antes de cambios
- Integración con GitLab CI/CD

### 2. Estructura de Archivos

```
workflows/development/
├── Customer_Onboarding_1001.json      # Workflow principal
├── Payment_Processing_1002.json       # Workflow principal  
├── _metadata.json                     # Metadatos del sync
├── _backup/                          # Backups automáticos
│   ├── 1001_20240115-142300_backup.json
│   └── 1002_20240115-142300_backup.json
└── README.md                         # Documentación específica
```

### 3. Convenciones de Naming

```
[BusinessDomain]_[Function]_[ID].json

Ejemplos:
- Customer_Onboarding_1001.json
- Payment_Processing_1002.json  
- Inventory_StockCheck_2001.json
- Marketing_EmailCampaign_3001.json
```

### 4. Gestión de Credenciales

```bash
# Configurar credenciales por ambiente
n8n-ops credentials set --env development --key STRIPE_KEY --value sk_test_...
n8n-ops credentials set --env production --key STRIPE_KEY --value sk_live_...

# Validar credenciales
n8n-ops credentials validate --env production
```

## Workflow de Desarrollo

### 1. Desarrollo de Nueva Funcionalidad

```bash
# 1. Crear branch de feature
git checkout -b feature/new-payment-gateway

# 2. Crear o modificar workflows
n8n-ops sync --from-n8n --env development

# 3. Desarrollar en n8n UI
# (hacer cambios en la interfaz de n8n)

# 4. Sincronizar cambios al repositorio
n8n-ops sync --to-git --env development

# 5. Validar workflows
n8n-ops validate ./workflows/development/

# 6. Commit y push
git add .
git commit -m "feat: add PayPal payment gateway workflow"
git push origin feature/new-payment-gateway
```

### 2. Deployment Pipeline

```yaml
# .gitlab-ci.yml
stages:
  - validate
  - test
  - deploy-dev
  - deploy-staging
  - deploy-production

validate:
  script:
    - n8n-ops validate ./workflows/

deploy-development:
  script:
    - n8n-ops deploy --env development --auto-backup
  only:
    - develop
    - /^feature\/.*/

deploy-staging:
  script:
    - n8n-ops deploy --env staging --auto-backup
  when: manual
  only:
    - staging

deploy-production:
  script:
    - n8n-ops deploy --env production --auto-backup
  when: manual
  only:
    - main
```

## Monitoreo y Observabilidad

### Configuración Automática

```bash
# Configurar observabilidad
n8n-ops observability setup --sentry-dsn $SENTRY_DSN --grafana-token $GRAFANA_TOKEN

# Activar monitoreo automático
n8n-ops monitor --env production --failure-threshold 3 --create-issues
```

### Dashboard de Monitoreo

El sistema incluye dashboard en tiempo real:
- URL: http://localhost:5000/dashboard.html
- Monitoreo de todos los ambientes
- Alertas automáticas via GitLab issues
- Métricas de performance y errores

## Resolución de Conflictos

### Conflictos de Workflow

```bash
# Detectar conflictos
n8n-ops status --check-conflicts

# Resolver manualmente
n8n-ops sync --resolve-conflicts --strategy merge

# Rollback si es necesario  
n8n-ops rollback --env staging --to-commit abc123
```

### Estrategias de Merge

1. **Merge Manual**: Para cambios críticos
2. **Auto-merge**: Para cambios compatibles
3. **Overwrite**: Forzar una versión específica

## Comandos Esenciales

```bash
# Inicialización
n8n-ops init --project mi-proyecto

# Sincronización bidireccional
n8n-ops sync --env development        # n8n → Git
n8n-ops sync --to-n8n --env staging   # Git → n8n

# Validación y testing
n8n-ops validate ./workflows/
n8n-ops test --workflow Payment_Processing_1002

# Deployment
n8n-ops deploy --env production --confirm

# Monitoreo
n8n-ops monitor --env production
# El comando monitor se ejecuta en modo continuo, por lo que no
# requiere la bandera `--daemon`.

# Troubleshooting
n8n-ops status --detailed
n8n-ops rollback --env staging --interactive
```

## Seguridad y Compliance

### 1. Gestión de Secretos

- Nunca commitear API keys en workflows
- Usar variables de entorno por ambiente
- Rotar credenciales regularmente

### 2. Control de Acceso

```yaml
# config/access-control.yaml
environments:
  development:
    access: ["developers", "qa"]
  staging: 
    access: ["developers", "qa", "product"]
  production:
    access: ["ops", "senior-developers"]
    requires_approval: true
```

### 3. Auditoría

El sistema automáticamente registra:
- Todos los deployments
- Cambios de workflow
- Acceso a credenciales
- Fallas y recuperaciones

## Troubleshooting Común

### 1. Workflow No Se Sincroniza

```bash
# Verificar conectividad
n8n-ops status --check-connection

# Revisar permisos
n8n-ops credentials validate

# Forzar resync
n8n-ops sync --force --env development
```

### 2. Conflictos de Merge

```bash
# Ver diferencias
n8n-ops diff --env staging

# Resolver conflicts
n8n-ops merge --interactive --env staging
```

### 3. Performance Issues

```bash
# Ver métricas
n8n-ops observability dashboard

# Analizar logs
n8n-ops logs --workflow Payment_Processing --last 24h
```

## Conclusión

Este sistema está diseñado para ser flexible y adaptarse a diferentes metodologías de desarrollo. La clave está en:

1. **Organizar workflows por dominio de negocio** en repositorios
2. **Usar branching strategy** que se adapte al equipo
3. **Aprovechar el CLI** para todas las operaciones
4. **Monitoreo continuo** en producción
5. **Documentar decisiones** arquitecturales

El objetivo es mantener los workflows como código, con todas las ventajas del desarrollo de software moderno: control de versiones, testing, CI/CD, y colaboración eficiente.