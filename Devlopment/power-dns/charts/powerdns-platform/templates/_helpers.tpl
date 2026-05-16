{{- define "powerdns-platform.name" -}}
powerdns-platform
{{- end -}}

{{- define "powerdns-platform.fullname" -}}
{{- printf "%s-%s" .Release.Name (include "powerdns-platform.name" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "powerdns-platform.image" -}}
{{- if .Values.global.imageRegistry -}}
{{- printf "%s/%s" .Values.global.imageRegistry .Values.operator.image -}}
{{- else -}}
{{- .Values.operator.image -}}
{{- end -}}
{{- end -}}