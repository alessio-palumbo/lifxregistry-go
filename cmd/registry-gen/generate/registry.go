package generate

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/alessio-palumbo/lifxregistry-go/cmd/registry-gen/decode"
	"golang.org/x/tools/imports"
)

//go:embed templates/*
var templates embed.FS

var generatorHeader = `
// Code generated. DO NOT EDIT.
// Source: https://github.com/LIFX/public-protocol@%s
// Generated: %s
`

var now = time.Now

func GenerateProductsRegistry(products []decode.Product, outputRoot string, sourceCommit string) error {
	header := fmt.Sprintf(generatorHeader, sourceCommit, now().UTC().Format(time.RFC3339))
	if err := generateFromTemplate("templates/products.tmpl", filepath.Join(outputRoot, "products.go"), header, products); err != nil {
		return err
	}
	return nil
}

func generateFromTemplate(tmplPath, outputPath string, header string, data any) error {
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating output directory [%s]: %w", dir, err)
	}

	tmplBytes, err := templates.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("reading template %s: %w", tmplPath, err)
	}

	tmpl, err := template.New(filepath.Base(tmplPath)).Funcs(template.FuncMap{
		"header": func() string { return header },
	}).Parse(string(tmplBytes))
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	formatted, err := imports.Process(outputPath, buf.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("goimports failed [%s]: %w", tmplPath, err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("creating output file [%s]: %w", outputPath, err)
	}
	defer f.Close()

	if _, err := f.Write(formatted); err != nil {
		return fmt.Errorf("writing formatted output: %w", err)
	}

	return nil
}
