// Copyright 2019 DeepMap, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/discord-gophers/goapi-gen/pkg/codegen"
	"github.com/urfave/cli/v2"
)

const (
	PackageKey        = "package"
	GenerateKey       = "generate"
	OutKey            = "out"
	IncludeTagsKey    = "include-tags"
	ExcludeTagsKey    = "exclude-tags"
	TemplatesKey      = "templates"
	ImportMappingKey  = "import-mapping"
	ExcludeSchemasKey = "exclude-schemas"
	AliasKey          = "alias"
	ConfigKey         = "config"
)

func run(c *cli.Context, cfg *config) error {
	if c.Args().Len() == 0 && (cfg.Package == "" || cfg.Package == "-") {
		return errors.New("package required when reading from stdin")
	}

	if cfg.Package == "" {
		path := flag.Arg(0)
		baseName := filepath.Base(path)
		nameParts := strings.Split(baseName, ".")
		cfg.Package = codegen.ToCamelCase(nameParts[0])
	}

	var err error
	in := os.Stdin
	if file := c.Args().Get(0); file != "" {
		in, err = os.Open(file)
		if err != nil {
			return fmt.Errorf("could not not open %s: %v", file, err)
		}
		defer in.Close()
	}

	templates, err := parseTemplateOverrides(cfg.Templates)
	if err != nil {
		return fmt.Errorf("could not open templates: %s", err)
	}

	opts := codegen.Options{
		IncludeTags:    cfg.IncludeTags,
		ExcludeTags:    cfg.ExcludeTags,
		ExcludeSchemas: cfg.ExcludeSchemas,
		UserTemplates:  templates,
		ImportMapping:  cfg.ImportMapping,
	}

	for _, tgt := range cfg.Generate {
		switch tgt {
		case "client":
			opts.GenerateClient = true
		case "server":
			opts.GenerateChiServer = true
		case "types":
			opts.GenerateTypes = true
		case "spec":
			opts.EmbedSpec = true
		case "skip-fmt":
			opts.SkipFmt = true
		case "skip-prune":
			opts.SkipPrune = true
		default:
			return fmt.Errorf("unknown generation option: %s", tgt)
		}
	}

	swagger, err := parseSwagger(in)
	if err != nil {
		return fmt.Errorf("could not load spec: %v", err)
	}
	code, err := codegen.Generate(swagger, cfg.Package, opts)
	if err != nil {
		return fmt.Errorf("could not generate code: %v", err)
	}

	out := os.Stdout
	if cfg.Out != "" {
		out, err = os.Create(cfg.Out)
		if err != nil {
			return fmt.Errorf("could not open output file: %v", err)
		}
		defer out.Close()
	}

	_, err = out.WriteString(code)
	if err != nil {
		return fmt.Errorf("could not write code: %v", err)
	}

	return nil
}

func main() {
	f := &flagConfig{
		GenerateTargets: cli.NewStringSlice("types", "client", "server", "spec"),
		IncludeTags:     &cli.StringSlice{},
		ExcludeTags:     &cli.StringSlice{},
		ImportMapping:   &cli.StringSlice{},
		ExcludeSchemas:  &cli.StringSlice{},
	}
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        PackageKey,
				Aliases:     []string{"p"},
				Usage:       "The package name for generated code.",
				DefaultText: "swagger file name",
				Destination: &f.PackageName,
			},
			&cli.StringSliceFlag{
				Name:        GenerateKey,
				Aliases:     []string{"g"},
				Value:       cli.NewStringSlice("types", "client", "server", "spec"),
				Usage:       `List of generation options.`,
				DefaultText: "types,client,server,spec",
				Destination: f.GenerateTargets,
			},
			&cli.StringFlag{
				Name:        OutKey,
				Aliases:     []string{"o"},
				Usage:       "Output file",
				DefaultText: "<stdout>",
				Destination: &f.OutputFile,
			},
			&cli.StringSliceFlag{
				Name:        IncludeTagsKey,
				Aliases:     []string{"t"},
				Usage:       "Only include matching operations in the given tags.",
				DefaultText: "<all>",
				Destination: f.IncludeTags,
			},
			&cli.StringSliceFlag{
				Name:        ExcludeTagsKey,
				Aliases:     []string{"T"},
				Usage:       "Exclude matching operations in the given tags",
				DefaultText: "<none>",
				Destination: f.ExcludeTags,
			},
			&cli.StringFlag{
				Name:        TemplatesKey,
				Aliases:     []string{"s"},
				Usage:       "Generate templates from a different directory",
				DefaultText: "<builtin>",
				Destination: &f.TemplatesDir,
			},
			&cli.StringSliceFlag{
				Name:        ImportMappingKey,
				Aliases:     []string{"i"},
				Usage:       "A dict from the external reference to golang package path",
				Destination: f.ImportMapping,
			},
			&cli.StringSliceFlag{
				Name:        ExcludeSchemasKey,
				Aliases:     []string{"S"},
				Usage:       "Exclude matching schemas from generation",
				DefaultText: "<none>",
				Destination: f.ExcludeSchemas,
			},
			&cli.BoolFlag{
				Name:        AliasKey,
				Aliases:     []string{"a"},
				Usage:       "Alias type declerations when possible",
				Destination: &f.AliasTypes,
			},
			&cli.StringFlag{
				Name:        ConfigKey,
				Aliases:     []string{"c"},
				Usage:       "Read configuration from a config file",
				DefaultText: "<none>",
			},
		},
		EnableBashCompletion: true,
		Version:              "v0.0.1-alpha",
		Usage:                "Generate Go code from OpenAPI specification YAML",
		Action: func(c *cli.Context) error {
			cfg, err := parseConfig(c, f)
			if err != nil {
				return fmt.Errorf("could not parse args: %v", err)
			}
			return run(c, cfg)
		},

		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "list available generation options",
				Action: func(_ *cli.Context) error {
					fmt.Println(generateList)
					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// //go:embed list.txt
var generateList string
