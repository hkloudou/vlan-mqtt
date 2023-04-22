{{/*
Expand the name of the chart.
*/}}
{{- define "basic.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "basic.fullname" -}}
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
{{- define "basic.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "basic.labels" -}}
helm.sh/chart: {{ include "basic.chart" . }}
{{ include "basic.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "basic.selectorLabels" -}}
app.kubernetes.io/name: {{ include "basic.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "basic.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "basic.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}


{{- define "basic.volumes" -}}
{{- if .Values.secretFile.enable}}
- name: auto-secret-files
  secret:
    secretName: {{ printf "%s-secret-files" (include "basic.fullname" .) }}
{{- end -}}
{{- end -}}


{{- define "basic.volumeMounts" -}}
{{- if .Values.secretFile.enable}}
- name: auto-secret-files
  mountPath: {{ .Values.secretFile.mounted }}
{{- end -}}
{{- end -}}

{{- define "basic.envs"}}
{{- range $key,$val := .Values.envs }}
- name: {{$key}}
  value: {{$val | quote}}
{{- end }}
{{- end }}