{{/*
Common labels
*/}}
{{- define "saber.labels" -}}
app.kubernetes.io/name: saber
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Selector labels for controller
*/}}
{{- define "saber.controller.selectorLabels" -}}
app.kubernetes.io/name: saber
app.kubernetes.io/component: controller
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Selector labels for transfer
*/}}
{{- define "saber.transfer.selectorLabels" -}}
app.kubernetes.io/name: saber
app.kubernetes.io/component: transfer
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Selector labels for probe
*/}}
{{- define "saber.probe.selectorLabels" -}}
app.kubernetes.io/name: saber
app.kubernetes.io/component: probe
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Controller full name
*/}}
{{- define "saber.controller.fullname" -}}
{{ .Release.Name }}-controller
{{- end -}}

{{/*
Transfer full name
*/}}
{{- define "saber.transfer.fullname" -}}
{{ .Release.Name }}-transfer
{{- end -}}

{{/*
Probe full name
*/}}
{{- define "saber.probe.fullname" -}}
{{ .Release.Name }}-probe
{{- end -}}

{{/*
Image
*/}}
{{- define "saber.image" -}}
{{ .Values.image.repository }}:{{ .Values.image.tag | default .Chart.AppVersion }}
{{- end -}}
