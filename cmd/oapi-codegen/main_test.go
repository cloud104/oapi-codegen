package main

import (
	"flag"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/util"
)

const examplesDir = "../../examples"

var writeResults = flag.Bool("write-results", false, "write generated results locally")

func TestGenerate(t *testing.T) {
	var specs []string

	err := filepath.WalkDir(examplesDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		switch strings.ToLower(filepath.Ext(path)) {
		case ".yaml", ".yml":
			specs = append(specs, path)
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(specs) == 0 {
		t.Fatalf("no YAML files found in %s", examplesDir)
	}

	for _, specPath := range specs {
		t.Run(specPath, func(t *testing.T) {
			// Arrange
			spec, err := util.LoadSwagger(specPath)
			if err != nil {
				t.Fatalf("load spec %q: %v", specPath, err)
			}

			// Act
			result, err := codegen.Generate(spec, codegen.Configuration{
				PackageName: "examples",
				Generate: codegen.GenerateOptions{
					Client: true,
					Models: true,
				},
				OutputOptions: codegen.OutputOptions{
					UserTemplates: map[string]string{
						"client.tmpl":                "../../templates/client.tmpl",
						"client-with-responses.tmpl": "../../templates/client-with-responses.tmpl",
					},
				},
			})

			// Assert
			if err != nil {
				t.Fatalf("generate %q: %v", specPath, err)
			}

			if !*writeResults {
				return
			}

			outputPath := strings.TrimSuffix(specPath, filepath.Ext(specPath)) + ".gen.go"
			if err := os.WriteFile(outputPath, []byte(result), 0o644); err != nil {
				t.Fatalf("write %q: %v", outputPath, err)
			}
		})
	}
}
