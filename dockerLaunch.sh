#!/bin/bash

#build imagen dokkerfile
docker build -t block-p .

#para el contenedor por si acaso lo tenias lanzado
docker stop block-p-container

#lo borra por si lo tenias lanzado
docker rm block-p-container

#lo vuelve a crear
docker run --name block-p-container -d -p 8081:8081 -p 8082:8082 -p 8080:8080 block-p