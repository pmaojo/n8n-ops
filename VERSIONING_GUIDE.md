# Guía de Versionado y Flujo de Trabajo n8n-ops

## 🏷️ Convención de Versionado

### Semantic Versioning (SemVer)
Usamos **Semantic Versioning** para mantener consistencia:

```
MAJOR.MINOR.PATCH-ENVIRONMENT
```

**Ejemplos:**
- `v1.2.3` - Versión estable de producción
- `v1.2.4-dev` - Versión en desarrollo
- `v1.3.0-rc1` - Release candidate
- `v2.0.0` - Breaking changes

### Estructura de Versiones

#### 1. **MAJOR** (cambios incompatibles)
- Cambios que rompen workflows existentes
- Nuevas dependencias obligatorias
- Cambios en estructura de datos

#### 2. **MINOR** (nuevas funcionalidades)
- Nuevos nodos o funcionalidades
- Mejoras compatibles hacia atrás
- Optimizaciones de performance

#### 3. **PATCH** (correcciones)
- Corrección de bugs
- Pequeñas mejoras
- Actualizaciones de seguridad

## 🔄 Flujo de Trabajo (Git Flow)

### Ramas Principales

```
main (production)     ─────●─────●─────●─────●
                           │     │     │     │
staging                    ●─────●─────●─────●
                           │     │     │     │
develop                    ●─────●─────●─────●
                           │     │     │     │
feature/nueva-funcionalidad ●─────●
```

### 1. **Branch: `main`** (Producción)
- Solo código estable y probado
- Deploy automático a producción
- Protegida: requiere pull request + revisión
- Versiones: `v1.0.0`, `v1.1.0`, `v2.0.0`

### 2. **Branch: `staging`** (Pruebas)
- Testing antes de producción
- Deploy automático a staging
- Versiones: `v1.1.0-rc1`, `v1.1.0-rc2`

### 3. **Branch: `develop`** (Desarrollo)
- Integración de nuevas funcionalidades
- Deploy automático a development
- Versiones: `v1.1.0-dev`, `v1.1.0-alpha`

### 4. **Feature Branches**
- `feature/customer-onboarding`
- `feature/payment-integration`
- `hotfix/urgent-bug-fix`

## 📋 Proceso de Deployment

### Flujo Completo

```
1. Desarrollo Local
   ├── Crear feature branch: git checkout -b feature/nueva-funcionalidad
   ├── Desarrollar workflows
   ├── Validar local: ./n8n-ops validate ./workflows/
   └── Commit cambios

2. Testing en Development
   ├── Push a develop: git push origin develop
   ├── GitLab CI/CD despliega automáticamente
   ├── Verificar: ./n8n-ops status --env development
   └── Probar funcionalidad

3. Staging (Pre-producción)
   ├── Merge a staging: git checkout staging && git merge develop
   ├── Deploy manual en GitLab
   ├── Testing completo: ./n8n-ops status --env staging
   └── Validación stakeholders

4. Producción
   ├── Merge a main: git checkout main && git merge staging
   ├── Tag versión: git tag v1.2.0 && git push --tags
   ├── Deploy manual protegido
   └── Verificar: ./n8n-ops status --env production
```

## 🎯 Comandos por Ambiente

### Development (Automático)
```bash
# Verificar estado
./n8n-ops status --env development

# Sincronizar desde servidor
./n8n-ops sync --env development

# Deploy automático via GitLab CI
git push origin develop
```

### Staging (Manual)
```bash
# Verificar antes de deploy
./n8n-ops status --env staging

# Deploy manual en GitLab CI/CD
# O localmente:
./n8n-ops deploy --env staging --dry-run
./n8n-ops deploy --env staging
```

### Production (Manual + Protegido)
```bash
# Verificar estado actual
./n8n-ops status --env production

# Deploy protegido (requiere aprobación)
# Solo via GitLab CI/CD con revisión obligatoria

# Rollback en emergencia
./n8n-ops rollback --env production
```

## 📊 Tracking de Versiones

### Ver Versiones Actuales
```bash
# Estado de todos los workflows
./n8n-ops status --env production --json

# Comparar ambientes
./n8n-ops status --env staging
./n8n-ops status --env production
```

### Ejemplo de Output
```
Workflow Status - PRODUCTION Environment

NAME                    VERSION    LAST DEPLOY   GIT COMMIT   NODES
────                    ───────    ───────────   ──────────   ─────
Customer Onboarding     v1.2.3     2025-07-20    a1b2c3d      8
Payment Processing      v1.1.8     2025-07-19    e4f5g6h      12
Email Notifications     v1.0.5     2025-07-18    i7j8k9l      5
```

## 🔍 Convenciones de Nomenclatura

### Workflows
```
kebab-case con propósito claro:
- customer-onboarding
- payment-processing  
- email-notifications
- data-backup-daily
```

### Git Tags
```
v1.0.0      - Versión de producción
v1.1.0-rc1  - Release candidate
v1.1.0-dev  - Versión desarrollo
v1.0.1      - Hotfix
```

### Commits
```
feat: add customer onboarding workflow
fix: resolve payment timeout issue  
docs: update deployment guide
refactor: optimize email notification logic
```

## 🚀 Mejores Prácticas

### 1. **Antes de cada Deploy**
```bash
# Validar workflows
./n8n-ops validate ./workflows/production/

# Verificar estado actual
./n8n-ops status --env production

# Dry-run del deploy
./n8n-ops deploy --env production --dry-run
```

### 2. **Después de cada Deploy**
```bash
# Verificar que se desplegó correctamente
./n8n-ops status --env production

# Probar workflows críticos
# Monitorear logs por 15-30 minutos
```

### 3. **En caso de Problemas**
```bash
# Rollback inmediato
./n8n-ops rollback --env production

# Verificar rollback
./n8n-ops status --env production
```

Este sistema te permite tener control total sobre qué versión de cada workflow está en cada ambiente, con trazabilidad completa desde desarrollo hasta producción.