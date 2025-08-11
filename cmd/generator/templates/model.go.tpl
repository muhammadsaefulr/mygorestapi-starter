package model

import "time"

type {{.PascalName}} struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	{{- range .Fields }}
	{{ .PascalName }} {{ .Type }} `gorm:"column:{{ .SnakeCaseName }}" json:"{{ .SnakeCaseName }}"`
	{{- end }}
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (m *{{.PascalName}}) TableName() string {
	return "{{.Name}}s"
}