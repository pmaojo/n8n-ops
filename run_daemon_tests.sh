#!/bin/bash

echo "🚀 EJECUTANDO TESTS COMPREHENSIVOS DEL DAEMON N8N-OPS"
echo "======================================================"
echo ""

# Colores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m' 
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Verificar que los servicios necesarios estén corriendo
echo -e "${BLUE}🔍 Verificando servicios necesarios...${NC}"

# Verificar mock n8n server
if curl -s http://localhost:3001/health > /dev/null; then
    echo -e "${GREEN}✅ Mock n8n server está corriendo${NC}"
else
    echo -e "${RED}❌ Mock n8n server NO está corriendo en puerto 3001${NC}"
    echo "Inicia el servidor con: go run mock-n8n-server/main.go"
    exit 1
fi

echo ""
echo -e "${YELLOW}🧪 Ejecutando tests del daemon...${NC}"
echo ""

# Ejecutar tests con verbose output
go test -v ./tests/daemon_test.go -timeout=30s

TEST_EXIT_CODE=$?

echo ""
echo "======================================================"

if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}🎉 ¡TODOS LOS TESTS PASARON!${NC}"
    echo -e "${GREEN}✅ El daemon funciona perfectamente${NC}"
    echo ""
    echo -e "${BLUE}📋 Funcionalidades verificadas:${NC}"
    echo "   • Detección de cambios en archivos JSON"
    echo "   • Conexión con API de n8n"
    echo "   • Creación de backups automáticos"
    echo "   • Manejo de múltiples workflows"
    echo "   • Rendimiento con cambios rápidos"
    echo "   • Estructura correcta de metadatos"
    echo ""
    echo -e "${YELLOW}🤖 Para usar el daemon:${NC}"
    echo "   go run main.go --daemon --demo --env development"
    echo ""
else
    echo -e "${RED}❌ ALGUNOS TESTS FALLARON${NC}"
    echo -e "${RED}Ver detalles arriba para más información${NC}"
    exit 1
fi