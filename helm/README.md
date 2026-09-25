# Penates Helm Charts

This directory contains Helm charts for deploying Penates on Kubernetes.

## Available Charts

- [penates](./penates/) - Main Penates application chart

## Quick Start

### Install from source

```bash
cd helm/penates
helm install penates . --namespace penates --create-namespace
```

### Install with custom configuration

```bash
helm install penates ./penates --namespace penates --create-namespace -f ./penates/values-prod.yaml
```

## Chart Features

The Penates Helm chart includes:

- **PostgreSQL** database with persistent storage
- **Backend** service (Go REST API)
- **Frontend** service (Vue.js SPA with nginx)
- Optional **Ingress** configuration with TLS support
- **Horizontal Pod Autoscaling** support
- **Monitoring** support (ServiceMonitor for Prometheus)
- **Network Policies** for enhanced security
- **Pod Disruption Budget** for high availability
- **RBAC** support with service accounts
- Multiple environment configurations (dev, staging, prod)

## Prerequisites

- Kubernetes cluster (v1.19+)
- Helm 3+
- kubectl
- Container registry with Penates images

## Configuration

See the [penates chart README](./penates/README.md) for detailed configuration options.

## Development

For chart development, see the [Makefile](./penates/Makefile) for common operations.

```bash
cd helm/penates
make lint        # Lint the chart
make test        # Test install (dry-run)
make template    # Render templates
make package     # Package the chart
```

## Directory Structure

```
helm/
└── penates/
    ├── Chart.yaml              # Chart metadata
    ├── values.yaml            # Default values
    ├── values-dev.yaml        # Development values
    ├── values-staging.yaml    # Staging values
    ├── values-prod.yaml       # Production values
    ├── values.schema.json     # JSON schema for values validation
    ├── README.md              # Chart documentation
    ├── Makefile              # Common operations
    ├── .helmignore            # Files to ignore when packaging
    └── templates/
        ├── _helpers.tpl       # Template helpers
        ├── serviceaccount.yaml
        ├── backend-deployment.yaml
        ├── backend-service.yaml
        ├── backend-pvc.yaml
        ├── backend-hpa.yaml
        ├── frontend-deployment.yaml
        ├── frontend-service.yaml
        ├── frontend-hpa.yaml
        ├── postgres-statefulset.yaml
        ├── postgres-service.yaml
        ├── postgres-secret.yaml
        ├── postgres-configmap.yaml
        ├── ingress.yaml
        ├── networkpolicy.yaml
        ├── pdb.yaml
        ├── servicemonitor.yaml
        └── NOTES.txt
```
