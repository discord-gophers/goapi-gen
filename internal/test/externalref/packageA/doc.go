package packageA

//go:generate go run github.com/discord-gophers/goapi-gen/cmd/oapi-codegen -generate types,skip-prune,spec --package=packageA -o externalref.gen.go --import-mapping=../packageB/spec.yaml:github.com/deepmap/oapi-codegen/internal/test/externalref/packageB spec.yaml
