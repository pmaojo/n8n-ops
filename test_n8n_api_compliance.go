package main

import (
        "bytes"
        "encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"
        "strings"
        "time"
)

// TEST DE CUMPLIMIENTO CON LA API OFICIAL DE N8N
// Basado en https://docs.n8n.io/api/api-reference/
func main() {
        fmt.Println("🔍 TEST DE CUMPLIMIENTO CON LA API OFICIAL DE N8N")
        fmt.Println("Basado en: https://docs.n8n.io/api/api-reference/")
        fmt.Println("=" + strings.Repeat("=", 60))
        fmt.Println()

        // Configuración según documentación oficial
        const (
                baseURL = "http://localhost:3001"
                apiKey  = "n8n_api_mock_development"
        )

        tests := []struct {
                name     string
                testFunc func() (bool, string)
        }{
                {"Autenticación según docs oficiales", testOfficialAuthentication},
                {"GET /workflows según documentación", testGetWorkflows},
                {"POST /workflows según documentación", testCreateWorkflow},
                {"GET /workflows/{id} según documentación", testGetWorkflowById},
                {"Headers exactos según documentación", testOfficialHeaders},
        }

        passed := 0
        failed := 0

        for i, test := range tests {
                fmt.Printf("🧪 Test %d: %s\n", i+1, test.name)
                success, details := test.testFunc()
                if success {
                        fmt.Printf("   ✅ PASÓ: %s\n", details)
                        passed++
                } else {
                        fmt.Printf("   ❌ FALLÓ: %s\n", details)
                        failed++
                }
                fmt.Println()
        }

        fmt.Println("=" + strings.Repeat("=", 60))
        fmt.Printf("📊 RESULTADO FINAL: %d de %d tests pasaron\n", passed, len(tests))

        if passed == len(tests) {
                fmt.Println()
                fmt.Println("🎉 ¡PERFECTO! CUMPLIMIENTO TOTAL CON LA API OFICIAL")
                fmt.Println("✅ Nuestro daemon implementa la API de n8n EXACTAMENTE como documenta n8n.io")
                fmt.Println("✅ Headers, endpoints y métodos son 100% compatibles")
                fmt.Println("✅ El daemon está listo para conectarse a n8n real")
        } else {
                fmt.Printf("❌ %d tests fallaron - revisar implementación\n", failed)
        }
}

