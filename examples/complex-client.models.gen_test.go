package examples

import (
	"reflect"
	"testing"
)

func TestBeer_ReadOnlyChanges(t *testing.T) {
	type fields struct {
		Brewery *Brewery
		ID      int
		Name    string
	}
	type args struct {
		patch map[string]any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "empty patch",
			fields: fields{
				ID:   7,
				Name: "Midnight Hops",
			},
			args: args{
				patch: map[string]any{},
			},
			want: nil,
		},
		{
			name: "unchanged read-only ID",
			fields: fields{
				ID:   7,
				Name: "Midnight Hops",
			},
			args: args{
				patch: map[string]any{
					"id": float64(7),
				},
			},
			want: nil,
		},
		{
			name: "changed read-only ID",
			fields: fields{
				ID:   7,
				Name: "Midnight Hops",
			},
			args: args{
				patch: map[string]any{
					"id": float64(8),
				},
			},
			want: []string{"id"},
		},
		{
			name: "changed writable field",
			fields: fields{
				ID:   7,
				Name: "Midnight Hops",
			},
			args: args{
				patch: map[string]any{
					"name": "Updated Beer",
				},
			},
			want: nil,
		},
		{
			name: "unchanged ID with changed writable field",
			fields: fields{
				ID:   7,
				Name: "Midnight Hops",
			},
			args: args{
				patch: map[string]any{
					"id":   float64(7),
					"name": "Updated Beer",
				},
			},
			want: nil,
		},
		{
			name: "changed ID with changed writable field",
			fields: fields{
				ID:   7,
				Name: "Midnight Hops",
			},
			args: args{
				patch: map[string]any{
					"id":   float64(8),
					"name": "Updated Beer",
				},
			},
			want: []string{"id"},
		},
		{
			name: "read-only ID set to null",
			fields: fields{
				ID:   7,
				Name: "Midnight Hops",
			},
			args: args{
				patch: map[string]any{
					"id": nil,
				},
			},
			want: []string{"id"},
		},
		{
			name: "read-only ID with invalid type",
			fields: fields{
				ID:   7,
				Name: "Midnight Hops",
			},
			args: args{
				patch: map[string]any{
					"id": "7",
				},
			},
			want: []string{"id"},
		},
		{
			name: "read-only ID with fractional value",
			fields: fields{
				ID: 7,
			},
			args: args{
				patch: map[string]any{
					"id": float64(7.5),
				},
			},
			want: []string{"id"},
		},
		{
			name: "unchanged nested read-only ID",
			fields: fields{
				ID: 7,
				Brewery: &Brewery{
					ID: new(10),
				},
			},
			args: args{
				patch: map[string]any{
					"brewery": map[string]any{
						"id": float64(10),
					},
				},
			},
			want: nil,
		},
		{
			name: "changed nested read-only ID",
			fields: fields{
				ID: 7,
				Brewery: &Brewery{
					ID: new(10),
				},
			},
			args: args{
				patch: map[string]any{
					"brewery": map[string]any{
						"id": float64(11),
					},
				},
			},
			want: []string{"brewery.id"},
		},
		{
			name: "changed nested writable field",
			fields: fields{
				ID: 7,
				Brewery: &Brewery{
					ID:   new(10),
					Name: new("Northern Trail"),
				},
			},
			args: args{
				patch: map[string]any{
					"brewery": map[string]any{
						"name": "Updated Brewery",
					},
				},
			},
			want: nil,
		},
		{
			name: "changed root and nested read-only IDs",
			fields: fields{
				ID: 7,
				Brewery: &Brewery{
					ID: new(10),
				},
			},
			args: args{
				patch: map[string]any{
					"id": float64(8),
					"brewery": map[string]any{
						"id": float64(11),
					},
				},
			},
			want: []string{"brewery.id", "id"},
		},
		{
			name: "nested object set to null",
			fields: fields{
				ID: 7,
				Brewery: &Brewery{
					ID: new(10),
				},
			},
			args: args{
				patch: map[string]any{
					"brewery": nil,
				},
			},
			want: []string{"brewery"},
		},
		{
			name: "nested object with invalid type",
			fields: fields{
				ID: 7,
				Brewery: &Brewery{
					ID: new(10),
				},
			},
			args: args{
				patch: map[string]any{
					"brewery": "invalid",
				},
			},
			want: []string{"brewery"},
		},
		{
			name: "nested object missing from current value",
			fields: fields{
				ID:      7,
				Brewery: nil,
			},
			args: args{
				patch: map[string]any{
					"brewery": map[string]any{
						"id": float64(10),
					},
				},
			},
			want: []string{"brewery"},
		},
		{
			name: "null nested object unchanged",
			fields: fields{
				ID:      7,
				Brewery: nil,
			},
			args: args{
				patch: map[string]any{
					"brewery": nil,
				},
			},
			want: nil,
		},
		{
			name: "unknown patch field",
			fields: fields{
				ID: 7,
			},
			args: args{
				patch: map[string]any{
					"unknown": "value",
				},
			},
			want: nil,
		},
		{
			name: "nil patch",
			fields: fields{
				ID: 7,
			},
			args: args{
				patch: nil,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Beer{
				Brewery: tt.fields.Brewery,
				ID:      tt.fields.ID,
				Name:    tt.fields.Name,
			}
			if got := v.ReadOnlyChanges(tt.args.patch); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadOnlyChanges() = %v, want %v", got, tt.want)
			}
		})
	}
}
