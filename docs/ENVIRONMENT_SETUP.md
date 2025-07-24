# Configuración de Variables de Entorno

## Problema Común: Variables de Entorno No Encontradas

Si ves errores como:
```
API key not found in environment variable N8N_DEVELOPMENT_API_KEY
```

Esto significa que las variables de entorno del archivo `.env` no se han cargado en tu shell.

## Solución Rápida

Antes de ejecutar cualquier comando de n8n-ops, carga las variables de entorno:

```bash
# Cargar variables del archivo .env
source .env

# Verificar que se cargaron correctamente
echo $N8N_DEV_API_KEY
echo $N8N_DEVELOPMENT_API_KEY

# Para el comando sync, necesitas exportar variables adicionales
export N8N_URL=$N8N_DEV_URL
export N8N_API_KEY=$N8N_DEV_API_KEY

# Ahora ejecutar n8n-ops
./n8n-ops onboard
./n8n-ops sync --env development
```

## Configuración Permanente

Para evitar tener que ejecutar `source .env` cada vez, puedes añadir las variables a tu perfil de shell:

### Para Bash (~/.bashrc)
```bash
echo 'export N8N_DEV_API_KEY="tu_api_key_aqui"' >> ~/.bashrc
echo 'export N8N_DEV_URL="http://localhost:5678"' >> ~/.bashrc
source ~/.bashrc
```

### Para Zsh (~/.zshrc)
```bash
echo 'export N8N_DEV_API_KEY="tu_api_key_aqui"' >> ~/.zshrc
echo 'export N8N_DEV_URL="http://localhost:5678"' >> ~/.zshrc
source ~/.zshrc
```

## Verificación de Variables

Para verificar que las variables están configuradas correctamente:

```bash
# Verificar variables de API
echo "DEV API KEY: $N8N_DEV_API_KEY"
echo "DEV URL: $N8N_DEV_URL"

# Verificar todas las variables de n8n
env | grep N8N_
```

## Variables Requeridas

Para el funcionamiento básico necesitas:

```bash
# Desarrollo
export N8N_DEV_API_KEY="tu_api_key_de_desarrollo"
export N8N_DEV_URL="http://localhost:5678"

# Staging (opcional)
export N8N_STAGING_API_KEY="tu_api_key_de_staging"
export N8N_STAGING_URL="https://n8n-staging.example.com"

# Producción (opcional)
export N8N_PROD_API_KEY="tu_api_key_de_produccion"
export N8N_PROD_URL="https://n8n-prod.example.com"
```

## Compatibilidad de Nombres

n8n-ops es compatible con múltiples formatos de nombres de variables:

| Formato Corto | Formato Largo | Uso |
|---------------|---------------|-----|
| `N8N_DEV_API_KEY` | `N8N_DEVELOPMENT_API_KEY` | Desarrollo |
| `N8N_STAGING_API_KEY` | `N8N_STAGING_API_KEY` | Staging |
| `N8N_PROD_API_KEY` | `N8N_PRODUCTION_API_KEY` | Producción |

## Variables para pruebas

Define `MOCK_SERVER_TIMEOUT` para ajustar el tiempo de espera al iniciar el
servidor `mock` utilizado en las pruebas. El valor debe seguir el formato de
duración de Go (`5s`, `1m`, ...). Si no se especifica, se usa `5s`.

## Solución de Problemas

### Error: "API key not found"
1. Verifica que el archivo `.env` existe en el directorio actual
2. Ejecuta `source .env` antes de usar n8n-ops
3. Verifica que la variable existe: `echo $N8N_DEV_API_KEY`

### Error: "Connection failed"
1. Verifica que n8n está corriendo: `curl http://localhost:5678/`
2. Verifica que la URL es correcta en tu `.env`
3. Prueba la API manualmente:
   ```bash
   curl -H "X-N8N-API-KEY: $N8N_DEV_API_KEY" http://localhost:5678/api/v1/workflows
   ```

### Error: "Invalid API key"
1. Verifica que el API key es correcto en n8n
2. Asegúrate de que no hay espacios extra en el `.env`
3. Regenera el API key en n8n si es necesario

## Ejemplo Completo

```bash
# 1. Crear archivo .env
cat > .env << 'EOF'
N8N_DEV_API_KEY="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
N8N_DEV_URL="http://localhost:5678"
EOF

# 2. Cargar variables
source .env

# 3. Verificar
echo $N8N_DEV_API_KEY

# 4. Usar n8n-ops
./n8n-ops onboard
```

## Automatización

Para automatizar la carga de variables, puedes crear un script:

```bash
#!/bin/bash
# run-n8n-ops.sh

# Cargar variables de entorno
if [ -f .env ]; then
    source .env
    echo "✅ Variables de entorno cargadas desde .env"
else
    echo "❌ Archivo .env no encontrado"
    exit 1
fi

# Ejecutar n8n-ops con los argumentos pasados
./n8n-ops "$@"
```

Uso:
```bash
chmod +x run-n8n-ops.sh
./run-n8n-ops.sh onboard
./run-n8n-ops.sh sync --env development
```