# OpenAPI Server Code Generator

[![Go Reference](https://pkg.go.dev/badge/github.com/discord-gophers/goapi-gen.svg)](https://pkg.go.dev/github.com/discord-gophers/goapi-gen)

This package contains a set of utilities for generating Go boilerplate code for
services based on
[OpenAPI 3.0](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.0.md)
API definitions. When working with services, it's important to have an API
contract which servers and clients both implement to minimize the chances of
incompatibilities. It's tedious to generate Go models which precisely correspond to
OpenAPI specifications, so let our code generator do that work for you, so that
you can focus on implementing the business logic for your service.

We have chosen to use [Chi](https://github.com/go-chi/chi) as
our HTTP routing engine, due to its speed, simplicity, and compatibility
with `net/http`.

This package tries to be too simple rather than too generic, so we've made some
design decisions in favor of simplicity, knowing that we can't generate strongly
typed Go code for all possible OpenAPI Schemas.

This repository is a hard fork of [deepmap/oapi-codegen](https://github.com/deepmap/oapi-codegen).
This new version plans to diverge from the original repository with different design goals and more
emphasis on `go-chi`.

## Overview

We're going to use the OpenAPI example of the
[Expanded Petstore](https://github.com/OAI/OpenAPI-Specification/blob/master/examples/v3.0/petstore-expanded.yaml)
in the descriptions below, please have a look at it.

In order to create a Go server to serve this exact schema, you would have to
write a lot of boilerplate code to perform all the marshalling and unmarshalling
into objects which match the OpenAPI 3.0 definition. The code generator in this
directory does a lot of that for you. You would run it like so:

    go install github.com/discord-gophers/goapi-gen@latest
    goapi-gen petstore-expanded.yaml > petstore.gen.go

Let's go through that `petstore.gen.go` file to show you everything which was
generated.

## Generated Server Boilerplate

The `/components/schemas` section in OpenAPI defines reusable objects, so Go
types are generated for these. The Pet Store example defines `Error`, `Pet`,
`Pets` and `NewPet`, so we do the same in Go:

```go
// Type definition for component schema "Error"
type Error struct {
    Code    int32  `json:"code"`
    Message string `json:"message"`
}

// Type definition for component schema "NewPet"
type NewPet struct {
    Name string  `json:"name"`
    Tag  *string `json:"tag,omitempty"`
}

// Type definition for component schema "Pet"
type Pet struct {
    // Embedded struct due to allOf(#/components/schemas/NewPet)
    NewPet
    // Embedded fields due to inline allOf schema
    Id int64 `json:"id"`
}

// Type definition for component schema "Pets"
type Pets []Pet
```

It's best to define objects under `/components` field in the schema, since
those will be turned into named Go types. If you use inline types in your
handler definitions, we will generate inline, anonymous Go types, but those
are more tedious to deal with since you will have to redeclare them at every
point of use.

For each element in the `paths` map in OpenAPI, we will generate a Go handler
function in an interface object. Here is the generated Go interface for our
Chi server.

```go
type ServerInterface interface {
    //  (GET /pets)
    FindPets(w http.ResponseWriter, r *http.Request, params FindPetsParams)
    //  (POST /pets)
    AddPet(w http.ResponseWriter, r *http.Request)
    //  (DELETE /pets/{id})
    DeletePet(w http.ResponseWriter, r *http.Request, id int64)
    //  (GET /pets/{id})
    FindPetById(w http.ResponseWriter, r *http.Request, id int64)
}
```

These are the functions which you will implement yourself in order to create
a server conforming to the API specification.

Notice that `FindPetById` takes a parameter `id int64`. All path arguments
will be passed as arguments to your function, since they are mandatory.

Remaining arguments can be passed in headers, query arguments or cookies. Those
will be written to a `params` object. Look at the `FindPets` function above, it
takes as input `FindPetsParams`, which is defined as follows:

```go
// Parameters object for FindPets
type FindPetsParams struct {
    Tags  *[]string `json:"tags,omitempty"`
    Limit *int32   `json:"limit,omitempty"`
}
```

The HTTP query parameter `limit` turns into a Go field named `Limit`. It is
passed by pointer, since it is an optional parameter. If the parameter is
specified, the pointer will be non-`nil`, and you can read its value.

If you changed the OpenAPI specification to make the parameter required, the
`FindPetsParams` structure will contain the type by value:

```go
type FindPetsParams struct {
    Tags  *[]string `json:"tags,omitempty"`
    Limit int32     `json:"limit"`
}
```

### Registering handlers

You can register handlers when generating a server with `-generate server`.

<details><summary><code>Chi</code></summary>

Code generated using `-generate server`.

```go
type PetStoreImpl struct {}
func (*PetStoreImpl) GetPets(w http.ResponseWriter, r *http.Request) {
    // Implement me
}

func SetupHandler() {
    var myApi PetStoreImpl

    r := chi.NewRouter()
    r.Mount("/", Handler(&myApi))
}
```

</summary></details>

<details><summary><code>net/http</code></summary>

[Chi](https://github.com/go-chi/chi) is 100% compatible with `net/http` allowing the following with code generated using `-generate server`.

```go
type PetStoreImpl struct {}
func (*PetStoreImpl) GetPets(w http.ResponseWriter, r *http.Request) {
    // Implement me
}

func SetupHandler() {
    var myApi PetStoreImpl

    http.Handle("/", Handler(&myApi))
}
```

</summary></details>

#### Additional Properties in type definitions

[OpenAPI Schemas](https://swagger.io/specification/#schemaObject) implicitly
accept `additionalProperties`, meaning that any fields provided, but not explicitly
defined via properties on the schema are accepted as input, and propagated. When
unspecified, the `additionalProperties` field is assumed to be `true`.

Additional properties are tricky to support in Go with typing, and require
lots of boilerplate code, so in this library, we assume that `additionalProperties`
defaults to `false` and we don't generate this boilerplate. If you would like
an object to accept `additionalProperties`, specify a schema for `additionalProperties`.

Say we declared `NewPet` above like so:

```yaml
NewPet:
  required:
    - name
  properties:
    name:
      type: string
    tag:
      type: string
  additionalProperties:
    type: string
```

The Go code for `NewPet` would now look like this:

```go
// NewPet defines model for NewPet.
type NewPet struct {
	Name                 string            `json:"name"`
	Tag                  *string           `json:"tag,omitempty"`
	AdditionalProperties map[string]string `json:"-"`
}
```

The additionalProperties, of type `string` become `map[string]string`, which maps
field names to instances of the `additionalProperties` schema.

```go
// Getter for additional properties for NewPet. Returns the specified
// element and whether it was found
func (a NewPet) Get(fieldName string) (value string, found bool) {...}

// Setter for additional properties for NewPet
func (a *NewPet) Set(fieldName string, value string) {...}

// Override default JSON handling for NewPet to handle additionalProperties
func (a *NewPet) UnmarshalJSON(b []byte) error {...}

// Override default JSON handling for NewPet to handle additionalProperties
func (a NewPet) MarshalJSON() ([]byte, error) {...}w
```

There are many special cases for `additionalProperties`, such as having to
define types for inner fields which themselves support additionalProperties, and
all of them are tested via the `internal/test/components` schemas and tests. Please
look through those tests for more usage examples.

## Extensions

`goapi-gen` supports the following extended properties:

- `x-go-type`: specifies Go type name. It allows you to specify the type name for a schema, and
  will override any default value. This extended property isn't supported in all parts of
  OpenAPI, so please refer to the spec as to where it's allowed. Swagger validation tools will
  flag incorrect usage of this property.
- `x-go-extra-tags`: adds extra Go field tags to the generated struct field. This is
  useful for interfacing with tag based ORM or validation libraries. The extra tags that
  are added are in addition to the regular json tags that are generated. If you specify your
  own `json` tag, you will override the default one.

  ```yaml
  components:
    schemas:
      Object:
        properties:
          name:
            type: string
            x-go-extra-tags:
              tag1: value1
              tag2: value2
  ```

  In the example above, field `name` will be declared as:

  ```go
  Name string `json:"name" tag1:"value1" tag2:"value2"`
  ```

- `x-go-middlewares`: specifies a list of tagged middlewares. These can be specific
  middlewares that are operation-specific, as well as path-specific. This is very useful when you
  want to give a specific routes middleware, but not to all operations. The middleware are always
  called in the order of definition. If the tagged middleware is not defined, panic will be called while calling `Handler`.

  ```yaml
  /pets:
    x-go-middlewares: [validateJSON]
      get:
        x-go-middlewares: [limit]
  ```

  In the example above, the following middleware calls will be added to your handler:

  ```go
  // Operation specific middleware
  handler = siw.Middlewares["validateJSON"](handler).ServeHTTP
  handler = siw.Middlewares["limit"](handler).ServeHTTP
  ```

- `x-go-optional-value`: boolean, forces the generator to output value types (as
  opposed to pointer nil-able types) for all optional fields. This is
  particularly useful for when there is no practical difference between an empty
  omitted `string` and an omitted nil `*string`.

  This property can go in either the property value or the object attribute
  itself:

  ```yaml
  components:
    schemas:
      Object:
        x-go-optional-value: true # valid
        properties:
          name:
          type: string
          x-go-optional-value: false # valid, overrides
  ```

- `x-go-string`: boolean, makes the generator add a `,string` attribute into the
  JSON struct tag of an object's field. This is useful for sending large numbers
  over JSON to JavaScript (or other dynamically-typed languages) where they
  would be parsed to floats otherwise.

- `x-go-omitempty`: boolean, makes the generator add a `,omitempty` attribute
  into the JSON struct tag. Note that if a field doesn't appear in `required`,
  then it will have an `omitempty` by default and its type will be a nil-able
  pointer.

## Using `goapi-gen`

[Usage details](docs.md)

The default options for `goapi-gen` will generate everything; server,
type definitions and embedded swagger spec, but you can generate subsets of
those via the `-generate` flag. It defaults to `types,server,spec`, but
you can specify any combination of those.

- `types`: generate all type definitions for all types in the OpenAPI spec. This
  will be everything under `#components`, as well as request parameter, request
  body, and response type objects.
- `server`: generate the Chi server boilerplate. This code is dependent on
  that produced by the `types` target.
- `spec`: embed the OpenAPI spec into the generated code as a gzipped blob. This
- `skip-fmt`: skip running `goimports` on the generated code. This is useful for debugging
  the generated file in case the spec contains weird strings.
- `skip-prune`: skip pruning unused components from the spec prior to generating
  the code.
- `import-mapping`: specifies a map of references external OpenAPI specs to go
  Go include paths. Please see below.

So, for example, if you would like to produce only the server code, you could
run `goapi-gen --generate types,server`. You could generate `types` and
`server` into separate files, but both are required for the server code.

`goapi-gen` can filter paths base on their tags in the openapi definition.
Use either `--include-tags` or `--exclude-tags` followed by a comma-separated list
of tags. For instance, to generate a server that serves all paths except those
tagged with `auth` or `admin`, use the argument, `--exclude-tags="auth,admin"`.
To generate a server that only handles `admin` paths, use the argument
`--include-tags="admin"`. When neither of these arguments is present, all paths
are generated.

`goapi-gen` can filter schemas based on the option `--exclude-schemas`, which is
a comma separated list of schema names. For instance, `--exclude-schemas=Pet,NewPet`
will exclude from generation schemas `Pet` and `NewPet`. This allow to have a
in the same package a manually defined structure or interface and refer to it
in the openapi spec.

Since `go generate` commands must be a single line, all the options above can make
them pretty unwieldy, so you can specify all of the options in a configuration
file via the `--config` option. Please see the test under
[`/internal/test/externalref/`](https://github.com/discord-gophers/goapi-gen/blob/main/internal/test/externalref/externalref.cfg.yaml)
for an example. The structure of the file is as follows:

```yaml
output: externalref.gen.go
package: externalref
generate:
  - types
  - skip-prune
import-mapping:
  ./packageA/spec.yaml: github.com/discord-gophers/goapi-gen/internal/test/externalref/packageA
  ./packageB/spec.yaml: github.com/discord-gophers/goapi-gen/internal/test/externalref/packageB
```

Have a look at [`parse.go`](https://github.com/discord-gophers/goapi-gen/blob/main/parse.go#L28-L39)
to see all the fields on the configuration structure.

### Import Mappings

OpenAPI specifications may contain references to other OpenAPI specifications,
and we need some additional information in order to be able to generate correct
Go code.

An external reference looks like this:

    $ref: ./some_spec.yaml#/components/schemas/Type

We assume that you have already generated the boilerplate code for `./some_spec.yaml`
using `goapi-gen`, and you have a package which contains the generated code,
let's call it `github.com/discord-gophers/some-package`. You need to tell `goapi-gen` that
`some_spec.yaml` corresponds to this package, and you would do it by specifying
this command line argument:

    -import-mapping=./some_spec.yaml:github.com/discord-gophers/some-package

This tells us that in order to resolve references generated from `some_spec.yaml` we
need to import `github.com/discord-gophers/some-package`. You may specify multiple mappings
by comma separating them in the form `key1:value1,key2:value2`.

## What's missing or incomplete

This code is still young, and not complete, since we're filling it in as we
need it. We've not yet implemented several things:

- `oneOf`, `anyOf` are not supported with strong Go typing. This schema:

        schema:
          oneOf:
            - $ref: '#/components/schemas/Cat'
            - $ref: '#/components/schemas/Dog'

  will result in a Go type of `interface{}`. It will be up to you
  to validate whether it conforms to `Cat` and/or `Dog`, depending on the
  keyword. It's not clear if we can do anything much better here given the
  limits of Go typing.

  `allOf` is supported, by taking the union of all the fields in all the
  component schemas. This is the most useful of these operations, and is
  commonly used to merge objects with an identifier, as in the
  `petstore-expanded` example.

- `patternProperties` isn't yet supported and will exit with an error. Pattern
  properties were defined in JSONSchema, and the `kin-openapi` Swagger object
  knows how to parse them, but they're not part of OpenAPI 3.0, so we've left
  them out, as support is very complicated.

## Making changes to code generation

After updating any files under the `pkg/codegen/templates` directory, run `go generate ./...`, and the templates will be updated accordingly.

Alternatively, you can provide custom templates to override built-in ones using
the `-templates` flag specifying a path to a directory containing templates
files. These files **must** be named identically to built-in template files
(see `pkg/codegen/templates/*.tmpl` in the source code), and will be interpreted
on-the-fly at run time. Example:

    $ ls -1 my-templates/
    typedef.tmpl
    $ goapi-gen \
        -templates my-templates/ \
        -generate types \
        petstore-expanded.yaml
