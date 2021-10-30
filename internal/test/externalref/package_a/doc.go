package package_a

//go:generate go run github.com/discord-gophers/goapi-gen --generate types,skip-prune,spec --package=package_a -o externalref.gen.go --import-mapping=../package_b/spec.yaml:github.com/discord-gophers/goapi-gen/internal/test/externalref/package_b spec.yaml
