// Copyright 2019 DeepMap, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package codegen

import (
	"net/http"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestGenerateDefaultOperationID(t *testing.T) {
	type test struct {
		op      string
		path    string
		want    string
		wantErr bool
	}

	suite := []test{
		{
			op:      http.MethodGet,
			path:    "/v1/foo/bar",
			want:    "GetV1FooBar",
			wantErr: false,
		},
		{
			op:      http.MethodGet,
			path:    "/v1/foo/bar/",
			want:    "GetV1FooBar",
			wantErr: false,
		},
		{
			op:      http.MethodPost,
			path:    "/v1",
			want:    "PostV1",
			wantErr: false,
		},
		{
			op:      http.MethodPost,
			path:    "v1",
			want:    "PostV1",
			wantErr: false,
		},
		{
			path:    "v1",
			want:    "",
			wantErr: true,
		},
		{
			path:    "",
			want:    "PostV1",
			wantErr: true,
		},
	}

	for _, test := range suite {
		got, err := generateDefaultOperationID(test.op, test.path)
		if err != nil {
			if !test.wantErr {
				t.Fatalf("did not expected error but got %v", err)
			}
		}

		if test.wantErr {
			return
		}
		if got != test.want {
			t.Fatalf("Operation ID generation error. Want [%v] Got [%v]", test.want, got)
		}
	}
}

func TestParameterDefinition_GoVariableName(t *testing.T) {
	type fields struct {
		ParamName string
		In        string
		Required  bool
		Spec      *openapi3.Parameter
		Schema    Schema
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "simple",
			fields: fields{
				ParamName: "foo",
			},
			want: "foo",
		},
		{
			name: "starts with number",
			fields: fields{
				ParamName: "1foo",
			},
			want: "n1foo",
		},
		{
			name: "contains underscode",
			fields: fields{
				ParamName: "foo_bar",
			},
			want: "fooBar",
		},
		{
			name: "contains dash",
			fields: fields{
				ParamName: "foo-bar",
			},
			want: "fooBar",
		},
		{
			name: "contains dash and underscore",
			fields: fields{
				ParamName: "foo-bar_baz",
			},
			want: "fooBarBaz",
		},
		{
			name: "contains URI",
			fields: fields{
				ParamName: "foo-bar_baz_uri",
			},
			want: "fooBarBazURI",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pd := ParameterDefinition{
				ParamName: tt.fields.ParamName,
				In:        tt.fields.In,
				Required:  tt.fields.Required,
				Spec:      tt.fields.Spec,
				Schema:    tt.fields.Schema,
			}
			if got := pd.GoVariableName(); got != tt.want {
				t.Errorf("ParameterDefinition.GoVariableName() = %v, want %v", got, tt.want)
			}
		})
	}
}
