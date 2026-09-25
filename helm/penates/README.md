# Penates Helm Chart

A Helm chart for deploying the Penates inventory management system on Kubernetes.

## Prerequisites

- Kubernetes cluster (v1.19+)
- Helm 3+ installed
- kubectl configured to access your cluster
- Container registry with Penates backend and frontend images

## Quick Start

### 1. Add the Helm repository (if published)

```bash
helm repo add penates https://leander-wendt.github.io/Penates
helm repo update
```

### 2. Install from local chart

```bash
cd helm/penates
helm install penates . --namespace penates --create-namespace
```

### 3. Install with custom values

```bash
helm install penates . --namespace penates --create-namespace -f values-production.yaml
```

## Configuration

### Required Configuration

Before deploying to production, you **MUST** update the following values:

1. **JWT Secret** - Set a strong, random JWT secret:
   ```yaml
   backend:
     env:
       JWT_SECRET: "your-strong-random-secret-here"
   ```

2. **PostgreSQL Password** - Set strong database credentials:
   ```yaml
   postgres:
     auth:
       password: "your-strong-db-password-here"
   ```

3. **Admin Credentials** - Change the default admin credentials:
   ```yaml
   backend:
     env:
       SEED_ADMIN_EMAIL: "admin@yourdomain.com"
       SEED_ADMIN_PASSWORD: "strong-admin-password"
       SEED_ADMIN_ORG: "Your Organisation"
   ```

### Common Configuration Options

#### Ingress Configuration

```yaml
ingress:
  enabled: true
  className: "nginx"  # or "traefik", "istio", etc.
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
  hosts:
    - host: penates.yourdomain.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: penates-tls
      hosts:
        - penates.yourdomain.com
```

#### External PostgreSQL

To use an external PostgreSQL database:

```yaml
postgres:
  enabled: false
  external: true
  externalHost: "your-postgres-host"
  externalPort: 5432
  externalUsername: "penates-user"
  externalPassword: "your-password"
  externalDatabase: "penates"
```

#### Resource Limits

```yaml
backend:
  resources:
    requests:
      memory: "256Mi"
      cpu: "250m"
    limits:
      memory: "512Mi"
      cpu: "500m"

frontend:
  resources:
    requests:
      memory: "128Mi"
      cpu: "100m"
    limits:
      memory: "256Mi"
      cpu: "250m"

postgres:
  resources:
    requests:
      memory: "512Mi"
      cpu: "250m"
    limits:
      memory: "1Gi"
      cpu: "500m"
```

#### Autoscaling

```yaml
autoscaling:
  enabled: true
  backend:
    enabled: true
    minReplicas: 2
    maxReplicas: 5
    targetCPUUtilizationPercentage: 80
  frontend:
    enabled: true
    minReplicas: 2
    maxReplicas: 5
    targetCPUUtilizationPercentage: 70
```

#### Using Existing Secrets

```yaml
postgres:
  auth:
    existingSecret: "my-postgres-secret"
    secretKeys:
      adminPasswordKey: "postgres-password"
```

### Custom Image Tags

```yaml
image:
  backend:
    repository: ghcr.io/your-org/penates-backend
    tag: "v1.0.0"
    pullPolicy: IfNotPresent
  frontend:
    repository: ghcr.io/your-org/penates-frontend
    tag: "v1.0.0"
    pullPolicy: IfNotPresent
```

### Private Container Registry

```yaml
imagePullSecrets:
  - name: regcred
```

## Uninstalling

```bash
# Remove the release
helm uninstall penates --namespace penates

# Optionally, remove the namespace
kubectl delete namespace penates

# To also remove PVCs (all data will be lost!)
helm uninstall penates --namespace penates --preserve-pvc=false
```

## Upgrading

```bash
# Get the current values
helm get values penates --namespace penates > current-values.yaml

# Edit current-values.yaml as needed
nano current-values.yaml

# Upgrade with new values
helm upgrade penates . --namespace penates -f current-values.yaml
```

## Development

### Building Images

Before deploying, ensure your images are built and pushed to a registry:

```bash
# Build backend
cd backend
docker build -t ghcr.io/your-org/penates-backend:latest .
docker push ghcr.io/your-org/penates-backend:latest

# Build frontend
cd frontend
docker build -t ghcr.io/your-org/penates-frontend:latest .
docker push ghcr.io/your-org/penates-frontend:latest
```

### Linting and Testing the Chart

```bash
# Install helm plugins
helm plugin install https://github.com/quintush/helm-unittest

# Lint the chart
helm lint .

# Test install (dry-run)
helm install penates . --dry-run --debug
```

## Architecture

The Helm chart deploys the following components:

1. **PostgreSQL** (StatefulSet) - Database for Penates
2. **Backend** (Deployment) - Go REST API
3. **Frontend** (Deployment) - Vue.js SPA with nginx
4. **Services** - ClusterIP services for each component
5. **Ingress** (optional) - External access to the frontend

### Network Flow

```
Client → Ingress → Frontend Service → Frontend Pod (nginx) → Backend Service → Backend Pod → PostgreSQL
```

## Values Reference

See `values.yaml` for all configurable options and their defaults.

## Security

The chart follows security best practices:

- Runs containers as non-root users
- Drops all capabilities
- Uses read-only root filesystems
- Implements pod security contexts
- Supports network policies
- Supports pod disruption budgets

## License

MIT License - see the main [Penates LICENSE](LICENSE) file.
