#!/bin/bash

kubectl config use-context docker-desktop || exit 1

# kubectl -n handson delete -f k8s/shared/jaeger/crd.yaml
# kubectl -n handson delete -f k8s/shared/jaeger/rbac.yaml
# kubectl -n handson delete -f k8s/shared/jaeger/operator.yaml
# kubectl -n handson delete -f k8s/shared/jaeger_cr.yaml

kubectl -n handson delete -f k8s/web
kubectl -n handson delete -f k8s/user-service

# kubectl delete -f k8s/shared/nginx-ingress/mandatory.yaml
# kubectl delete -f k8s/shared/nginx-ingress/service.yaml

# kubectl delete --ignore-not-found=true -f k8s/shared/kube-prometheus/manifests/ -f k8s/shared/kube-prometheus/manifests/setup
