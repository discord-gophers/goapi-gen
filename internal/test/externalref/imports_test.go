package externalref

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/discord-gophers/goapi-gen/internal/test/externalref/package_a"
	"github.com/discord-gophers/goapi-gen/internal/test/externalref/package_b"
)

func TestParameters(t *testing.T) {
	b := &package_b.ObjectB{}
	_ = Container{
		ObjectA: &package_a.ObjectA{ObjectB: b},
		ObjectB: b,
	}
}

func TestGetSwagger(t *testing.T) {
	_, err := package_b.GetSwagger()
	require.Nil(t, err)

	_, err = package_b.GetSwagger()
	require.Nil(t, err)

	_, err = GetSwagger()
	require.Nil(t, err)
}
