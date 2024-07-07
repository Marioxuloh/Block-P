#!/bin/bash

# Obtener el porcentaje de uso de la CPU en macOS utilizando ps
cpu_percentage=$(ps -A -o %cpu | awk '{s+=$1} END {print int(s)}')

# Construir la cadena de salida
output="${cpu_percentage}%"

# Imprimir la cadena de salida
echo "$output"

