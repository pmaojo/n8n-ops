# API Monitoring Integration Guide

## Descripción del Sistema

El sistema de monitoreo utiliza la **API real de n8n** para detectar fallos de workflows y crear issues automáticamente en GitLab. La integración se basa en consultas periódicas a los endpoints de executions de n8n.

## Flujo de Detección de Fallos

### 1. Consulta de Workflows Activos
```bash
GET /api/v1/workflows
```

El sistema obtiene la lista de workflows activos y filtra solo los que están habilitados.

### 2. Consulta de Executions por Workflow
```bash
GET /api/v1/executions?workflowId={id}&limit=10
```

Para cada workflow activo, el sistema consulta las últimas 10 executions para detectar fallos recientes.

### 3. Detección de Errores
El sistema verifica si la última execution tiene:
- **Status**: `"error"`
- **Error Details**: Mensaje específico, nodo que falló, stack trace

### 4. Creación de Issues
Cuando se detectan fallos consecutivos que superan el umbral configurado, se crea automáticamente una issue en GitLab con:

```json
{
  "title": "🚨 Workflow Failure: Payment Processing (production)",
  "labels": ["workflow-failure", "automated", "env:production", "severity:high"],
  "description": "## Error Details\n**Message**: Database connection timeout\n**Node**: PostgreSQL\n**Execution ID**: exec_1002"
}
```

## Estructura de Datos de Execution

### ExecutionResult (n8n API Response)
```go
type ExecutionResult struct {
    ID           string          `json:"id"`
    WorkflowID   string          `json:"workflowId"`
    WorkflowName string          `json:"workflowName"`
    Status       string          `json:"status"`        // "success", "error", "running"
    StartedAt    time.Time       `json:"startedAt"`
    StoppedAt    *time.Time      `json:"stoppedAt"`
    Mode         string          `json:"mode"`          // "trigger", "manual", "retry"
    Retries      int             `json:"retries"`
    Error        *ExecutionError `json:"error"`         // Null if no error
}
```

### ExecutionError (Detalles del Error)
```go
type ExecutionError struct {
    Message string `json:"message"`    // "Database connection timeout"
    Node    string `json:"node"`       // "PostgreSQL" 
    Stack   string `json:"stack"`      // Stack trace completo
}
```

## Endpoints de n8n Utilizados

### GET /api/v1/executions
- **Propósito**: Obtener lista de executions
- **Parámetros**:
  - `workflowId`: ID del workflow específico
  - `status`: Filtrar por estado (`error`, `success`, etc.)
  - `limit`: Número máximo de resultados

### GET /api/v1/executions/{id}
- **Propósito**: Obtener detalles de execution específica
- **Respuesta**: ExecutionResult completo con error details

## Configuración de Monitoreo

### Variables de Entorno Necesarias
```bash
# n8n API Integration
N8N_API_KEY=your-n8n-api-key
N8N_URL=https://your-n8n-instance.com

# GitLab Integration  
GITLAB_TOKEN=your-gitlab-token
GITLAB_PROJECT_ID=12345
GITLAB_URL=https://gitlab.com
```

### Configuración de Umbrales
```bash
n8n-ops monitor \
  --failure-threshold 3 \      # Crear issue después de 3 fallos
  --check-interval 1m \        # Verificar cada minuto
  --env production            # Ambiente de producción
```

## Mock Server para Testing

Para pruebas y desarrollo, se incluye un mock server que simula la API de n8n:

### Datos de Prueba
```json
{
  "exec_1002": {
    "id": "exec_1002",
    "workflowId": "1001", 
    "workflowName": "Customer Onboarding",
    "status": "error",
    "startedAt": "2025-07-23T09:00:00Z",
    "stoppedAt": "2025-07-23T09:02:00Z",
    "mode": "trigger",
    "retries": 2,
    "error": {
      "message": "Database connection timeout",
      "node": "PostgreSQL",
      "stack": "Error: Database connection timeout\n    at PostgreSQL.execute()"
    }
  }
}
```

