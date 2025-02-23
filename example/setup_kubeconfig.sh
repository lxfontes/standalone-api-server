#!/bin/bash
set -xeuo pipefail

export KUBECONFIG=./kubeconfig
CA_FILE=./minica.pem
TOKEN="abc-123"

kubectl config set-cluster standalone --server=https://localhost:6443/ --certificate-authority="$CA_FILE"
kubectl config set-context standalone --cluster=standalone
kubectl config set-credentials standalone-user --token="${TOKEN}"
kubectl config set-context standalone --user=standalone-user
kubectl config use-context standalone
