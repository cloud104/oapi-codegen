package naming

import "testing"

func TestCaser_Camel(t *testing.T) {
	tests := []struct {
		name        string
		initialisms []string
		input       string
		want        string
	}{
		{
			name:  "simple snake case",
			input: "hello_world",
			want:  "helloWorld",
		},
		{
			name:        "single initialism",
			initialisms: []string{"ID"},
			input:       "user_ID",
			want:        "userID",
		},
		{
			name:        "multiple initialisms",
			initialisms: []string{"ID", "HTTP"},
			input:       "user_ID_HTTP",
			want:        "userIDHTTP",
		},
		{
			name:        "embedded initialism",
			initialisms: []string{"HTTP"},
			input:       "HTTPServer",
			want:        "httpServer",
		},
		{
			name:        "starts with two consecutive initialisms",
			initialisms: []string{"HTTP", "ID"},
			input:       "HTTPIDServer",
			want:        "httpIDServer",
		},
		{
			name:        "overlapping initialisms require longest match first",
			initialisms: []string{"HTTP", "HTTPS"},
			input:       "HTTPSClient",
			want:        "httpsClient",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			caser := NewCaser(tt.initialisms)

			if got := caser.Camel(tt.input); got != tt.want {
				t.Errorf("Camel(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestCaser_Pascal(t *testing.T) {
	tests := []struct {
		name        string
		initialisms []string
		input       string
		want        string
	}{
		{
			name:  "simple snake case",
			input: "hello_world",
			want:  "HelloWorld",
		},
		{
			name:        "single initialism",
			initialisms: []string{"ID"},
			input:       "user_ID",
			want:        "UserID",
		},
		{
			name:        "multiple initialisms",
			initialisms: []string{"ID", "HTTP"},
			input:       "user_ID_HTTP",
			want:        "UserIDHTTP",
		},
		{
			name:        "embedded initialism",
			initialisms: []string{"HTTP"},
			input:       "HTTPServer",
			want:        "HTTPServer",
		},
		{
			name:        "starts with two consecutive initialisms",
			initialisms: []string{"HTTP", "ID"},
			input:       "HTTPIDServer",
			want:        "HTTPIDServer",
		},
		{
			name:        "overlapping initialisms require longest match first",
			initialisms: []string{"HTTP", "HTTPS"},
			input:       "HTTPSClient",
			want:        "HTTPSClient",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			caser := NewCaser(tt.initialisms)

			if got := caser.Pascal(tt.input); got != tt.want {
				t.Errorf("Pascal(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
