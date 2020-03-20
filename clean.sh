#!/bin/bash

kubectl config use-context docker-desktop || exit 1

kubectl -n handson delete -f k8s/web
kubectl -n handson delete -f k8s/user-service
kubectl delete -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/mandatory.yaml
kubectl delete -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.30.0/deploy/static/provider/cloud-generic.yaml