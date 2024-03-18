#!/bin/bash

docker stop block-p-container

docker rm block-p-container

docker run --name block-p-container -d -p 8081:8081 -p 8082:8082 -p 8080:8080 block-p