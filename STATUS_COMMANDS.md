# Comandos para Ver Versiones de Workflows

## 🔍 Ver Status de Workflows en Producción

### Comando `status` - Estado completo ✅ FUNCIONANDO
```bash
# Ver estado de todos los workflows en producción
./n8n-ops status --env production

# Incluir workflows inactivos también  
./n8n-ops status --env production --all

# Salida en formato JSON
./n8n-ops status --env production --json
```

**Muestra:**
- 📋 Nombre del workflow
- ✅ Estado (activo/inactivo)  
- 📊 Versión (local/desplegada)
- 📅 Última modificación
- 🔢 Cantidad de nodos

### Comando `list` - Lista simple
```bash  
# Listar workflows en producción
./n8n-ops list --env production

# Solo workflows activos
./n8n-ops list --env production --active-only

# Buscar workflows específicos
./n8n-ops list --env production --search "customer"

# Mostrar IDs de workflows
./n8n-ops list --env production --show-ids
```

## 📊 Ejemplo de Salida del Comando Status

```
        ad88888ba       88888888ba          ,ad8888ba,   
       d8"     "8b      88      "8b        d8"'    "8b  
       88       88      88      ,8P       d8'        8b 
       ...

Workflow Status - PRODUCTION Environment
Last updated: 2025-07-22 21:15:30

NAME                    STATUS      VERSION   LAST DEPLOY   GIT COMMIT   NODES
────                    ──────      ───────   ───────────   ──────────   ─────
Customer Onboarding     🟢 Active   v1.2.3    2025-07-20    a1b2c3d      8
Payment Processing      🟢 Active   v2.1.0    2025-07-19    e4f5g6h      12
Email Notifications     🟢 Active   v1.0.5    2025-07-18    i7j8k9l      5
Data Backup             🔴 Inactive v1.0.0    2025-07-15    m1n2o3p      3

Total workflows: 4
Use --all to show inactive workflows
```

## 🎯 Casos de Uso

### 1. **Verificar versiones antes de deploy**
```bash
./n8n-ops status --env production
# Ve qué versión está actualmente
# Compara con tu versión local antes de hacer deploy
```

### 2. **Auditoría de producción**  
```bash
./n8n-ops status --env production --all --json > prod-status.json
# Exporta estado completo para auditoría
```

### 3. **Encontrar workflows específicos**
```bash
./n8n-ops list --env production --search "payment" --show-ids
# Busca workflows relacionados con pagos
```

### 4. **Comparar entre ambientes**
```bash
./n8n-ops status --env staging
./n8n-ops status --env production
# Compara qué versiones tienes en cada ambiente
```

## 💾 Datos Mostrados

Los comandos obtienen información de:

1. **n8n API**: Estado actual, nombres, nodos
2. **Base de datos local**: Historial de deployments, versiones, commits
3. **Git**: Información de commits y versiones

La información se almacena localmente cada vez que haces un deploy, permitiendo rastrear:
- Cuándo se hizo cada deploy
- Qué commit Git se desplegó
- Qué versión se asignó
- Estado del deployment (exitoso/fallido)

¡Ahora puedes ver exactamente qué versión de cada workflow está corriendo en producción!