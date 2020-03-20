#!/bin/bash

set -e

kubectl config use-context docker-desktop || exit 1

echo "deploying web..."
kubectl -n handson create -f k8s/web

echo "deploying user-service..."
kubectl -n handson create -f k8s/user-service
