#!/bin/bash

# Cargar variables de entorno
source .env

# Exportar variables para el comando sync
export N8N_URL=$N8N_DEV_URL
export N8N_API_KEY=$N8N_DEV_API_KEY

echo "Testing n8n connection..."
echo "N8N_URL: $N8N_URL"
echo "N8N_API_KEY: ${N8N_API_KEY:0:10}..."

# Probar la conexión a la API
echo -e "\nTesting API connection..."
curl -s -H "X-N8N-API-KEY: $N8N_API_KEY" $N8N_URL/api/v1/workflows | jq .

# Ejecutar el comando sync con verbose
echo -e "\nRunning sync command..."
./n8n-ops sync --env development --verbose