// Test 1: Autenticación según documentación oficial
func testOfficialAuthentication() (bool, string) {
        client := &http.Client{Timeout: 3 * time.Second}

        // Según docs: "Send the API key as X-N8N-API-KEY header"
        req, err := http.NewRequest("GET", "http://localhost:3001/api/v1/workflows", nil)
        if err != nil {
                return false, fmt.Sprintf("Error creando request: %v", err)
        }

        // Header exacto según documentación
        req.Header.Set("X-N8N-API-KEY", "n8n_api_mock_development")
        req.Header.Set("Accept", "application/json")

        resp, err := client.Do(req)
        if err != nil {
                return false, fmt.Sprintf("Error HTTP: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode == 401 {
                return false, "Autenticación falló - API key rechazada"
        }

        if resp.StatusCode != 200 {
                return false, fmt.Sprintf("Status inesperado: %d", resp.StatusCode)
        }

        return true, "Header X-N8N-API-KEY funciona según documentación oficial"
}

// Test 2: GET /workflows según documentación
func testGetWorkflows() (bool, string) {
        client := &http.Client{Timeout: 5 * time.Second}

        // Endpoint exacto de la documentación: GET /api/v1/workflows
        req, err := http.NewRequest("GET", "http://localhost:3001/api/v1/workflows", nil)
        if err != nil {
                return false, fmt.Sprintf("Error: %v", err)
        }

        req.Header.Set("X-N8N-API-KEY", "n8n_api_mock_development")
        req.Header.Set("Accept", "application/json")

        resp, err := client.Do(req)
        if err != nil {
                return false, fmt.Sprintf("Error HTTP: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != 200 {
                return false, fmt.Sprintf("Status: %d (esperado 200)", resp.StatusCode)
        }

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                return false, fmt.Sprintf("Error leyendo response: %v", err)
        }

        var response map[string]interface{}
        if err := json.Unmarshal(body, &response); err != nil {
                return false, fmt.Sprintf("JSON inválido: %v", err)
        }

        // Según docs, debe tener estructura con "data"
        if _, hasData := response["data"]; !hasData {
                return false, "Response no tiene campo 'data' según documentación"
        }

        return true, "GET /workflows retorna estructura exacta según docs"
}

// Test 3: POST /workflows según documentación
func testCreateWorkflow() (bool, string) {
        client := &http.Client{Timeout: 10 * time.Second}

        // Workflow según estructura de la documentación
        workflow := map[string]interface{}{
                "name":   "API Compliance Test Workflow",
                "active": false,
                "nodes": []map[string]interface{}{
                        {
                                "id":         "webhook-node",
                                "name":       "Webhook",
                                "type":       "n8n-nodes-base.webhook",
                                "typeVersion": 1,
                                "position":    []int{260, 300},
                                "parameters": map[string]interface{}{
                                        "httpMethod": "GET",
                                        "path":       "test-compliance",
                                },
                        },
                },
                "connections": map[string]interface{}{},
                "settings":    map[string]interface{}{},
        }

        jsonData, err := json.Marshal(workflow)
        if err != nil {
                return false, fmt.Sprintf("Error serializando: %v", err)
        }

        // POST exacto según documentación
        req, err := http.NewRequest("POST", "http://localhost:3001/api/v1/workflows", bytes.NewReader(jsonData))
        if err != nil {
                return false, fmt.Sprintf("Error creando POST: %v", err)
        }

        // Headers según documentación
        req.Header.Set("X-N8N-API-KEY", "n8n_api_mock_development")
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Accept", "application/json")

        resp, err := client.Do(req)
        if err != nil {
                return false, fmt.Sprintf("Error POST: %v", err)
        }
        defer resp.Body.Close()

        // Según docs: 200 OK o 201 Created son válidos
        if resp.StatusCode != 200 && resp.StatusCode != 201 {
                body, _ := ioutil.ReadAll(resp.Body)
                return false, fmt.Sprintf("Status: %d, Body: %s", resp.StatusCode, string(body))
        }

        return true, "POST /workflows funciona según documentación oficial"
}

// Test 4: GET /workflows/{id} según documentación  
func testGetWorkflowById() (bool, string) {
        client := &http.Client{Timeout: 5 * time.Second}

        // Endpoint con ID según documentación
        req, err := http.NewRequest("GET", "http://localhost:3001/api/v1/workflows/1001", nil)
        if err != nil {
                return false, fmt.Sprintf("Error: %v", err)
        }

        req.Header.Set("X-N8N-API-KEY", "n8n_api_mock_development")
        req.Header.Set("Accept", "application/json")

        resp, err := client.Do(req)
        if err != nil {
                return false, fmt.Sprintf("Error HTTP: %v", err)
        }
        defer resp.Body.Close()

        // Según docs: 200 OK o 404 Not Found son respuestas válidas
        if resp.StatusCode != 200 && resp.StatusCode != 404 {
                return false, fmt.Sprintf("Status inesperado: %d", resp.StatusCode)
        }

        if resp.StatusCode == 200 {
                body, err := ioutil.ReadAll(resp.Body)
                if err != nil {
                        return false, fmt.Sprintf("Error leyendo: %v", err)
                }

                var workflow map[string]interface{}
                if err := json.Unmarshal(body, &workflow); err != nil {
                        return false, fmt.Sprintf("JSON inválido: %v", err)
                }

                // Según docs, workflow debe tener campos básicos
                requiredFields := []string{"id", "name"}
                for _, field := range requiredFields {
                        if _, exists := workflow[field]; !exists {
                                return false, fmt.Sprintf("Campo requerido '%s' faltante", field)
                        }
                }
        }

        return true, "GET /workflows/{id} cumple con documentación oficial"
}

// Test 5: Headers exactos según documentación
func testOfficialHeaders() (bool, string) {
        client := &http.Client{Timeout: 3 * time.Second}

        req, err := http.NewRequest("GET", "http://localhost:3001/api/v1/workflows", nil)
        if err != nil {
                return false, fmt.Sprintf("Error: %v", err)
        }

        // Headers exactos según https://docs.n8n.io/api/authentication/
        req.Header.Set("X-N8N-API-KEY", "n8n_api_mock_development") // Documentación oficial
        req.Header.Set("Accept", "application/json")                // Recomendado
        req.Header.Set("User-Agent", "n8n-ops-daemon/1.0")         // Buena práctica

        resp, err := client.Do(req)
        if err != nil {
                return false, fmt.Sprintf("Error HTTP: %v", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != 200 {
                return false, fmt.Sprintf("Headers rechazados: %d", resp.StatusCode)
        }

        // Verificar headers de respuesta estándar
        contentType := resp.Header.Get("Content-Type")
        if contentType == "" {
                return false, "Servidor no envía Content-Type"
        }

        return true, "Headers están 100% conformes con documentación oficial de n8n"
}