#!/bin/bash

set -e

kubectl config use-context docker-desktop || exit 1

kubectl -n handson delete -f k8s/web
kubectl -n handson delete -f k8s/user-service
