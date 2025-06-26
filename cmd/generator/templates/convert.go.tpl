package convert_types

import (
	"{{.ModulePath}}/internal/domain/dto/{{.Name}}/request"
	model "{{.ModulePath}}/internal/domain/model"
)

func Create{{.PascalName}}ToModel(req *request.Create{{.PascalName}}) *model.{{.PascalName}} {
	return &model.{{.PascalName}}{
		// TODO: sesuaikan field sesuai model
		// Example:
		// Name: req.Name,
	}
}

func Update{{.PascalName}}ToModel(req *request.Update{{.PascalName}}) *model.{{.PascalName}} {
	return &model.{{.PascalName}}{
		// TODO: sesuaikan field sesuai model
		// Example:
		// Name: req.Name,
	}
}
