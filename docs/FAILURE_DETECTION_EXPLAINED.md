# Cómo Funciona la Detección de Fallos en n8n-ops

## 📊 Resumen del Sistema

El sistema de detección de fallos de n8n-ops utiliza **múltiples métodos** para identificar cuando workflows están fallando y necesitan atención:

### 🔍 **Métodos de Detección**

#### 1. **Consulta Periódica de Ejecuciones (Polling)**
```bash
# El sistema consulta cada minuto (configurable)
n8n-ops monitor --check-interval 30s  # Cada 30 segundos
```

**Proceso:**
- Conecta a la API de n8n cada intervalo configurado
- Consulta las ejecuciones recientes de todos los workflows activos
- Analiza el estado de cada ejecución: `success`, `error`, `waiting`, `running`

#### 2. **Análisis de Estado de Ejecuciones**
```json
{
  "id": "exec_1003",
  "workflowId": "1002", 
  "status": "error",
  "error": {
    "message": "API rate limit exceeded",
    "node": "Stripe API",
    "stack": "Stripe Error: API rate limit exceeded\n    at Stripe.makeRequest()"
  },
  "startedAt": "2025-07-23T09:00:00Z",
  "stoppedAt": "2025-07-23T09:02:00Z"
}
```

#### 3. **Contador de Fallos Consecutivos**
El sistema mantiene un contador por workflow:
```go
failureCounts := map[string]int{
    "1001": 1,  // Customer Onboarding - 1 fallo
    "1002": 3,  // Payment Processing - 3 fallos → CREAR ISSUE
    "1003": 0,  // Order Fulfillment - funcionando
}
```

#### 4. **Umbral Configurable**
```bash
# Crear issue después de 2 fallos consecutivos
n8n-ops monitor --failure-threshold 2

# Para workflows críticos - issue inmediato
n8n-ops monitor --failure-threshold 1
```

## 🎯 **Algoritmo de Detección Paso a Paso**

### **Paso 1: Consulta de Ejecuciones**
```go
// Cada intervalo (ej: 10 segundos)
executions, err := n8nClient.GetExecutions(ctx, workflowID, 10)
```

### **Paso 2: Análisis de Estado**
```go
for _, execution := range executions {
    if execution.Status == "error" {
        // Incrementar contador de fallos
        failureCounts[workflowID]++
        
        logger.WithFields(logrus.Fields{
            "workflowId":   workflowID,
            "executionId":  execution.ID,
            "error":        execution.Error.Message,
        }).Warn("Workflow execution failure detected")
    }
}
```

### **Paso 3: Evaluación del Umbral**
```go
if failureCounts[workflowID] >= retryThreshold {
    // CREAR ISSUE EN GITLAB
    failure := &WorkflowFailure{
        WorkflowID:   workflowID,
        WorkflowName: workflow.Name,
        ErrorMessage: execution.Error.Message,
        RetryCount:   failureCounts[workflowID],
    }
    
    issueURL, err := issueManager.CreateWorkflowFailureIssue(ctx, failure)
    fmt.Printf("🚨 Issue created: %s\n", issueURL)
}
```

### **Paso 4: Detección de Recuperación**
```go
if execution.Status == "success" && failureCounts[workflowID] > 0 {
    // WORKFLOW SE RECUPERÓ
    logger.Info("Workflow recovered", "workflowId", workflowID)
    
    // Resetear contador
    failureCounts[workflowID] = 0
    
    // Actualizar issue existente
    issueManager.UpdateIssueWithRecovery(ctx, workflowID)
}
```

## 🎮 **Lo Que Estás Viendo en Demo**

En los logs que ves, el sistema está:

```bash
WARN[0010] Workflow execution failure detected           
error="&{API rate limit exceeded Stripe Error: API rate limit exceeded\n    at Stripe.makeRequest()}" 
executionId=exec_1003 workflowId=1002

WARN[0010] Workflow failure detected workflowId=1002
```

**Traducción:**
1. **Mock server** simula ejecución fallida en workflow 1002 (Payment Processing)
2. **Detector** encuentra el error y incrementa contador
3. **Después de 2 fallos consecutivos** → Crea issue automáticamente
4. **GitLab issue** se crea con detalles completos del fallo

## 🔧 **Configuraciones Avanzadas**

### **Por Ambiente**
```bash
# Producción - muy sensible
n8n-ops monitor --env production --failure-threshold 1 --check-interval 30s

# Desarrollo - más tolerante  
n8n-ops monitor --env development --failure-threshold 5 --check-interval 2m
```

### **Por Tipo de Workflow**
```yaml
# Configuración futura
critical-workflows:
  - id: "payment-processing"
    threshold: 1        # Issue inmediato
    interval: "15s"     # Verificar cada 15s
    
non-critical:
  - id: "email-notifications"  
    threshold: 10       # Más tolerante
    interval: "5m"      # Verificar cada 5 min
```

## 🏗️ **Arquitectura del Sistema**

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   n8n Instance  │    │  Failure        │    │  GitLab Issues  │
│                 │────│  Detector       │────│                 │
│ • Executions    │    │                 │    │ • Auto Creation │
│ • Workflows     │    │ • Polling       │    │ • Recovery      │
│ • API Status    │    │ • Counting      │    │ • Tracking      │
└─────────────────┘    │ • Alerting      │    └─────────────────┘
                       └─────────────────┘
                               │
                       ┌─────────────────┐
                       │  Observability  │
                       │                 │
                       │ • Sentry        │
                       │ • Grafana       │ 
                       │ • Metrics       │
                       └─────────────────┘
```

## 📈 **Tipos de Fallos Detectados**

### **1. Fallos de Ejecución**
- API timeouts
- Rate limits 
- Database connection errors
- Authentication failures

### **2. Fallos de Red**
- Connection timeouts
- DNS resolution errors
- SSL certificate issues

### **3. Fallos de Datos**
- Invalid JSON responses
- Missing required fields
- Data validation errors

### **4. Fallos de Lógica**
- Conditional branch errors
- Variable reference errors
- Expression evaluation failures

## 🚨 **Ejemplo en Vivo**

Lo que ves en tu monitor es exactamente esto funcionando:

```bash
# 🔍 Sistema detecta fallo
WARN[0020] Workflow execution failure detected 
error="API rate limit exceeded" workflowId=1002

# 📊 Contador incrementa
WARN[0020] Workflow failure detected workflowId=1002  

# 🎯 Alcanza umbral (2 fallos)
🔍 DEMO: Would create issue for workflow failure
   Workflow: Payment Processing (1002)
   Environment: development
   Error: Multiple consecutive execution failures detected
   Retry Count: 2

# 🚨 Issue creada automáticamente  
INFO[0020] Created failure issue for workflow
issueURL="https://gitlab.example.com/project/issues/1" workflowId=1002

🚨 Issue created for workflow failure: https://gitlab.example.com/project/issues/1
```

Este sistema te da **visibilidad completa** de la salud de tus workflows y **respuesta automática** cuando algo falla. ¿Te gustaría que ajuste algún parámetro o te explique otra parte del sistema?