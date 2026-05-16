{{- define "powerdns-stack.name" -}}
powerdns-stack
{{- end -}}

{{- define "powerdns-stack.fullname" -}}
{{- printf "%s-%s" .Release.Name (include "powerdns-stack.name" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "powerdns-stack.secretName" -}}
{{- printf "%s-%s" (include "powerdns-stack.fullname" .) .suffix -}}
{{- end -}}