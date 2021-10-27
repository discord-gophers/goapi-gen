package runtime

import (
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type InnerObject struct {
	Name string
	ID   int
}

// These are all possible field types, mandatory and optional.
type AllFields struct {
	I   int          `json:"i"`
	Oi  *int         `json:"oi,omitempty"`
	F   float32      `json:"f"`
	Of  *float32     `json:"of,omitempty"`
	B   bool         `json:"b"`
	Ob  *bool        `json:"ob,omitempty"`
	As  []string     `json:"as"`
	Oas *[]string    `json:"oas,omitempty"`
	O   InnerObject  `json:"o"`
	Oo  *InnerObject `json:"oo,omitempty"`
	D   MockBinder   `json:"d"`
	Od  *MockBinder  `json:"od,omitempty"`
}

func TestDeepObject(t *testing.T) {
	oi := int(5)
	of := float32(3.7)
	ob := true
	oas := []string{"foo", "bar"}
	oo := InnerObject{
		Name: "Marcin Romaszewicz",
		ID:   123,
	}
	d := MockBinder{Time: time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)}

	srcObj := AllFields{
		I:   12,
		Oi:  &oi,
		F:   4.2,
		Of:  &of,
		B:   true,
		Ob:  &ob,
		As:  []string{"hello", "world"},
		Oas: &oas,
		O: InnerObject{
			Name: "Joe Schmoe",
			ID:   456,
		},
		Oo: &oo,
		D:  d,
		Od: &d,
	}

	marshaled, err := MarshalDeepObject(srcObj, "p")
	require.NoError(t, err)
	t.Log(marshaled)

	params := make(url.Values)
	marshaledParts := strings.Split(marshaled, "&")
	for _, p := range marshaledParts {
		parts := strings.Split(p, "=")
		require.Equal(t, 2, len(parts))
		params.Set(parts[0], parts[1])
	}

	var dstObj AllFields
	err = UnmarshalDeepObject(&dstObj, "p", params)
	require.NoError(t, err)
	assert.EqualValues(t, srcObj, dstObj)
}

func TestUnmarshalDeepObject(t *testing.T) {
	type args struct {
		dst       interface{}
		paramName string
		params    url.Values
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "simple",
			args: args{
				dst: &struct {
					I int `json:"i"`
				}{},
				paramName: "p",
				params: url.Values{
					"p[i]": {"12"},
				},
			},
			want: &struct {
				I int `json:"i"`
			}{I: 12},
			wantErr: false,
		},
		{
			name: "no values",
			args: args{
				dst: &struct {
					I int `json:"i"`
				}{},
				paramName: "p",
				params: url.Values{
					"p": {},
				},
			},
			want: &struct {
				I int `json:"i"`
			}{I: 0},
			wantErr: false,
		},
		{
			name: "advanced",
			args: args{
				paramName: "deepObj",
				dst: &struct {
					ID      int  `json:"Id"`
					IsAdmin bool `json:"IsAdmin"`
					Object  struct {
						FirstName string `json:"firstName"`
						Role      string `json:"role"`
					} `json:"Object"`
				}{},
				params: url.Values{
					"deepObj[Id]":                {"12345"},
					"deepObj[IsAdmin]":           {"true"},
					"deepObj[Object][firstName]": {"Alex"},
					"deepObj[Object][role]":      {"admin"},
				},
			},
			want: &struct {
				ID      int  `json:"Id"`
				IsAdmin bool `json:"IsAdmin"`
				Object  struct {
					FirstName string `json:"firstName"`
					Role      string `json:"role"`
				} `json:"Object"`
			}{
				ID:      12345,
				IsAdmin: true,
				Object: struct {
					FirstName string `json:"firstName"`
					Role      string `json:"role"`
				}{
					FirstName: "Alex",
					Role:      "admin",
				},
			},
		},
		{
			name: "invalid",
			args: args{
				dst:       &struct{}{},
				paramName: "p",
				params: url.Values{
					"p[[i]": {"12"},
				},
			},
			want:    &struct{}{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalDeepObject(tt.args.dst, tt.args.paramName, tt.args.params)

			if (err != nil) != tt.wantErr {
				require.NoError(t, err, "UnmarshalDeepObject() error")
				return
			}

			require.Equal(t, tt.want, tt.args.dst, "UnmarshalDeepObject() invalid")
		})
	}
}
