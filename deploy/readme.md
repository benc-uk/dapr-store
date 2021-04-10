# Deployment of Dapr Store to Kubernetes

This is a brief guide to deploying Dapr Store to Kubernetes.

Assumptions:

- kubectl is installed, and configured to access your Kubernetes cluster
- dapr CLI is installed - https://docs.dapr.io/getting-started/install-dapr-cli/
- helm is installed - https://helm.sh/docs/intro/install/

This guide does not cover more advanced deployment scenarios such as deploying behind a DNS name, or with HTTPS enabled or with used identity enabled.

For more details see the [documentation for the Dapr Store Helm chart](./helm/daprstore/readme.md)

## 🥾 Initial Setup

### Deploy Dapr to Kubernetes

Skip this if the Dapr control plane is already deployed

```bash
dapr init --kubernetes
kubectl get pod --namespace dapr-system
```

Full instructions here:  
📃 https://docs.dapr.io/operations/hosting/kubernetes/kubernetes-overview/

Optional - If you wish to view or check the Dapr dashboard

```bash
kubectl port-forward deploy/dapr-dashboard --namespace dapr-system 8080:8080
```

Open the dashboard at http://localhost:8080/

### Create namespace for Dapr Store app

```bash
namespace=dapr-store
kubectl create namespace $namespace
```

### Add Helm repos

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
```

## 💾 Deploy Redis

```bash
helm install dapr-redis bitnami/redis --values deploy/config/redis-values.yaml --namespace $namespace
```

Validate & check status

```bash
helm list --namespace $namespace
kubectl get pod daprstore-redis-master-0 --namespace $namespace
```

## 🌐 Deploy NGINX Ingress Controller (API Gateway)

```bash
helm install api-gateway ingress-nginx/ingress-nginx --values deploy/config/ingress-values.yaml --namespace $namespace
```

Validate & check status

```bash
helm list --namespace $namespace
kubectl get pod -l app.kubernetes.io/instance=api-gateway --namespace $namespace
kubectl get svc --namespace $namespace
```

## 🚀 Deploy Dapr Store

Now deploy the Dapr Store application and all services using Helm

```bash
helm install store ./deploy/helm/daprstore --namespace $namespace
```

Validate & check status

```bash
helm list --namespace $namespace
kubectl get pod -l app.kubernetes.io/instance=store --namespace $namespace
```

To get the URL of the deployed store run the following command:

```bash
echo -e "Access Dapr Store here: http://$(kubectl get svc -l "purpose=daprstore-api-gateway" -o jsonpath="{.items[0].status.loadBalancer.ingress[0].ip}")/"
```
