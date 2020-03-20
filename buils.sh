#!/bin/bash

echo "bulding web"
docker build -t hands-on-k8s-web:latest -f web/Dockerfile web

echo "bulding user-service"
docker build -t hands-on-k8s-user-service:latest -f user-service/Dockerfile user-service