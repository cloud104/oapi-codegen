package examples

import (
	"strings"
	"testing"
)

func TestJSONValidator_Validate(t *testing.T) {
	validator := NewJSONValidator()

	if err := validator.Register("CreateUser", CreateUserJSONRequestBodySchema); err != nil {
		t.Fatalf("failed to register schema: %v", err)
	}

	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			name:    "valid user",
			data:    `{"name":"John","email":"john@example.com"}`,
			wantErr: false,
		},
		{
			name:    "missing required field",
			data:    `{"name":"John"}`,
			wantErr: true,
		},
		{
			name:    "invalid email format",
			data:    `{"name":"John","email":"invalid-email"}`,
			wantErr: true,
		},
		{
			name:    "unexpected property",
			data:    `{"name":"John","email":"john@example.com","age":30}`,
			wantErr: true,
		},
		{
			name:    "invalid field type",
			data:    `{"name":123,"email":"john@example.com"}`,
			wantErr: true,
		},
		{
			name:    "unregistered schema",
			data:    `{"name":"John","email":"john@example.com"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schemaName := "CreateUser"
			if tt.name == "unregistered schema" {
				schemaName = "UnknownSchema"
			}

			err := validator.Validate(schemaName, strings.NewReader(tt.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate(%q, %q) error = %v, wantErr %v", schemaName, tt.data, err, tt.wantErr)
			}
		})
	}
}
