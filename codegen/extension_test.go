package codegen

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_extImportPath(t *testing.T) {
	type args struct {
		extPropValue interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    extImportPathDetails
		wantErr bool
	}{
		{
			name:    "success",
			args:    args{json.RawMessage(`"uint64"`)},
			want:    extImportPathDetails{Type: "uint64"},
			wantErr: false,
		},
		{
			name:    "type conversion error",
			args:    args{nil},
			want:    extImportPathDetails{},
			wantErr: true,
		},
		{
			name:    "json unmarshal error",
			args:    args{json.RawMessage("invalid json format")},
			want:    extImportPathDetails{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extImportPath(tt.args.extPropValue)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
