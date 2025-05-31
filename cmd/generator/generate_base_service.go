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
		fmt.Println("Usage: go run main.go generate module <name>")
		return
	}

	moduleName := os.Args[3]
	data := ModuleData{
		Name:       moduleName,
		PascalName: strings.Title(moduleName),
		ModulePath: "github.com/muhammadsaefulr/NimeStreamAPI", // 💡 hardcoded
	}

	files := []struct {
		Template string
		Path     string
	}{
		{"handler.go.tpl", "internal/delivery/http/handler/{{.Name}}_handler.go"},
		{"repository.go.tpl", "internal/repository/{{.Name}}/repository.go"},
		{"repository_impl.go.tpl", "internal/repository/{{.Name}}/repository_impl.go"},
		{"service.go.tpl", "internal/service/{{.Name}}_service/service.go"},
		{"service_impl.go.tpl", "internal/service/{{.Name}}_service/service_impl.go"},
		{"convert.go.tpl", "internal/domain/convert/{{.Name}}_converter.go"},
		{"dto/request.go.tpl", "internal/domain/dto/{{.Name}}/request/request.go"},
		{"dto/response.go.tpl", "internal/domain/dto/{{.Name}}/response/response.go"},
		{"model.go.tpl", "internal/domain/model/{{.Name}}/model.go"},
	}

	for _, f := range files {
		renderAndCreate(f.Template, f.Path, data)
	}
}

func renderAndCreate(templateFile, targetPath string, data ModuleData) {
	tmplPath := filepath.Join("templates", templateFile)
	tmplBytes, err := os.ReadFile(tmplPath)
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
