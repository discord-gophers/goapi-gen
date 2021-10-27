package packageA

//go:generate go run github.com/discord-gophers/goapi-gen/cmd/goapi-gen -generate types,skip-prune,spec -package=packageA -o externalref.gen.go -import-mapping=../packageB/spec.yaml:github.com/discord-gophers/goapi-gen/internal/test/externalref/packageB spec.yaml
