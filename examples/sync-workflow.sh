#!/bin/bash
# Ejemplo de script para sincronizar workflows desde n8n

set -e  # Exit on any error

echo "🔄 n8n Workflow Sync Script"
echo "=========================="

# Verificar que n8n-ops existe
if [[ ! -f "./n8n-ops" ]]; then
    echo "❌ n8n-ops binary not found. Run 'go build -o n8n-ops main.go' first"
    exit 1
fi

# Verificar variables de entorno
check_env_vars() {
    local env=$1
    local url_var="N8N_URL_${env^^}"
    local key_var="N8N_API_KEY_${env^^}"
    
    if [[ -z "${!url_var}" || -z "${!key_var}" ]]; then
        echo "❌ Missing environment variables for $env:"
        echo "   Required: $url_var and $key_var"
        return 1
    fi
    echo "✅ Environment variables for $env are set"
}

# Función para sincronizar un ambiente
sync_environment() {
    local env=$1
    echo ""
    echo "📥 Syncing $env environment..."
    echo "----------------------------------------"
    
    # Verificar variables de entorno
    if ! check_env_vars "$env"; then
        echo "⚠️  Skipping $env due to missing credentials"
        return 0
    fi
    
    # Crear directorio si no existe
    mkdir -p "./workflows/$env"
    
    # Hacer backup de archivos existentes
    if [[ -n "$(ls -A "./workflows/$env" 2>/dev/null)" ]]; then
        echo "💾 Creating backup of existing workflows..."
        backup_dir="./workflows/$env.backup.$(date +%Y%m%d_%H%M%S)"
        cp -r "./workflows/$env" "$backup_dir"
        echo "   Backup saved to: $backup_dir"
    fi
    
    # Ejecutar sync
    echo "🔄 Downloading workflows from $env..."
    if ./n8n-ops sync --env "$env" --verbose; then
        echo "✅ Successfully synced $env workflows"
        
        # Mostrar status
        echo ""
        echo "📊 Current status:"
        ./n8n-ops status --env "$env"
        
        # Mostrar cambios en Git
        if git status --porcelain "./workflows/$env" | grep -q .; then
            echo ""
            echo "📝 Git changes detected:"
            git status "./workflows/$env"
            echo ""
            echo "💡 To commit changes:"
            echo "   git add ./workflows/$env/"
            echo "   git commit -m \"sync: update $env workflows\""
        else
            echo ""
            echo "ℹ️  No changes detected in Git"
        fi
    else
        echo "❌ Failed to sync $env workflows"
        return 1
    fi
}

# Función principal
main() {
    # Argumentos
    ENVIRONMENT=${1:-""}
    
    if [[ -n "$ENVIRONMENT" ]]; then
        # Sincronizar ambiente específico
        sync_environment "$ENVIRONMENT"
    else
        # Sincronizar todos los ambientes
        echo "🔄 Syncing all environments..."
        echo ""
        
        for env in "development" "staging" "production"; do
            sync_environment "$env"
            echo ""
        done
        
        echo "🎉 Sync process completed for all environments"
        echo ""
        echo "📋 Summary:"
        echo "   Use './n8n-ops status --env [ENV]' to check workflow status"
        echo "   Use 'git status' to see all changes"
        echo "   Use 'git add . && git commit -m \"sync: update workflows\"' to save changes"
    fi
}

# Mostrar ayuda
if [[ "$1" == "--help" || "$1" == "-h" ]]; then
    echo "Usage: $0 [environment]"
    echo ""
    echo "Arguments:"
    echo "  environment    Sync specific environment (development, staging, production)"
    echo "                 If not specified, syncs all environments"
    echo ""
    echo "Examples:"
    echo "  $0                    # Sync all environments"
    echo "  $0 production        # Sync only production"
    echo "  $0 development       # Sync only development"
    echo ""
    echo "Required Environment Variables:"
    echo "  N8N_URL_DEV, N8N_API_KEY_DEV          # For development"
    echo "  N8N_URL_STAGING, N8N_API_KEY_STAGING  # For staging"
    echo "  N8N_URL_PROD, N8N_API_KEY_PROD        # For production"
    exit 0
fi

# Ejecutar función principal
main "$@"