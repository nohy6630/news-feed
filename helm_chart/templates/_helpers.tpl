{{/*
Common labels - simplified
*/}}
{{- define "news-feed.labels" -}}
app.kubernetes.io/name: {{ .Values.common.name | quote }}
app.kubernetes.io/instance: {{ .Release.Name | quote }}
app.kubernetes.io/version: {{ .Values.common.version | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service | quote }}
{{- end }}
