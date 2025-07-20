{{/*
Expand the name of the chart.
*/}}
{{- define "landing.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "landing.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "landing.labels" -}}
helm.sh/chart: {{ include "landing.chart" . }}
{{ include "landing.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "landing.selectorLabels" -}}
app.kubernetes.io/name: {{ include "landing.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
