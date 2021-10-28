package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	client3, err := NewClient(
		"",
		WithClientBaseURL(withTrailingSlash),
	)
	assert.NoError(t, err)

	client4, err := NewClient(
		"",
		WithClientBaseURL(withoutTrailingSlash),
	)
	assert.NoError(t, err)

	expectedURL := withTrailingSlash

	assert.Equal(t, expectedURL, client1.Server)
	assert.Equal(t, expectedURL, client2.Server)
	assert.Equal(t, expectedURL, client3.Server)
	assert.Equal(t, expectedURL, client4.Server)
}
