package packageB

//go:generate go run github.com/discord-gophers/goapi-gen/cmd/goapi-gen --generate types,skip-prune,spec --package=packageB -o externalref.gen.go spec.yaml
