{{- define "stack.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "stack.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name (include "stack.name" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "ggwd.image" -}}
{{- printf "%s:%s" .Values.worker.image.repository (default .Chart.AppVersion .Values.worker.image.tag) -}}
{{- end -}}

{{- define "stack.labels" -}}
app.kubernetes.io/name: {{ include "stack.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
{{- end -}}

{{- define "stack.selectorLabels" -}}
app.kubernetes.io/name: {{ include "stack.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "stack.secretValue" -}}
{{- $secretName := index . 0 -}}
{{- $key := index . 1 -}}
{{- $provided := index . 2 -}}
{{- $root := index . 3 -}}
{{- if $provided -}}
{{- $provided -}}
{{- else -}}
{{- $existing := lookup "v1" "Secret" $root.Release.Namespace $secretName -}}
{{- if and $existing (hasKey $existing.data $key) -}}
{{- index $existing.data $key | b64dec -}}
{{- else -}}
{{- randAlphaNum 32 -}}
{{- end -}}
{{- end -}}
{{- end -}}
