package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

type flagConfig struct {
	PackageName     string
	GenerateTargets *cli.StringSlice
	OutputFile      string
	IncludeTags     *cli.StringSlice
	ExcludeTags     *cli.StringSlice
	TemplatesDir    string
	ImportMapping   *cli.StringSlice
	ExcludeSchemas  *cli.StringSlice
	AliasTypes      bool
	Initialisms     *cli.StringSlice
}

type config struct {
	Package        string            `yaml:"package"`
	Generate       []string          `yaml:"generate"`
	Out            string            `yaml:"output"`
	IncludeTags    []string          `yaml:"include-tags"`
	ExcludeTags    []string          `yaml:"exclude-tags"`
	Templates      string            `yaml:"templates"`
	ImportMapping  map[string]string `yaml:"import-mapping"`
	ExcludeSchemas []string          `yaml:"exclude-schemas"`
	Alias          bool              `yaml:"alias"`
	Initialisms    []string          `yaml:"initialisms"`
}

// parseConfig parses the flags and configuration file (if provided). all
// configuration entries via a file will be overridden via the cli flags.
func parseConfig(c *cli.Context, f *flagConfig) (*config, error) {
	cfg := config{}

	// Load the configuration file first.
	if c.IsSet("config") {
		f, err := os.Open(c.String("config"))
		if err != nil {
			return nil, fmt.Errorf("could not open configuration file: %v", err)
		}
		defer f.Close()

		if err = yaml.NewDecoder(f).Decode(&cfg); err != nil {
			return nil, fmt.Errorf("could not decode yaml configuration: %v", err)
		}
	}

	if cfg.Package == "" || c.IsSet(PackageKey) {
		cfg.Package = f.PackageName
	}
	if cfg.Out == "" || c.IsSet(OutKey) {
		cfg.Out = f.OutputFile
	}
	if cfg.Generate == nil || c.IsSet(GenerateKey) {
		cfg.Generate = splitString(f.GenerateTargets, ',')
	}
	if cfg.IncludeTags == nil || c.IsSet(IncludeTagsKey) {
		cfg.IncludeTags = splitString(f.IncludeTags, ',')
	}
	if cfg.ExcludeTags == nil || c.IsSet(ExcludeTagsKey) {
		cfg.ExcludeTags = splitString(f.ExcludeTags, ',')
	}
	if cfg.Templates == "" || c.IsSet(TemplatesKey) {
		cfg.Templates = f.TemplatesDir
	}
	if cfg.ImportMapping == nil || c.IsSet(ImportMappingKey) {
		mappings, err := parseMappings(f.ImportMapping)
		if err != nil {
			return nil, fmt.Errorf("could not parse import mappings: %v", err)
		}
		cfg.ImportMapping = mappings
	}
	if cfg.ExcludeSchemas == nil || c.IsSet(ExcludeSchemasKey) {
		cfg.ExcludeSchemas = splitString(f.ExcludeSchemas, ',')
	}
	if c.IsSet(AliasKey) {
		cfg.Alias = f.AliasTypes
	}
	if cfg.Initialisms == nil || c.IsSet(InitialismsKey) {
		cfg.Initialisms = splitString(f.Initialisms, ',')
	}

	return &cfg, nil
}

func parseTemplateOverrides(templatesDir string) (map[string]string, error) {
	templates := make(map[string]string)

	if templatesDir == "" {
		return templates, nil
	}

	dir, err := os.ReadDir(templatesDir)
	if err != nil {
		return nil, fmt.Errorf("could not open directory: %v", err)
	}

	for _, f := range dir {
		data, err := os.ReadFile(path.Join(templatesDir, f.Name()))
		if err != nil {
			return nil, fmt.Errorf("could not open file: %v", err)
		}
		templates[f.Name()] = string(data)
	}

	return templates, nil
}

func parseMappings(slice *cli.StringSlice) (map[string]string, error) {
	if slice == nil {
		return nil, nil
	}

	result := make(map[string]string)
	mappings := slice.Value()
	for _, t := range mappings {
		kv := strings.Split(t, ":")
		if len(kv) != 2 {
			return nil, fmt.Errorf("expected key:value, got: %q", t)
		}
		result[kv[0]] = kv[1]
	}
	return result, nil
}

func parseSwagger(in io.Reader) (swagger *openapi3.T, err error) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	buf, err := io.ReadAll(in)
	if err != nil {
		return nil, fmt.Errorf("could not read: %v", err)
	}

	return loader.LoadFromData(buf)
}

// This function splits a string along the specifed separator, but it
// ignores anything between double quotes for splitting. We do simple
// inside/outside quote counting. Quotes are not stripped from output.
func splitString(slice *cli.StringSlice, sep rune) []string {
	const escapeChar rune = '"'

	var parts []string
	var part string
	inQuotes := false

	if slice == nil {
		return nil
	}
	for _, s := range slice.Value() {
		for _, c := range s {
			if c == escapeChar {
				if inQuotes {
					inQuotes = false
				} else {
					inQuotes = true
				}
			}

			// If we've gotten the separator rune, consider the previous part
			// complete, but only if we're outside of quoted sections
			if c == sep && !inQuotes {
				parts = append(parts, part)
				part = ""
				continue
			}
			part = part + string(c)
		}
		parts = append(parts, part)
		part = ""
	}
	return parts
}