### Pruebas con Mock Server
```bash
# Iniciar mock server
cd mock-n8n-server && go run main.go

# Probar endpoint de executions
curl -H "X-N8N-API-KEY: n8n_api_mock_development" \
  "http://localhost:3001/api/v1/executions?workflowId=1002"

# Iniciar monitoreo en modo demo
n8n-ops monitor --demo --env development
```

## Detección Inteligente de Fallos

### Algoritmo de Detección
1. **Consulta periódica** de executions para workflows activos
2. **Análisis de la última execution** para detectar status "error"
3. **Contador de fallos consecutivos** por workflow
4. **Threshold configurable** antes de crear issue
5. **Detección de recuperación** automática

### Estados de Workflow
- **Healthy**: Última execution exitosa
- **Failing**: Una o más executions fallidas, pero bajo el threshold
- **Critical**: Fallos consecutivos superan threshold → **Issue creada**
- **Recovered**: Workflow vuelve a funcionar → **Issue actualizada**

## Integración con GitLab

### Creación Automática de Issues
Cuando un workflow supera el threshold de fallos:

```go
failure := &WorkflowFailure{
    WorkflowID:   "wf_payment_123",
    WorkflowName: "Payment Processing",
    ExecutionID:  "exec_1002", 
    Environment:  "production",
    ErrorMessage: "Database connection timeout",
    NodeName:     "PostgreSQL",
    RetryCount:   3,
}

issue, err := issueManager.CreateWorkflowFailureIssue(ctx, failure)
```

### Actualización por Recuperación
Cuando el workflow se recupera:

```go
recovery := &RecoveryInfo{
    RecoveredAt:  time.Now(),
    RecoveryType: "auto",        // "auto", "manual", "rollback"
    Notes:        "Database connection restored",
}

err := issueManager.UpdateIssueWithRecovery(ctx, issue.ID, recovery)
```

## Ventajas de la Integración API

### 1. **Datos Reales**
- No simulaciones, usa execution logs reales de n8n
- Errores específicos con stack traces completos
- Información precisa de nodos que fallan

### 2. **Tiempo Real**
- Detección inmediata de fallos (intervalo configurable)
- No depende de logs externos o webhooks
- Polling directo de la API de n8n

### 3. **Contexto Completo**
- ID de execution específica para debugging
- Información del nodo que causó el error
- Número de reintentos y modo de execution

### 4. **Integración GitOps**
- Issues automáticas con toda la información técnica
- Labels automáticas para categorización
- Tracking de recuperación y resolución

## Troubleshooting

### Errores Comunes

**1. API Key Inválida**
```bash
Error: failed to get workflows: HTTP 401 Unauthorized
```
**Solución**: Verificar N8N_API_KEY en variables de entorno

**2. n8n Instance No Disponible**
```bash
Error: failed to connect to n8n: connection refused
```
**Solución**: Verificar N8N_URL y conectividad de red

**3. Executions Vacías**
```bash
Warning: No executions found for workflow wf_123
```
**Solución**: Normal para workflows nuevos, no es un error

### Logs de Debug
```bash
n8n-ops monitor --verbose --env development
```

El sistema proporciona logs detallados de cada consulta API y decisión de monitoring.

## Mejores Prácticas

### 1. **Configuración por Ambiente**
- **Development**: Thresholds altos, intervals largos
- **Staging**: Configuración intermedia
- **Production**: Thresholds bajos, monitoreo frecuente

### 2. **Gestión de Rate Limits**
- Intervals apropiados para evitar sobrecargar la API de n8n
- Configuración de timeouts adecuados
- Handling de errores de conectividad

### 3. **Categorización de Issues**
- Labels automáticas por severidad y ambiente
- Asignación de equipos según tipo de workflow
- Templates de troubleshooting específicos

La integración con la API de n8n proporciona un sistema de monitoreo robusto que detecta fallos reales y proporciona el contexto necesario para resolución rápida de problemas.