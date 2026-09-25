{{/*
Expand the name of the chart.
*/}}
{{- define "penates.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "penates.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "penates.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "penates.labels" -}}
helm.sh/chart: {{ include "penates.chart" . }}
{{ include "penates.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "penates.selectorLabels" -}}
app.kubernetes.io/name: {{ include "penates.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Create a label selector for a specific component
*/}}
{{- define "penates.componentSelectorLabels" -}}
{{- $component := default .Values.nameOverride .Values.nameOverride -}}
app.kubernetes.io/name: {{ include "penates.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: {{ .component | default "" }}
{{- end -}}

{{/*
Create a secret name for PostgreSQL
*/}}
{{- define "penates.postgresSecretName" -}}
{{- if .Values.postgres.auth.existingSecret -}}
{{- .Values.postgres.auth.existingSecret -}}
{{- else -}}
{{- printf "%s-postgres" (include "penates.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Create database URL for backend
*/}}
{{- define "penates.databaseURL" -}}
{{- if .Values.postgres.external -}}
postgres://{{ .Values.postgres.externalUsername }}:{{ .Values.postgres.externalPassword }}@{{ .Values.postgres.externalHost }}:{{ .Values.postgres.externalPort }}/{{ .Values.postgres.externalDatabase | default "penates" }}?sslmode=disable
{{- else -}}
postgres://{{ .Values.postgres.auth.username }}:{{ .Values.postgres.auth.password }}@{{ include "penates.postgresServiceName" . }}:{{ .Values.service.postgres.port }}>{{ include "penates.postgresDatabase" . }}?sslmode=disable
{{- end -}}
{{- end -}}

{{/*
PostgreSQL service name
*/}}
{{- define "penates.postgresServiceName" -}}
{{- printf "%s-postgres" (include "penates.fullname" .) -}}
{{- end -}}

{{/*
PostgreSQL database name
*/}}
{{- define "penates.postgresDatabase" -}}
/{{ .Values.postgres.auth.database }}
{{- end -}}

{{/*
Backend service name
*/}}
{{- define "penates.backendServiceName" -}}
{{- printf "%s-backend" (include "penates.fullname" .) -}}
{{- end -}}

{{/*
Frontend service name
*/}}
{{- define "penates.frontendServiceName" -}}
{{- printf "%s-frontend" (include "penates.fullname" .) -}}
{{- end -}}

{{/*
Create TLS secret name
*/}}
{{- define "penates.tlsSecretName" -}}
{{- if .Values.ingress.tls -}}
{{- index .Values.ingress.tls 0 | default (printf "%s-tls" (include "penates.fullname" .)) -}}
{{- else -}}
{{- printf "%s-tls" (include "penates.fullname" .) -}}
{{- end -}}
{{- end -}}
