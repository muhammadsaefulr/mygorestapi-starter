package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type ModuleData struct {
	Name       string
	PascalName string
	ModulePath string
}

func main() {
	if len(os.Args) < 4 || os.Args[1] != "generate" || os.Args[2] != "module" {
		fmt.Println("Usage: go run generate_base_service.go generate module <module_name>")
		return
	}

	moduleName := os.Args[3]
	data := ModuleData{
		Name:       moduleName,
		PascalName: ToPascalCase(moduleName),
		ModulePath: "github.com/muhammadsaefulr/mygorestapi-starter",
	}

	files := []struct {
		Template string
		Path     string
	}{
		{"controller.go.tpl", "internal/delivery/http/controller/{{.Name}}_controller/{{.Name}}_controller.go"},
		{"router.go.tpl", "internal/delivery/http/router/{{.Name}}_router.go"},
		{"repository.go.tpl", "internal/repository/{{.Name}}/{{.Name}}_repository.go"},
		{"repository_impl.go.tpl", "internal/repository/{{.Name}}/{{.Name}}_repository_impl.go"},
		{"service.go.tpl", "internal/service/{{.Name}}_service/{{.Name}}_service.go"},
		{"service_impl.go.tpl", "internal/service/{{.Name}}_service/{{.Name}}_service_impl.go"},
		{"convert.go.tpl", "internal/shared/convert_types/{{.Name}}_converter.go"},
		{"dto/request.go.tpl", "internal/domain/dto/{{.Name}}/request/request.go"},
		{"dto/response.go.tpl", "internal/domain/dto/{{.Name}}/response/response.go"},
		{"model.go.tpl", "internal/domain/model/{{.Name}}_model.go"},
	}

	for _, f := range files {
		renderAndCreate(f.Template, f.Path, data)
	}
}

func renderAndCreate(templateFile, targetPath string, data ModuleData) {

	cwd, _ := os.Getwd()
	tmplPath := filepath.Join(cwd, "cmd/generator/templates", templateFile)
	tmplBytes, err := os.ReadFile(tmplPath)

	if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("Template not found: %s", tmplPath))
	}

	if err != nil {
		panic(err)
	}

	tmpl, err := template.New(templateFile).Parse(string(tmplBytes))
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		panic(err)
	}

	targetPathParsed := parsePath(targetPath, data)
	os.MkdirAll(filepath.Dir(targetPathParsed), os.ModePerm)

	if _, err := os.Stat(targetPathParsed); err == nil {
		fmt.Println("SKIP (already exists):", targetPathParsed)
		return
	}

	if err := os.WriteFile(targetPathParsed, buf.Bytes(), 0644); err != nil {
		panic(err)
	}

	fmt.Println("Created:", targetPathParsed)
}

func parsePath(path string, data ModuleData) string {
	path = strings.ReplaceAll(path, "{{.Name}}", data.Name)
	path = strings.ReplaceAll(path, "{{.PascalName}}", data.PascalName)
	return path
}

func ToPascalCase(s string) string {
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	words := strings.Fields(s)

	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, "")
}
