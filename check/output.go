package check

import (
	"encoding/json"
	"html/template"
	"io"
)

var (
	output *template.Template
)

const (
	outputRaw = `Healthz Status: {{.Status}}
{{ range $name, $data := .Components -}}
Component: {{ $name }}
 - Healthy: {{ if .Errors }}❌{{ else }}✔{{ end }}
 - Metadata:
{{ range $mkey, $mval :=  .Metadata }}    - {{ $mkey }}: {{ $mval }}
{{ end -}}
{{ if .Errors }} - Errors:
{{ range .Errors }}   - {{ .Description }}
{{ end }}
{{ end }}
{{- end }}`
)

func init() {
	output = template.Must(template.New("output").Parse(outputRaw))
}

type Output struct {
	Format string
	Dest   io.Writer
}

func (o *Output) write(r *response) {
	if o.Dest == nil {
		return
	}

	switch o.Format {
	case "json":
		b, _ := json.MarshalIndent(r, "", "  ")
		o.Dest.Write(b)
	default:
		output.Execute(o.Dest, r)
	}
}
