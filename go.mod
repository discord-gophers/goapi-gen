module github.com/discord-gophers/goapi-gen

replace (
	github.com/getkin/kin-openapi => github.com/diamondburned/kin-openapi v0.80.1-0.20240812070530-f9deb571c539
	github.com/ghodss/yaml/v2 => github.com/diamondburned/yaml/v2 v2.0.0-20240812065612-baf990d70122
)

require (
	github.com/getkin/kin-openapi v0.80.0
	github.com/go-chi/chi/v5 v5.0.4
	github.com/golangci/lint-1 v0.0.0-20181222135242-d2cdd8c08219
	github.com/kenshaw/snaker v0.1.6
	github.com/matryer/moq v0.3.1
	github.com/stretchr/testify v1.7.0
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/tools v0.3.0
	gopkg.in/yaml.v3 v3.0.1
)

require github.com/ghodss/yaml/v2 v2.0.0-00010101000000-000000000000 // indirect

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-chi/render v1.0.1
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	golang.org/x/mod v0.7.0 // indirect
	golang.org/x/sys v0.2.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

go 1.17
