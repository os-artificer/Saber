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
Selector labels for databus
*/}}
{{- define "saber.databus.selectorLabels" -}}
app.kubernetes.io/name: saber
app.kubernetes.io/component: databus
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Selector labels for agent
*/}}
{{- define "saber.agent.selectorLabels" -}}
app.kubernetes.io/name: saber
app.kubernetes.io/component: agent
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Controller full name
*/}}
{{- define "saber.controller.fullname" -}}
{{ .Release.Name }}-controller
{{- end -}}

{{/*
Databus full name
*/}}
{{- define "saber.databus.fullname" -}}
{{ .Release.Name }}-databus
{{- end -}}

{{/*
Agent full name
*/}}
{{- define "saber.agent.fullname" -}}
{{ .Release.Name }}-agent
{{- end -}}

{{/*
Image
*/}}
{{- define "saber.image" -}}
{{ .Values.image.repository }}:{{ .Values.image.tag | default .Chart.AppVersion }}
{{- end -}}
