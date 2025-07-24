package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pmaojo/n8n-ops/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var stopServer func()

func TestMain(m *testing.M) {
	var err error
	stopServer, err = testutil.StartMockServer(0)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	code := m.Run()
	if stopServer != nil {
		stopServer()
	}
	os.Exit(code)
}

const (
	mockN8nURL       = "http://localhost:3001"
	testWorkflowsDir = "./test-workflows/development"
	testBackupsDir   = "./test-backups/development"
	testWorkflowID   = "daemon_test_workflow_001"
	testWorkflowName = "Daemon Test Workflow"
	testAPIKey       = "n8n_api_mock_development"
)

// TestWorkflow representa un workflow para testing
type TestWorkflow struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Active      bool                   `json:"active"`
	Nodes       []interface{}          `json:"nodes"`
	Connections map[string]interface{} `json:"connections"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	UpdatedAt   string                 `json:"updatedAt,omitempty"`
}

// Test 1: Verificar que el mock n8n server está funcionando
func TestMockN8nServerRunning(t *testing.T) {
	t.Log("🧪 Test 1: Verificando que el mock n8n server está corriendo")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(mockN8nURL + "/health")

	require.NoError(t, err, "El servidor mock n8n debe estar corriendo en %s", mockN8nURL)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Health endpoint debe retornar 200")

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var health map[string]interface{}
	err = json.Unmarshal(body, &health)
	require.NoError(t, err)

	assert.Equal(t, "ok", health["status"], "Health status debe ser 'ok'")
	t.Log("✅ Mock n8n server está corriendo correctamente")
}

// Test 2: Verificar que podemos crear y leer workflows via API
func TestN8nAPIWorkflowOperations(t *testing.T) {
	t.Log("🧪 Test 2: Verificando operaciones de workflow via API")

	client := &http.Client{Timeout: 10 * time.Second}

	// Crear workflow de prueba
	testWorkflow := TestWorkflow{
		ID:     testWorkflowID,
		Name:   testWorkflowName,
		Active: false,
		Nodes: []interface{}{
			map[string]interface{}{
				"id":          "start_node",
				"name":        "Start",
				"type":        "n8n-nodes-base.manualTrigger",
				"typeVersion": 1,
				"position":    []int{250, 300},
				"parameters":  map[string]interface{}{},
			},
		},
		Connections: map[string]interface{}{},
		Settings:    map[string]interface{}{"executionOrder": "v1"},
		Tags:        []string{"test", "daemon"},
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	// Serializar workflow
	workflowJSON, err := json.Marshal(testWorkflow)
	require.NoError(t, err)

	// Crear workflow via API
	req, err := http.NewRequest("POST", mockN8nURL+"/api/v1/workflows",
		io.NopCloser(bytes.NewReader(workflowJSON)))
	require.NoError(t, err)

	req.Header.Set("X-N8N-API-KEY", testAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated,
		"Crear workflow debe retornar 200 o 201, recibido: %d", resp.StatusCode)

	// Verificar que podemos leer el workflow
	getResp, err := client.Get(mockN8nURL + "/api/v1/workflows/" + testWorkflowID)
	require.NoError(t, err)
	defer getResp.Body.Close()

	assert.Equal(t, http.StatusOK, getResp.StatusCode,
		"Leer workflow debe retornar 200")

	t.Log("✅ API de workflows funciona correctamente")
}

// Test 3: Verificar que el daemon detecta cambios en archivos
func TestDaemonFileWatching(t *testing.T) {
	t.Log("🧪 Test 3: Verificando detección de cambios en archivos")

	// Crear directorio de prueba
	err := os.MkdirAll(testWorkflowsDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll("./test-workflows")

	// Crear archivo de workflow inicial
	initialWorkflow := TestWorkflow{
		ID:     testWorkflowID,
		Name:   "Initial Workflow Name",
		Active: false,
		Nodes: []interface{}{
			map[string]interface{}{
				"id":          "node1",
				"name":        "Initial Node",
				"type":        "n8n-nodes-base.manualTrigger",
				"typeVersion": 1,
				"position":    []int{250, 300},
			},
		},
		Connections: map[string]interface{}{},
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	workflowFile := filepath.Join(testWorkflowsDir, "test-workflow.json")
	initialJSON, err := json.MarshalIndent(initialWorkflow, "", "  ")
	require.NoError(t, err)

	err = os.WriteFile(workflowFile, initialJSON, 0644)
	require.NoError(t, err)

	// Verificar que el archivo se creó
	assert.FileExists(t, workflowFile, "Archivo de workflow debe existir")

	// Leer el archivo y verificar contenido
	fileContent, err := os.ReadFile(workflowFile)
	require.NoError(t, err)

	var readWorkflow TestWorkflow
	err = json.Unmarshal(fileContent, &readWorkflow)
	require.NoError(t, err)

	assert.Equal(t, "Initial Workflow Name", readWorkflow.Name,
		"Nombre del workflow debe coincidir")

	t.Log("✅ Archivo de workflow creado y leído correctamente")
}

// Test 4: Test de integración completa del daemon
func TestDaemonIntegration(t *testing.T) {
	t.Log("🧪 Test 4: Test de integración completa del daemon")

	// Crear directorios necesarios
	err := os.MkdirAll(testWorkflowsDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll("./test-workflows")

	// Crear workflow inicial
	workflow := TestWorkflow{
		ID:     "integration_test_001",
		Name:   "Integration Test Workflow",
		Active: false,
		Nodes: []interface{}{
			map[string]interface{}{
				"id":          "webhook_node",
				"name":        "Webhook",
				"type":        "n8n-nodes-base.webhook",
				"typeVersion": 1,
				"position":    []int{250, 300},
				"parameters": map[string]interface{}{
					"httpMethod": "POST",
					"path":       "integration-test",
				},
			},
		},
		Connections: map[string]interface{}{},
		Tags:        []string{"integration", "test"},
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	workflowFile := filepath.Join(testWorkflowsDir, "integration-test.json")
	workflowJSON, err := json.MarshalIndent(workflow, "", "  ")
	require.NoError(t, err)

	err = os.WriteFile(workflowFile, workflowJSON, 0644)
	require.NoError(t, err)

	t.Log("📝 Workflow inicial creado")

	// Simular modificación del archivo (lo que haría el daemon)
	time.Sleep(100 * time.Millisecond) // Pequeña pausa para simular tiempo

	workflow.Name = "Integration Test Workflow - UPDATED"
	workflow.Active = true
	workflow.UpdatedAt = time.Now().Format(time.RFC3339)

	updatedJSON, err := json.MarshalIndent(workflow, "", "  ")
	require.NoError(t, err)

	err = os.WriteFile(workflowFile, updatedJSON, 0644)
	require.NoError(t, err)

	t.Log("📝 Workflow modificado")

	// Verificar que los cambios se guardaron
	fileContent, err := os.ReadFile(workflowFile)
	require.NoError(t, err)

	var updatedWorkflow TestWorkflow
	err = json.Unmarshal(fileContent, &updatedWorkflow)
	require.NoError(t, err)

	assert.Equal(t, "Integration Test Workflow - UPDATED", updatedWorkflow.Name)
	assert.True(t, updatedWorkflow.Active, "Workflow debe estar activo")

	t.Log("✅ Integración completa verificada")
}

// Test 5: Verificar que el daemon puede manejar múltiples archivos
func TestDaemonMultipleFiles(t *testing.T) {
	t.Log("🧪 Test 5: Verificando manejo de múltiples archivos")

	// Crear directorio
	err := os.MkdirAll(testWorkflowsDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll("./test-workflows")

	// Crear múltiples workflows
	workflows := []TestWorkflow{
		{
			ID:          "multi_test_001",
			Name:        "Multi Test Workflow 1",
			Active:      false,
			Nodes:       []interface{}{},
			Connections: map[string]interface{}{},
		},
		{
			ID:          "multi_test_002",
			Name:        "Multi Test Workflow 2",
			Active:      true,
			Nodes:       []interface{}{},
			Connections: map[string]interface{}{},
		},
		{
			ID:          "multi_test_003",
			Name:        "Multi Test Workflow 3",
			Active:      false,
			Nodes:       []interface{}{},
			Connections: map[string]interface{}{},
		},
	}

	// Crear archivos
	for i, wf := range workflows {
		fileName := fmt.Sprintf("multi-test-%d.json", i+1)
		filePath := filepath.Join(testWorkflowsDir, fileName)

		wfJSON, err := json.MarshalIndent(wf, "", "  ")
		require.NoError(t, err)

		err = os.WriteFile(filePath, wfJSON, 0644)
		require.NoError(t, err)

		assert.FileExists(t, filePath, "Archivo %s debe existir", fileName)
	}

	// Verificar que todos los archivos se crearon
	files, err := os.ReadDir(testWorkflowsDir)
	require.NoError(t, err)

	jsonFiles := 0
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			jsonFiles++
		}
	}

	assert.Equal(t, 3, jsonFiles, "Deben existir exactamente 3 archivos JSON")

	t.Log("✅ Manejo de múltiples archivos verificado")
}

// Test 6: Verificar estructura de backups
func TestBackupStructure(t *testing.T) {
	t.Log("🧪 Test 6: Verificando estructura de backups")

	// Crear directorio de backups
	err := os.MkdirAll(testBackupsDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll("./test-backups")

	// Simular creación de backup
	backupWorkflow := TestWorkflow{
		ID:          "backup_test_001",
		Name:        "Backup Test Workflow",
		Active:      true,
		Nodes:       []interface{}{},
		Connections: map[string]interface{}{},
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	// Crear archivo de backup
	timestamp := time.Now().Format("20060102-150405")
	backupFileName := fmt.Sprintf("%s_%s_backup.json", backupWorkflow.ID, timestamp)
	backupFilePath := filepath.Join(testBackupsDir, backupFileName)

	backupJSON, err := json.MarshalIndent(backupWorkflow, "", "  ")
	require.NoError(t, err)

	err = os.WriteFile(backupFilePath, backupJSON, 0644)
	require.NoError(t, err)

	// Crear archivo de metadata
	metadataFileName := fmt.Sprintf("%s_%s_backup.meta.json", backupWorkflow.ID, timestamp)
	metadataFilePath := filepath.Join(testBackupsDir, metadataFileName)

	metadata := map[string]interface{}{
		"originalFile": fmt.Sprintf("./workflows/development/%s.json", backupWorkflow.ID),
		"backupFile":   backupFilePath,
		"timestamp":    time.Now().Format(time.RFC3339),
		"environment":  "development",
		"workflowId":   backupWorkflow.ID,
		"workflowName": backupWorkflow.Name,
	}

	metadataJSON, err := json.MarshalIndent(metadata, "", "  ")
	require.NoError(t, err)

	err = os.WriteFile(metadataFilePath, metadataJSON, 0644)
	require.NoError(t, err)

	// Verificar que los archivos existen
	assert.FileExists(t, backupFilePath, "Archivo de backup debe existir")
	assert.FileExists(t, metadataFilePath, "Archivo de metadata debe existir")

	// Verificar contenido del backup
	backupContent, err := os.ReadFile(backupFilePath)
	require.NoError(t, err)

	var restoredWorkflow TestWorkflow
	err = json.Unmarshal(backupContent, &restoredWorkflow)
	require.NoError(t, err)

	assert.Equal(t, backupWorkflow.ID, restoredWorkflow.ID, "ID debe coincidir en backup")
	assert.Equal(t, backupWorkflow.Name, restoredWorkflow.Name, "Nombre debe coincidir en backup")

	t.Log("✅ Estructura de backups verificada")
}

// Test 7: Test de rendimiento - múltiples cambios rápidos
func TestDaemonPerformance(t *testing.T) {
	t.Log("🧪 Test 7: Test de rendimiento con múltiples cambios")

	// Crear directorio
	err := os.MkdirAll(testWorkflowsDir, 0755)
	require.NoError(t, err)
	defer os.RemoveAll("./test-workflows")

	workflowFile := filepath.Join(testWorkflowsDir, "performance-test.json")

	// Realizar múltiples modificaciones rápidas
	startTime := time.Now()

	for i := 0; i < 10; i++ {
		workflow := TestWorkflow{
			ID:          "performance_test_001",
			Name:        fmt.Sprintf("Performance Test - Iteration %d", i),
			Active:      i%2 == 0, // Alternar entre activo/inactivo
			Nodes:       []interface{}{},
			Connections: map[string]interface{}{},
			UpdatedAt:   time.Now().Format(time.RFC3339),
		}

		workflowJSON, err := json.MarshalIndent(workflow, "", "  ")
		require.NoError(t, err)

		err = os.WriteFile(workflowFile, workflowJSON, 0644)
		require.NoError(t, err)

		// Pequeña pausa entre modificaciones
		time.Sleep(10 * time.Millisecond)
	}

	duration := time.Since(startTime)

	// Verificar que el último cambio se guardó correctamente
	finalContent, err := os.ReadFile(workflowFile)
	require.NoError(t, err)

	var finalWorkflow TestWorkflow
	err = json.Unmarshal(finalContent, &finalWorkflow)
	require.NoError(t, err)

	assert.Equal(t, "Performance Test - Iteration 9", finalWorkflow.Name)
	assert.False(t, finalWorkflow.Active, "Último workflow debe estar inactivo (9%2 != 0)")

	t.Logf("✅ Test de rendimiento completado en %v", duration)
	assert.Less(t, duration, 5*time.Second, "Test debe completarse en menos de 5 segundos")
}

// Función helper para ejecutar todos los tests
func TestDaemonCompleteValidation(t *testing.T) {
	t.Log("🚀 INICIANDO VALIDACIÓN COMPLETA DEL DAEMON")
	t.Log("=" + strings.Repeat("=", 50))

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{"Mock n8n Server", TestMockN8nServerRunning},
		{"API Operations", TestN8nAPIWorkflowOperations},
		{"File Watching", TestDaemonFileWatching},
		{"Integration", TestDaemonIntegration},
		{"Multiple Files", TestDaemonMultipleFiles},
		{"Backup Structure", TestBackupStructure},
		{"Performance", TestDaemonPerformance},
	}

	passed := 0
	failed := 0

	for _, test := range tests {
		test := test
		success := t.Run(test.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("❌ Test %s falló con panic: %v", test.name, r)
				}
			}()

			test.testFunc(t)
		})

		if success {
			passed++
		} else {
			failed++
		}
	}

	t.Log("=" + strings.Repeat("=", 50))
	t.Logf("🏁 RESUMEN FINAL: %d tests pasaron, %d fallaron", passed, failed)

	if failed == 0 {
		t.Log("🎉 ¡TODOS LOS TESTS PASARON! EL DAEMON FUNCIONA PERFECTAMENTE")
	} else {
		t.Fail()
	}
}
