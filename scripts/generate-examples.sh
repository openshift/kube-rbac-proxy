#!/usr/bin/env bash

version="$(cat VERSION)"

sed -e "s/KUBE_RBAC_PROXY_VERSION/${version}/g" scripts/templates/non-resource-url-deployment.yaml > examples/non-resource-url/deployment.yaml
sed -e "s/KUBE_RBAC_PROXY_VERSION/${version}/g" scripts/templates/non-resource-url-token-request-deployment.yaml > examples/non-resource-url-token-request/deployment.yaml
sed -e "s/KUBE_RBAC_PROXY_VERSION/${version}/g" scripts/templates/resource-attributes-deployment.yaml > examples/resource-attributes/deployment.yaml
sed -e "s/KUBE_RBAC_PROXY_VERSION/${version}/g" scripts/templates/rewrites-deployment.yaml > examples/rewrites/deployment.yaml
sed -e "s/KUBE_RBAC_PROXY_VERSION/${version}/g" scripts/templates/oidc-deployment.yaml > examples/oidc/deployment.yaml
sed -e "s/KUBE_RBAC_PROXY_VERSION/${version}/g" scripts/templates/static-auth-deployment.yaml > examples/static-auth/deployment.yaml
