package response

import "time"

type {{.PascalName}}Response struct {
	ID        uint      `json:"id"`
	{{- range .Fields }}
	{{ .PascalName }} {{ .Type }} `json:"{{ .SnakeCaseName }}"`
	{{- end }}
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}