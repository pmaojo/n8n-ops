package i18n

import (
	"fmt"
)

// Language represents a supported localization option.
type Language string

const (
	English Language = "en"
	Spanish Language = "es"
)

var currentLang Language = English

// SetLanguage sets the current language
func SetLanguage(lang string) {
	switch lang {
	case "es", "spanish", "español":
		currentLang = Spanish
	default:
		currentLang = English
	}
}

// GetLanguage returns the current language
func GetLanguage() Language {
	return currentLang
}

// Messages in multiple languages
var messages = map[Language]map[string]string{
	English: {
		"welcome_title":        "Welcome to the Future of Workflow Automation",
		"powered_by":           "Powered by n8n Technology",
		"lightning_fast":       "Lightning-fast Git Integration",
		"multi_env":            "Multi-environment Support",
		"enterprise_security":  "Enterprise-grade Security",
		"begin_journey":        "Type 'n8n-ops --help' to begin your journey...",
		"workflow_automation":  "Workflow Automation CLI",
		"current_branch":       "Current branch",
		"environment":          "Environment",
		"branch_env_mapping":   "Branch → Environment mapping",
		"branches":             "Branches",
		"workflow":             "Workflow",
		"status":               "Status",
		"env":                  "Env",
		"success":              "SUCCESS!",
		"error_detected":       "ERROR DETECTED",
		"starting_sync":        "Starting workflow sync",
		"sync_completed":       "Workflows synchronized successfully",
		"deploy_completed":     "Workflows deployed successfully",
		"validation_completed": "Workflow validation completed",
		"rollback_completed":   "Rollback completed successfully",
		"creating_branch":      "Creating new branch",
		"switching_branch":     "Switching to branch",
		"branch_created":       "Created and switched to branch",
		"branch_switched":      "Switched to branch",
		"workflows_processed":  "workflows processed",
		"command":              "COMMAND",
		"not_git_repo":         "not in a Git repository",
		"failed_get_branches":  "failed to get branch list",
		"failed_create_branch": "failed to create branch",
		"failed_switch_branch": "failed to switch to branch",
	},
	Spanish: {
		"welcome_title":        "Bienvenido al Futuro de la Automatización de Workflows",
		"powered_by":           "Impulsado por Tecnología n8n",
		"lightning_fast":       "Integración Git Ultra-rápida",
		"multi_env":            "Soporte Multi-entorno",
		"enterprise_security":  "Seguridad Nivel Empresarial",
		"begin_journey":        "Escribe 'n8n-ops --help' para comenzar tu viaje...",
		"workflow_automation":  "CLI de Automatización de Workflows",
		"current_branch":       "Rama actual",
		"environment":          "Entorno",
		"branch_env_mapping":   "Mapeo Rama → Entorno",
		"branches":             "Ramas",
		"workflow":             "Workflow",
		"status":               "Estado",
		"env":                  "Entorno",
		"success":              "¡ÉXITO!",
		"error_detected":       "ERROR DETECTADO",
		"starting_sync":        "Iniciando sincronización de workflows",
		"sync_completed":       "Workflows sincronizados exitosamente",
		"deploy_completed":     "Workflows desplegados exitosamente",
		"validation_completed": "Validación de workflows completada",
		"rollback_completed":   "Rollback completado exitosamente",
		"creating_branch":      "Creando nueva rama",
		"switching_branch":     "Cambiando a rama",
		"branch_created":       "Rama creada y cambiada a",
		"branch_switched":      "Cambiado a rama",
		"workflows_processed":  "workflows procesados",
		"command":              "COMANDO",
		"not_git_repo":         "no está en un repositorio Git",
		"failed_get_branches":  "falló al obtener lista de ramas",
		"failed_create_branch": "falló al crear rama",
		"failed_switch_branch": "falló al cambiar a rama",
	},
}

// T translates a message key to the current language
func T(key string) string {
	if msg, exists := messages[currentLang][key]; exists {
		return msg
	}
	// Fallback to English if key not found
	if msg, exists := messages[English][key]; exists {
		return msg
	}
	return key // Return key if no translation found
}

// Tf translates and formats a message
func Tf(key string, args ...interface{}) string {
	template := T(key)
	return fmt.Sprintf(template, args...)
}
