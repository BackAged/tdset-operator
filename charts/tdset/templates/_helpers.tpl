{{/*
Expand the name of the chart.
*/}}
{{- define "tdset.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "tdset.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "tdset.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "tdset.labels" -}}
helm.sh/chart: {{ include "tdset.chart" . }}
{{ include "tdset.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "tdset.selectorLabels" -}}
app.kubernetes.io/name: {{ include "tdset.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "tdset.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "tdset.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the name of the tdset controller role to use
*/}}
{{- define "tdset.controllerRoleName" -}}
{{- if .Values.rbac.create }}
{{- default (include "tdset.fullname" .) .Values.rbac.controller.name }}
{{- else }}
{{- default "default" .Values.rbac.controller.name }}
{{- end }}
{{- end }}

{{/*
Create the name of the tdset leader role to use
*/}}
{{- define "tdset.leaderRoleName" -}}
{{- if .Values.rbac.create }}
{{- default (printf "%s-leader" (include "tdset.fullname" .)) .Values.rbac.leader.name }}
{{- else }}
{{- default "default" .Values.rbac.leader.name }}
{{- end }}
{{- end }}
