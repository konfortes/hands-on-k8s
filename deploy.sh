#!/bin/bash


kubectl config use-context docker-desktop || exit 1

echo "deploying jaeger"
kubectl create namespace observability
set -e
kubectl -n handson create -f k8s/shared/jaeger/crd.yaml
kubectl -n handson create -f k8s/shared/jaeger/rbac.yaml
kubectl -n handson create -f k8s/shared/jaeger/operator.yaml
kubectl -n handson create -f k8s/shared/jaeger_cr.yaml


echo "deploying nginx ingress controller"
kubectl apply -f k8s/shared/nginx-ingress/mandatory.yaml
kubectl apply -f k8s/shared/nginx-ingress/service.yaml

kubectl create -f k8s/shared/kube-prometheus/manifests/setup
until kubectl get servicemonitors --all-namespaces ; do date; sleep 1; echo ""; done
kubectl create -f k8s/shared/kube-prometheus/manifests/

echo "deploying web..."
kubectl -n handson create -f k8s/web

echo "deploying user-service..."
kubectl -n handson create -f k8s/user-service