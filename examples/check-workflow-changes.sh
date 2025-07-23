#!/bin/bash
# Script para verificar cambios en workflows sin sincronizar

set -e

echo "🔍 Verificando cambios en workflows n8n"
echo "====================================="

# Función para verificar un ambiente específico
check_environment() {
    local env=$1
    echo ""
    echo "📊 Verificando ambiente: $env"
    echo "--------------------------------"
    
    # Verificar si el CLI existe
    if [[ ! -f "./n8n-ops" ]]; then
        echo "❌ n8n-ops binary not found. Run 'go build -o n8n-ops main.go' first"
        return 1
    fi
    
    # Verificar cambios
    if ./n8n-ops check --env "$env" --quiet; then
        echo "✅ Todos los workflows están sincronizados"
    else
        echo "⚠️  Cambios detectados en workflows:"
        ./n8n-ops check --env "$env"
        
        echo ""
        echo "🔄 ¿Quieres sincronizar ahora? (y/N)"
        read -r response
        if [[ "$response" =~ ^[Yy]$ ]]; then
            echo "Sincronizando workflows desde $env..."
            ./n8n-ops sync --env "$env" --verbose
            
            echo "📝 Mostrando cambios en Git:"
            git status "workflows/$env/"
            
            echo ""
            echo "💾 ¿Quieres hacer commit de los cambios? (y/N)"
            read -r commit_response
            if [[ "$commit_response" =~ ^[Yy]$ ]]; then
                git add "workflows/$env/"
                git commit -m "sync: update $env workflows from n8n instance"
                echo "✅ Cambios guardados en Git"
            fi
        fi
    fi
}

# Función para verificar todos los ambientes
check_all_environments() {
    echo "🌍 Verificando todos los ambientes..."
    
    for env in "development" "staging" "production"; do
        check_environment "$env"
    done
    
    echo ""
    echo "📋 Resumen completo:"
    echo "==================="
    
    # Generar reporte JSON de todos los ambientes
    for env in "development" "staging" "production"; do
        echo ""
        echo "📊 $env:"
        if ./n8n-ops check --env "$env" --json 2>/dev/null | jq -r '
            "  ✅ Sincronizados: \(.synchronized)",
            "  ⚠️  Modificados: \(.modified)", 
            "  📝 Total: \(.totalWorkflows)"
        '; then
            continue
        else
            echo "  ❌ Error verificando $env"
        fi
    done
}

# Función para monitoreo continuo
monitor_changes() {
    local env=${1:-"development"}
    local interval=${2:-300} # 5 minutos por defecto
    
    echo "👀 Monitoreando cambios en $env cada $interval segundos..."
    echo "Presiona Ctrl+C para detener"
    echo ""
    
    while true; do
        local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
        
        if ./n8n-ops check --env "$env" --quiet; then
            echo "[$timestamp] ✅ $env: Todos los workflows sincronizados"
        else
            echo "[$timestamp] ⚠️  $env: Cambios detectados"
            ./n8n-ops check --env "$env" --json | jq -r '
                "  Modificados: \(.modified) workflows",
                (.workflows.modified[] | "    - \(.name) (\(.timeAgo))")
            '
            
            # Opcional: Auto-sincronizar development
            if [[ "$env" == "development" ]]; then
                echo "🔄 Auto-sincronizando development..."
                ./n8n-ops sync --env development --verbose
            fi
        fi
        
        sleep "$interval"
    done
}

# Función para pre-commit hook
pre_commit_check() {
    echo "🔍 Verificando workflows antes del commit..."
    
    local has_changes=false
    for env in "development" "staging" "production"; do
        if ! ./n8n-ops check --env "$env" --quiet 2>/dev/null; then
            echo "⚠️  $env tiene workflows modificados sin sincronizar"
            has_changes=true
        fi
    done
    
    if [ "$has_changes" = true ]; then
        echo ""
        echo "❌ Commit bloqueado: hay workflows modificados en n8n sin sincronizar"
        echo "💡 Ejecuta el sync primero:"
        echo "   ./n8n-ops sync --env [environment]"
        echo "   git add workflows/ && git commit -m 'sync: update workflows'"
        exit 1
    else
        echo "✅ Todos los workflows están sincronizados"
        exit 0
    fi
}

# Mostrar ayuda
show_help() {
    echo "Uso: $0 [comando] [opciones]"
    echo ""
    echo "Comandos:"
    echo "  check [env]              Verificar ambiente específico (development, staging, production)"
    echo "  all                      Verificar todos los ambientes"
    echo "  monitor [env] [interval] Monitoreo continuo (default: development, 300s)"
    echo "  pre-commit              Verificación para pre-commit hook"
    echo "  --help                  Mostrar esta ayuda"
    echo ""
    echo "Ejemplos:"
    echo "  $0 check production     # Verificar solo producción"
    echo "  $0 all                  # Verificar todos los ambientes"
    echo "  $0 monitor staging 60   # Monitorear staging cada minuto"
    echo "  $0 pre-commit          # Verificar antes de commit"
}

# Lógica principal
case "${1:-check}" in
    "check")
        if [[ -n "$2" ]]; then
            check_environment "$2"
        else
            check_environment "development"
        fi
        ;;
    "all")
        check_all_environments
        ;;
    "monitor")
        monitor_changes "$2" "$3"
        ;;
    "pre-commit")
        pre_commit_check
        ;;
    "--help"|"-h")
        show_help
        ;;
    *)
        echo "❌ Comando desconocido: $1"
        show_help
        exit 1
        ;;
esac