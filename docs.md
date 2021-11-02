# NAME

goapi-gen - Generate Go code from OpenAPI specification YAML

# SYNOPSIS

goapi-gen

```
[--alias|-a]
[--config|-c]=[value]
[--exclude-schemas|-S]=[value]
[--exclude-tags|-T]=[value]
[--generate|-g]=[value]
[--help|-h]
[--import-mapping|-i]=[value]
[--include-tags|-t]=[value]
[--initialisms]=[value]
[--out|-o]=[value]
[--package|-p]=[value]
[--templates|-s]=[value]
[--version|-v]
```

**Usage**:

```
goapi-gen [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--alias, -a**: Alias type declerations when possible

**--config, -c**="": Read configuration from a config file

**--exclude-schemas, -S**="": Exclude matching schemas from generation (default: [])

**--exclude-tags, -T**="": Exclude matching operations in the given tags (default: [])

**--generate, -g**="": List of generation options. (default: [types server spec])

**--help, -h**: show help

**--import-mapping, -i**="": A dict from the external reference to golang package path (default: [])

**--include-tags, -t**="": Only include matching operations in the given tags. (default: [])

**--initialisms**="": Add custom initialisms (i.e ID, API, URI) (default: [])

**--out, -o**="": Output file

**--package, -p**="": The package name for generated code.

**--templates, -s**="": Generate templates from a different directory

**--version, -v**: print the version


# COMMANDS

## list

list available generation options

## help, h

Shows a list of commands or help for one command

