#!/bin/bash
# Script para workflows de equipo seguros

set -e

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🤝 n8n Team Collaboration Workflow${NC}"
echo "=================================="

# Función para mostrar estado del equipo
show_team_status() {
    echo -e "\n${BLUE}👥 Estado del Equipo${NC}"
    echo "-------------------"
    
    for env in "development" "staging" "production"; do
        echo -e "\n📊 ${env}:"
        
        if ./n8n-ops check --env "$env" --quiet 2>/dev/null; then
            echo -e "   ✅ Todos los workflows sincronizados"
        else
            echo -e "   ⚠️  Workflows modificados:"
            ./n8n-ops check --env "$env" --json 2>/dev/null | jq -r '
                .workflows.modified[]? | "     - \(.name) (\(.timeAgo))"
            ' || echo "     Error obteniendo detalles"
        fi
    done
}

# Función para workflow de development seguro
safe_development_workflow() {
    local branch_name="$1"
    
    if [[ -z "$branch_name" ]]; then
        echo -e "${RED}❌ Especifica el nombre de la branch${NC}"
        echo "   Uso: $0 dev-workflow feature/nueva-funcionalidad"
        return 1
    fi
    
    echo -e "\n${GREEN}🚀 Iniciando workflow de desarrollo seguro${NC}"
    echo "Branch: $branch_name"
    
    # 1. Verificar estado actual
    echo -e "\n1️⃣  Verificando estado actual..."
    ./n8n-ops check --env development
    
    if ! ./n8n-ops check --env development --quiet; then
        echo -e "${YELLOW}⚠️  Hay workflows modificados. ¿Sincronizar primero? (y/N)${NC}"
        read -r response
        if [[ "$response" =~ ^[Yy]$ ]]; then
            ./n8n-ops sync --env development
            git add workflows/development/
            git commit -m "sync: update development workflows before feature work"
        fi
    fi
    
    # 2. Crear/cambiar a branch
    echo -e "\n2️⃣  Configurando branch..."
    if git show-ref --verify --quiet refs/heads/"$branch_name"; then
        echo "Branch existente, cambiando..."
        git checkout "$branch_name"
    else
        echo "Creando nueva branch..."
        git checkout -b "$branch_name"
    fi
    
    # 3. Trabajar en n8n
    echo -e "\n3️⃣  ${GREEN}Listo para trabajar en n8n web interface${NC}"
    echo "   👉 Edita tus workflows en: https://n8n-dev.company.com"
    echo "   👉 Cuando termines, presiona ENTER para continuar..."
    read -r
    
    # 4. Sync y commit
    echo -e "\n4️⃣  Sincronizando cambios..."
    ./n8n-ops sync --env development --verbose
    
    if git status --porcelain workflows/development/ | grep -q .; then
        echo -e "\n5️⃣  Cambios detectados, haciendo commit..."
        git add workflows/development/
        echo "Describe tus cambios:"
        read -r commit_message
        git commit -m "$commit_message"
        
        echo -e "\n6️⃣  ¿Push a remote? (Y/n)"
        read -r push_response
        if [[ ! "$push_response" =~ ^[Nn]$ ]]; then
            git push origin "$branch_name"
            echo -e "${GREEN}✅ Cambios pusheados. Puedes crear Merge Request en GitLab${NC}"
        fi
    else
        echo -e "${YELLOW}ℹ️  No se detectaron cambios en workflows${NC}"
    fi
}

# Función para pre-deploy checks
pre_deploy_checks() {
    local env="$1"
    
    echo -e "\n${BLUE}🔍 Verificaciones Pre-Deploy: $env${NC}"
    echo "================================"
    
    # 1. Check sync status
    echo "1️⃣  Verificando estado de sincronización..."
    if ./n8n-ops check --env "$env" --quiet; then
        echo -e "   ✅ Todos los workflows sincronizados"
    else
        echo -e "   ${RED}❌ Hay workflows modificados sin sincronizar${NC}"
        ./n8n-ops check --env "$env"
        echo -e "\n   ${YELLOW}¿Continuar deploy de todas formas? (y/N)${NC}"
        read -r response
        if [[ ! "$response" =~ ^[Yy]$ ]]; then
            return 1
        fi
    fi
    
    # 2. Validate workflows
    echo -e "\n2️⃣  Validando workflows..."
    if ./n8n-ops validate "./workflows/$env/" --verbose; then
        echo -e "   ✅ Validación exitosa"
    else
        echo -e "   ${RED}❌ Errores en validación${NC}"
        return 1
    fi
    
    # 3. Show current status
    echo -e "\n3️⃣  Estado actual del ambiente:"
    ./n8n-ops status --env "$env"
    
    echo -e "\n${GREEN}✅ Pre-deploy checks completados para $env${NC}"
    return 0
}

