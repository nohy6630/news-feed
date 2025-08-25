{{- define "news-feed.labels" -}}
app.kubernetes.io/name: {{ .Values.common.name }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Values.common.version }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}
