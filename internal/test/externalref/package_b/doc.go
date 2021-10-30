package package_b

//go:generate go run github.com/discord-gophers/goapi-gen --generate types,skip-prune,spec --package=package_b -o externalref.gen.go spec.yaml