# Función para rollback de emergencia
emergency_rollback() {
    local env="$1"
    
    echo -e "\n${RED}🚨 ROLLBACK DE EMERGENCIA: $env${NC}"
    echo "========================="
    
    if [[ "$env" == "production" ]]; then
        echo -e "${RED}⚠️  ADVERTENCIA: Rollback en PRODUCCIÓN${NC}"
        echo "   ¿Estás seguro? Escribe 'ROLLBACK' para confirmar:"
        read -r confirmation
        if [[ "$confirmation" != "ROLLBACK" ]]; then
            echo "Rollback cancelado"
            return 1
        fi
    fi
    
    echo "Ejecutando rollback..."
    if ./n8n-ops rollback --env "$env" --verbose; then
        echo -e "${GREEN}✅ Rollback completado${NC}"
        
        # Sync after rollback to update local files
        echo "Sincronizando estado después del rollback..."
        ./n8n-ops sync --env "$env"
        
        # Commit the rollback state
        git add "workflows/$env/"
        git commit -m "rollback: emergency rollback for $env environment"
        
        echo -e "${GREEN}✅ Estado local actualizado después del rollback${NC}"
    else
        echo -e "${RED}❌ Error durante rollback${NC}"
        return 1
    fi
}

# Función para monitoring en tiempo real
live_monitoring() {
    local env="${1:-development}"
    local interval="${2:-30}"
    
    echo -e "\n${BLUE}👀 Monitoring en Tiempo Real: $env${NC}"
    echo "============================="
    echo "Verificando cada $interval segundos (Ctrl+C para salir)"
    echo ""
    
    while true; do
        local timestamp=$(date '+%H:%M:%S')
        
        if ./n8n-ops check --env "$env" --quiet; then
            echo -e "[$timestamp] ${GREEN}✅ $env: Sincronizado${NC}"
        else
            echo -e "[$timestamp] ${YELLOW}⚠️  $env: Cambios detectados${NC}"
            ./n8n-ops check --env "$env" --json 2>/dev/null | jq -r '
                .workflows.modified[]? | "    📝 \(.name) - \(.timeAgo)"
            ' 2>/dev/null || echo "    Error obteniendo detalles"
            
            # Auto-notify si está configurado
            if [[ -n "$SLACK_WEBHOOK" ]]; then
                curl -s -X POST "$SLACK_WEBHOOK" \
                    -d "{\"text\":\"⚠️ Cambios detectados en $env workflows\"}" \
                    >/dev/null || true
            fi
        fi
        
        sleep "$interval"
    done
}

# Mostrar ayuda
show_help() {
    echo "Uso: $0 [comando] [opciones]"
    echo ""
    echo "Comandos para colaboración en equipo:"
    echo "  team-status                    Ver estado de todos los ambientes"
    echo "  dev-workflow <branch>          Workflow completo de desarrollo"
    echo "  pre-deploy <env>              Verificaciones antes de deploy"
    echo "  emergency-rollback <env>      Rollback de emergencia"
    echo "  monitor [env] [interval]      Monitoring en tiempo real"
    echo "  --help                        Mostrar esta ayuda"
    echo ""
    echo "Ejemplos:"
    echo "  $0 team-status"
    echo "  $0 dev-workflow feature/new-payment-flow"
    echo "  $0 pre-deploy staging"
    echo "  $0 emergency-rollback production"
    echo "  $0 monitor development 60"
}

# Lógica principal
case "${1:-team-status}" in
    "team-status")
        show_team_status
        ;;
    "dev-workflow")
        safe_development_workflow "$2"
        ;;
    "pre-deploy")
        if [[ -z "$2" ]]; then
            echo -e "${RED}❌ Especifica el ambiente: development, staging, production${NC}"
            exit 1
        fi
        pre_deploy_checks "$2"
        ;;
    "emergency-rollback")
        if [[ -z "$2" ]]; then
            echo -e "${RED}❌ Especifica el ambiente para rollback${NC}"
            exit 1
        fi
        emergency_rollback "$2"
        ;;
    "monitor")
        live_monitoring "$2" "$3"
        ;;
    "--help"|"-h")
        show_help
        ;;
    *)
        echo -e "${RED}❌ Comando desconocido: $1${NC}"
        show_help
        exit 1
        ;;
esac