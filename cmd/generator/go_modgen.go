package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

// Field sekarang menyimpan nama PascalCase, SnakeCase, tipe, dan validasi
type Field struct {
	PascalName    string // Contoh: ProductName
	SnakeCaseName string // Contoh: product_name
	Type          string // Contoh: string
	Validation    string // Contoh: required,min=3
}

type ModuleData struct {
	Name       string
	PascalName string
	ModulePath string
	Fields     []Field
}

// Regex untuk konversi ke snake_case
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func main() {
	// Definisikan flag
	modelFields := flag.String("model", "", "Model fields in format: \"FieldName1:type1:validation1,FieldName2:type2:validation2\"")

	// Parsing flags
	flag.Parse()

	args := flag.Args()
	if len(args) < 3 || args[0] != "generate" || args[1] != "module" {
		fmt.Println("Usage: go run main.go generate module <module_name> --model=\"FieldName:type:validation,...\"")
		fmt.Println("Example: go run main.go generate module product --model=\"ProductName:string:required,ProductCategory:string:required\"")
		return
	}

	if *modelFields == "" {
		fmt.Println("Error: --model flag is required.")
		fmt.Println("Example: --model=\"ProductName:string:required\"")
		return
	}

	moduleName := args[2]

	// Parsing string dari flag --model
	fieldDefs := strings.Split(*modelFields, ",")
	fields := make([]Field, len(fieldDefs))

	for i, def := range fieldDefs {
		parts := strings.Split(strings.TrimSpace(def), ":")
		if len(parts) < 2 {
			fmt.Printf("Invalid field format: %s. Must be 'FieldName:type' or 'FieldName:type:validation'.\n", def)
			return
		}

		var validation string
		if len(parts) > 2 {
			validation = parts[2]
		}

		pascalName := parts[0]
		fields[i] = Field{
			PascalName:    pascalName,
			SnakeCaseName: ToSnakeCase(pascalName),
			Type:          parts[1],
			Validation:    validation,
		}
	}

	data := ModuleData{
		Name:       moduleName,
		PascalName: ToPascalCase(moduleName),
		ModulePath: "github.com/muhammadsaefulr/mygorestapi-starter", // Sesuaikan jika perlu
		Fields:     fields,
	}

	// (Sisa kode untuk generate file sama seperti sebelumnya)
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

// Fungsi renderAndCreate, parsePath, dan ToPascalCase tetap sama
// ... (salin fungsi-fungsi tersebut dari kode Anda sebelumnya)
func renderAndCreate(templateFile, targetPath string, data ModuleData) {
	cwd, _ := os.Getwd()
	tmplPath := filepath.Join(cwd, "cmd/generator/templates", templateFile)

	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		if os.IsNotExist(err) {
			panic(fmt.Sprintf("Template not found: %s", tmplPath))
		}
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
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")
	words := strings.Fields(s)

	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, "")
}
