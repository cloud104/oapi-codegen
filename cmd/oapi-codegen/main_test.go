package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/util"
)

func TestGenerate(t *testing.T) {
	spec, err := util.LoadSwagger("../../examples/minimal-client/api.yaml")
	if err != nil {
		t.Fatal(err)
	}

	outputDir := "generated"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}

	type args struct {
		spec *openapi3.T
		opts codegen.Configuration
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "minimal client",
			args: args{
				spec: spec,
				opts: codegen.Configuration{
					PackageName: "api",
					Generate: codegen.GenerateOptions{
						Client: true,
						Models: true,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := codegen.Generate(tt.args.spec, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Generate() error = %v, wantErr %v", err, tt.wantErr)
			}

			output := filepath.Join(outputDir, "api.gen.go")

			if err := os.WriteFile(output, []byte(got), 0644); err != nil {
				t.Fatal(err)
			}
		})
	}
}
