# n8n Workflows - Git Integration

## 🚨 Sistema de Detección de Workflows No Commiteados

Este directorio contiene workflows de n8n organizados por ambiente. El sistema n8n-ops detecta automáticamente workflows que han sido editados pero no commiteados a Git.

### Estructura

```
workflows/
├── development/     # Workflows de desarrollo
├── staging/        # Workflows de staging  
├── production/     # Workflows de producción
└── README.md
```

### Detección Automática

n8n-ops detecta cambios mediante:

1. **Git Status**: Compara archivos modificados, agregados, eliminados
2. **Workflow Analysis**: Identifica archivos .json en directorios de workflows
3. **Environment Mapping**: Determina el ambiente basado en la estructura de carpetas

### Comandos de Verificación

```bash
# Verificar cambios sin commitear
n8n-ops status --check-uncommitted

# Estado completo del sistema
n8n-ops status

# Sync bloqueado si hay cambios
n8n-ops sync --env development

# Forzar sync ignorando cambios
n8n-ops sync --env development --force
```



### Mejores Prácticas

1. **Commitear Frecuentemente**: Guarda cambios después de editar workflows
2. **Verificar Antes de Sync**: Usa `--check-uncommitted` antes de sincronizar
4. **Usar Branches**: Crea branches para features nuevos

El sistema previene pérdida de datos y facilita la colaboración en equipos.
