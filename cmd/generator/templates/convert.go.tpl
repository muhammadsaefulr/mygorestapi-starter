package convert_types

import (
	"{{.ModulePath}}/internal/domain/dto/{{.Name}}/request"
	model "{{.ModulePath}}/internal/domain/model"
)

func Create{{.PascalName}}ToModel(req *request.Create{{.PascalName}}) *model.{{.PascalName}} {
	return &model.{{.PascalName}}{
		{{- range .Fields }}
		{{ .PascalName }}: req.{{ .PascalName }},
		{{- end }}
	}
}


func Update{{.PascalName}}ToModel(req *request.Update{{.PascalName}}) *model.{{.PascalName}} {
	return &model.{{.PascalName}}{
		ID: req.ID,
		{{- range .Fields }}
		{{ .PascalName }}: req.{{ .PascalName }},
		{{- end }}
	}
}
