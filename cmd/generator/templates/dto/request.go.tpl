package request

type Create{{.PascalName}} struct {
	{{- range .Fields }}
	{{ .PascalName }} {{ .Type }} `json:"{{ .SnakeCaseName }}" validate:"{{ .Validation }}"`
	{{- end }}
}

type Update{{.PascalName}} struct {
	ID uint `json:"-"`
	{{- range .Fields }}
	{{ .PascalName }} {{ .Type }} `json:"{{ .SnakeCaseName }}" validate:"{{ .Validation }}"`
	{{- end }}
}

type Query{{.PascalName}} struct {
	Page   int   `query:"page"`
	Limit  int   `query:"limit"`
	Search string `query:"search"`
}