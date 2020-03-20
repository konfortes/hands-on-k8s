#!/bin/bash

set -e

kubectl config use-context docker-desktop || exit 1

echo "deploying web..."
kubectl -n handson create -f k8s/web

echo "deploying user-service..."
kubectl -n handson create -f k8s/user-service

echo "deploying nginx ingress controller"
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/mandatory.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/provider/cloud-generic.yaml

echo "deploying ingress"
kubectl -n handson apply -f k8s/shared/ingress.yaml

