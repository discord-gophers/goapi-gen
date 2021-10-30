package securityprovider

import (
	"testing"

	"github.com/discord-gophers/goapi-gen/internal/test/client"
	"github.com/stretchr/testify/assert"
)

var withTrailingSlash = "https://my-api.com/some-base-url/v1/"

func TestSecurityProviders(t *testing.T) {
	bearer, err := NewBearerToken("mytoken")
	assert.NoError(t, err)
	client1, err := client.NewClient(
		withTrailingSlash,
		client.WithRequestEditorFn(bearer.Intercept),
	)
	assert.NoError(t, err)

	apiKey, err := NewAPIKey("cookie", "apikey", "mykey")
	assert.NoError(t, err)
	client2, err := client.NewClient(
		withTrailingSlash,
		client.WithRequestEditorFn(apiKey.Intercept),
	)
	assert.NoError(t, err)

	basicAuth, err := NewBasicAuth("username", "password")
	assert.NoError(t, err)
	client3, err := client.NewClient(
		withTrailingSlash,
		client.WithRequestEditorFn(basicAuth.Intercept),
	)
	assert.NoError(t, err)

	assert.Equal(t, withTrailingSlash, client1.Server)
	assert.Equal(t, withTrailingSlash, client2.Server)
	assert.Equal(t, withTrailingSlash, client3.Server)
}
