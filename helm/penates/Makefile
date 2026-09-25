.PHONY: all install upgrade uninstall lint test package publish clean

NAME ?= penates
NAMESPACE ?= penates
RELEASE_NAME ?= penates
CHART_DIR ?= .
VALUES_FILE ?= 

# Include parent Makefile variables if available
-include ../../../Makefile

all: lint

# Install the chart
install: 
	@echo "Installing ${RELEASE_NAME} chart..."
	@helm install ${RELEASE_NAME} ${CHART_DIR} \
		--namespace ${NAMESPACE} \
		--create-namespace \
		${VALUES_FILE:+-f ${VALUES_FILE}}

# Upgrade the chart
upgrade: 
	@echo "Upgrading ${RELEASE_NAME} chart..."
	@helm upgrade ${RELEASE_NAME} ${CHART_DIR} \
		--namespace ${NAMESPACE} \
		${VALUES_FILE:+-f ${VALUES_FILE}}

# Uninstall the chart
uninstall: 
	@echo "Uninstalling ${RELEASE_NAME} chart..."
	@helm uninstall ${RELEASE_NAME} --namespace ${NAMESPACE}

# Lint the chart
lint: 
	@echo "Linting chart..."
	@helm lint ${CHART_DIR}

# Test install (dry-run)
test: 
	@echo "Testing chart installation (dry-run)..."
	@helm install ${RELEASE_NAME} ${CHART_DIR} \
		--namespace ${NAMESPACE} \
		--dry-run \
		--debug \
		${VALUES_FILE:+-f ${VALUES_FILE}}

# Template the chart (render without installing)
template: 
	@echo "Rendering chart templates..."
	@helm template ${RELEASE_NAME} ${CHART_DIR} \
		--namespace ${NAMESPACE} \
		${VALUES_FILE:+-f ${VALUES_FILE}}

# Package the chart
package: 
	@echo "Packaging chart..."
	@helm package ${CHART_DIR}

# Publish the chart (requires pages workflow or helm push)
publish: package
	@echo "Chart packaged. To publish:"
	@echo "  - For GitHub Pages: Run helm repo index and commit to gh-pages branch"
	@echo "  - For OCI: helm push ${RELEASE_NAME}-*.tgz oci://your-registry"

# Clean up
clean: 
	@echo "Removing chart package files..."
	@rm -f ${RELEASE_NAME}-*.tgz

# Get values from running release
get-values: 
	@echo "Getting values from running release..."
	@helm get values ${RELEASE_NAME} --namespace ${NAMESPACE}

# Get manifest from running release
get-manifest: 
	@echo "Getting manifest from running release..."
	@helm get manifest ${RELEASE_NAME} --namespace ${NAMESPACE}

# Get all secrets
get-secrets: 
	@echo "Getting secrets from running release..."
	@helm get secrets ${RELEASE_NAME} --namespace ${NAMESPACE}

# List releases
list: 
	@echo "Listing releases..."
	@helm list --namespace ${NAMESPACE}

# Status of release
status: 
	@echo "Getting status of release..."
	@helm status ${RELEASE_NAME} --namespace ${NAMESPACE}

# Rollback to previous revision
rollback: 
	@echo "Rolling back to previous revision..."
	@helm rollback ${RELEASE_NAME} 1 --namespace ${NAMESPACE}

# Help
help: 
	@echo "Available targets:"
	@echo "  install     - Install the chart"
	@echo "  upgrade     - Upgrade the chart"
	@echo "  uninstall   - Uninstall the chart"
	@echo "  lint        - Lint the chart"
	@echo "  test        - Test install (dry-run)"
	@echo "  template    - Render chart templates"
	@echo "  package     - Package the chart"
	@echo "  publish     - Package and show publish instructions"
	@echo "  clean       - Remove package files"
	@echo "  get-values  - Get values from running release"
	@echo "  get-manifest - Get manifest from running release"
	@echo "  get-secrets - Get secrets from running release"
	@echo "  list        - List releases"
	@echo "  status      - Get status of release"
	@echo "  rollback    - Rollback to previous revision"
	@echo ""
	@echo "Example usage:"
	@echo "  make install NAMESPACE=production VALUES_FILE=values-prod.yaml"
	@echo "  make upgrade RELEASE_NAME=penates NAMESPACE=staging"
