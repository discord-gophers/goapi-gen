package client

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildURL(t *testing.T) {
	mustURL := func(u *url.URL, _ error) *url.URL {
		return u
	}

	type testCase struct {
		name    string
		baseURL string
		path    map[string]string
		params  map[string]string
		want    string
	}

	tt := []testCase{
		{
			name:    "both",
			baseURL: "https://my-api.com/some-base-url/v1/{id}/{user}",
			path: map[string]string{
				"id":   "123",
				"user": "henry",
			},
			params: map[string]string{
				"message_id": "123456789",
			},
			want: "https://my-api.com/some-base-url/v1/123/henry?message_id=123456789",
		},
		{
			name:    "no path",
			baseURL: "https://my-api.com/some-base-url/v1/?id=123&message_id=123456789&user=henry",
			path:    nil,
			params: map[string]string{
				"message_id": "123456789",
				"id":         "123",
				"user":       "henry",
			},
			want: "https://my-api.com/some-base-url/v1/?id=123&message_id=123456789&user=henry",
		},
		{
			name:    "no params",
			baseURL: "https://my-api.com/some-base-url/v1/{id}/{message_id}/{user}",
			path: map[string]string{
				"message_id": "123456789",
				"id":         "123",
				"user":       "henry",
			},
			params: nil,
			want:   "https://my-api.com/some-base-url/v1/123/123456789/henry",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := buildURL(tc.baseURL, tc.path, tc.params)
			assert.Equal(t, mustURL(url.Parse(tc.want)).String(), got)
		})
	}
}

func TestTemp(t *testing.T) {
	var (
		withTrailingSlash    string = "https://my-api.com/some-base-url/v1/"
		withoutTrailingSlash string = "https://my-api.com/some-base-url/v1"
	)

	client1, err := NewClient(
		withTrailingSlash,
	)
	assert.NoError(t, err)

	client2, err := NewClient(
		withoutTrailingSlash,
	)
	assert.NoError(t, err)

	expectedURL := withTrailingSlash

	assert.Equal(t, expectedURL, client1.BaseURL)
	assert.Equal(t, expectedURL, client2.BaseURL)
}
