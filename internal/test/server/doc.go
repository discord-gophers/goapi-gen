package server

//go:generate go run github.com/discord-gophers/goapi-gen/cmd/goapi-gen --generate=types,server --package=server -o server.gen.go ../test-schema.yaml
//go:generate go run github.com/matryer/moq -out server_moq.gen.go . ServerInterface
