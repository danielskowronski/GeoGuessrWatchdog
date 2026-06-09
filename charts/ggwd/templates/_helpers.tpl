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

{{- define "ggwd.workerImage" -}}
{{- printf "%s:%s" .Values.worker.image.repository (default .Chart.AppVersion .Values.worker.image.tag) -}}
{{- end -}}

{{- define "ggwd.image" -}}
{{- include "ggwd.workerImage" . -}}
{{- end -}}

{{- define "ggwd.apiImage" -}}
{{- $repository := default .Values.worker.image.repository .Values.api.image.repository -}}
{{- $tag := default (default .Chart.AppVersion .Values.worker.image.tag) .Values.api.image.tag -}}
{{- printf "%s:%s" $repository $tag -}}
{{- end -}}

{{- define "ggwd.temporalFullname" -}}
{{- default (printf "%s-temporal" .Release.Name) .Values.temporal.fullnameOverride | trunc 63 | trimSuffix "-" -}}
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
{{- $cacheKey := printf "%s/%s" $secretName $key -}}
{{- $cache := get $root.Values "_ggwdSecretValues" -}}
{{- if not $cache -}}
{{- $cache = dict -}}
{{- $_ := set $root.Values "_ggwdSecretValues" $cache -}}
{{- end -}}
{{- if hasKey $cache $cacheKey -}}
{{- get $cache $cacheKey -}}
{{- else -}}
{{- $existing := lookup "v1" "Secret" $root.Release.Namespace $secretName -}}
{{- if and $existing (hasKey $existing.data $key) -}}
{{- $value := index $existing.data $key | b64dec -}}
{{- $_ := set $cache $cacheKey $value -}}
{{- $value -}}
{{- else -}}
{{- $value := randAlphaNum 32 -}}
{{- $_ := set $cache $cacheKey $value -}}
{{- $value -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
