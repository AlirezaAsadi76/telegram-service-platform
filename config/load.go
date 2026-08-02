package config

import (
	"os"

	"strings"

	"slices"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func Load(configPath string) Config {
	var cfg Config

	koans := koanf.New(".")

	_ = koans.Load(confmap.Provider(defaultValue, "."), nil)

	if err := koans.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		panic(err)
	}
	_ = koans.Load(env.Provider(".", env.Opt{
		Prefix: "PLATFORM_",
		TransformFunc: func(k, v string) (string, any) {
			k = strings.ReplaceAll(strings.ToLower(
				strings.TrimPrefix(k, "PLATFORM_")), "_", ".")
			k = strings.ReplaceAll(k, "..", "_")

			// If there is a space in the value, split the value into a slice by the space.
			if strings.Contains(v, " ") {
				return k, strings.Split(v, " ")
			}
			return k, v
		},
		EnvironFunc: func() []string {
			return slices.DeleteFunc(os.Environ(), func(s string) bool {
				return strings.HasPrefix(s, "PLATFORM_")
			})
		},
	}), nil)

	err := koans.Unmarshal("", &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
