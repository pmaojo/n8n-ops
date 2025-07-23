# Ejemplos de Flujo de Trabajo n8n-ops

## 📋 Ejemplo Práctico: Nueva Funcionalidad

### Escenario: Agregar workflow de "Notificaciones Push"

#### 1. **Crear Feature Branch**
```bash
# Crear rama desde develop
git checkout develop
git pull origin develop
git checkout -b feature/push-notifications

# Verificar rama actual
git branch
```

#### 2. **Desarrollo Local**
```bash
# Crear workflow en development
# Sincronizar desde servidor desarrollo
./n8n-ops sync --env development

# Ver workflows actuales
./n8n-ops status --env development

# Desarrollar en n8n web interface
# Cuando esté listo, sincronizar de vuelta
./n8n-ops sync --env development
```

#### 3. **Testing en Development**
```bash
# Validar workflow antes de commit
./n8n-ops validate ./workflows/development/

# Commit cambios
git add .
git commit -m "feat: add push notifications workflow"

# Push a develop branch
git push origin feature/push-notifications

# Crear Pull Request a develop
```

#### 4. **Deploy a Staging**
```bash
# Después de merge a develop
git checkout staging
git pull origin staging
git merge develop

# Verificar versión antes de deploy
./n8n-ops status --env staging

# Deploy a staging (manual en GitLab)
git push origin staging

# Verificar deployment
./n8n-ops status --env staging
```

#### 5. **Deploy a Producción**
```bash
# Después de testing exitoso
git checkout main
git pull origin main
git merge staging

# Tag nueva versión
git tag v1.3.0
git push origin main --tags

# Deploy manual protegido en GitLab
# Verificar después del deploy
./n8n-ops status --env production
```

## 🔄 Ejemplo Práctico: Hotfix Urgente

### Escenario: Bug crítico en producción

#### 1. **Crear Hotfix Branch**
```bash
# Crear desde main (producción)
git checkout main
git pull origin main
git checkout -b hotfix/payment-timeout-fix
```

#### 2. **Fix Rápido**
```bash
# Sincronizar workflows actuales
./n8n-ops sync --env production

# Editar workflow problemático
# Validar fix
./n8n-ops validate ./workflows/production/payment-processing.json

# Commit
git add .
git commit -m "fix: resolve payment timeout in production workflow"
```

#### 3. **Deploy Directo**
```bash
# Merge a main
git checkout main
git merge hotfix/payment-timeout-fix

# Tag hotfix
git tag v1.2.1
git push origin main --tags

# Deploy inmediato
# Verificar
./n8n-ops status --env production
```

#### 4. **Sync a Otras Ramas**
```bash
# Aplicar fix a develop también
git checkout develop
git merge hotfix/payment-timeout-fix
git push origin develop

# Aplicar a staging
git checkout staging  
git merge hotfix/payment-timeout-fix
git push origin staging
```

## 📊 Ejemplo: Auditoría de Versiones

### Ver Estado de Todos los Ambientes

```bash
# Script para comparar versiones
#!/bin/bash
echo "=== DEVELOPMENT ==="
./n8n-ops status --env development

echo -e "\n=== STAGING ==="
./n8n-ops status --env staging  

echo -e "\n=== PRODUCTION ==="
./n8n-ops status --env production
```

### Output Esperado
```
=== DEVELOPMENT ===
Workflow Status - development Environment

NAME                    VERSION    LAST MODIFIED     NODES
────                    ───────    ─────────────     ─────
Customer Onboarding     v1.4.0-dev 2025-07-22 15:30  8
Payment Processing      v1.3.5-dev 2025-07-22 14:20  12
Push Notifications      v1.0.0-dev 2025-07-22 16:45  6

=== STAGING ===
Workflow Status - staging Environment

NAME                    VERSION    LAST MODIFIED     NODES
────                    ───────    ─────────────     ─────  
Customer Onboarding     v1.3.2-rc1 2025-07-21 10:15  8
Payment Processing      v1.3.4     2025-07-20 18:30  12

=== PRODUCTION ===  
Workflow Status - production Environment

NAME                    VERSION    LAST MODIFIED     NODES
────                    ───────    ─────────────     ─────
Customer Onboarding     v1.3.1     2025-07-20 09:00  8  
Payment Processing      v1.3.4     2025-07-20 18:35  12
```

### Análisis del Output
- **Development**: Tiene la versión más nueva (v1.4.0-dev) y el nuevo workflow de Push Notifications
- **Staging**: Tiene versión candidata (v1.3.2-rc1) lista para testing  
- **Production**: Tiene versión estable (v1.3.1), una versión detrás de staging

## 🚨 Ejemplo: Rollback de Emergencia

### Escenario: Deploy con problemas en producción

#### 1. **Detectar Problema**
```bash
# Verificar estado actual
./n8n-ops status --env production

# Ver logs de deployment (en GitLab CI/CD)
# Confirmar necesidad de rollback
```

#### 2. **Rollback Inmediato**
```bash
# Rollback vía GitLab CI/CD (botón manual)
# O si tienes acceso local de emergencia:
./n8n-ops rollback --env production

# Verificar rollback exitoso  
./n8n-ops status --env production
```

#### 3. **Comunicación**
```bash
# Notificar al equipo
# Documentar incidente
# Crear issue para fix
```

## 🎯 Tips y Mejores Prácticas

### 1. **Naming Conventions**
```bash
# Branches
feature/customer-loyalty-program
hotfix/email-template-encoding  
release/v1.4.0

# Commits
feat: add customer loyalty points calculation
fix: resolve email template encoding issue
docs: update deployment procedures
refactor: optimize payment processing logic
```

### 2. **Testing Checklist**
```bash
# Antes de cada deploy
□ ./n8n-ops validate ./workflows/[env]/
□ ./n8n-ops status --env [env]  
□ ./n8n-ops deploy --env [env] --dry-run
□ Verificar dependencias externas
□ Backup de workflows críticos
```

### 3. **Monitoring Post-Deploy**
```bash
# Después de deploy
□ ./n8n-ops status --env [env]
□ Verificar workflows críticos manualmente
□ Monitorear logs por 30 minutos  
□ Confirmar métricas de negocio
□ Notificar éxito al equipo
```

Este flujo te da control total y trazabilidad completa de cada workflow en cada ambiente.