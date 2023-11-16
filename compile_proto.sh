#!/bin/bash

# Este es un script para compilar el archivo proto con protoc

# Verifica si el usuario tiene permisos de sudo
if [ "$EUID" -ne 0 ]; then
  echo "Por favor, ejecuta este script con sudo."
  exit 1
fi

# Ruta al archivo proto
PROTO_FILE="proto/monitoring.proto"

# Ejecuta protoc con las opciones necesarias
sudo protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    $PROTO_FILE

# Informa al usuario que la compilación se ha completado
echo "Compilación exitosa. Los archivos generados se encuentran junto a $PROTO_FILE."

