package config

import (
	"os"
	"slices"
	"strings"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
)

func Load(configPath string) Config {
	var cfg Config

	// Load .env into process environment.
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	koans := koanf.New(".")

	// 1. Load default values.
	if err := koans.Load(
		confmap.Provider(defaultValue, "."),
		nil,
	); err != nil {
		panic(err)
	}

	// 2. Read config file.
	configData, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	// 3. Resolve environment variables:
	//
	// ${TELEGRAM_BOT_TOKEN}
	// ${JPANAEL_API_KEY}
	//
	// will be replaced with their environment values.
	configData = []byte(
		os.ExpandEnv(string(configData)),
	)

	// 4. Load YAML.
	if err := koans.Load(
		rawbytes.Provider(configData),
		yaml.Parser(),
	); err != nil {
		panic(err)
	}

	// 5. Load PLATFORM_* environment variables.
	//
	// Example:
	//
	// PLATFORM_JUSTANOTHERPANEL_API_KEY=xxx
	//
	// becomes:
	//
	// justanotherpanel.api.key
	//
	if err := koans.Load(
		env.Provider(".", env.Opt{
			Prefix: "PLATFORM_",

			TransformFunc: func(k, v string) (string, any) {
				k = strings.TrimPrefix(k, "PLATFORM_")
				k = strings.ToLower(k)
				k = strings.ReplaceAll(k, "_", ".")

				k = strings.ReplaceAll(k, "..", "_")

				if strings.Contains(v, " ") {
					return k, strings.Split(v, " ")
				}

				return k, v
			},

			EnvironFunc: func() []string {
				return slices.DeleteFunc(
					os.Environ(),
					func(s string) bool {
						return !strings.HasPrefix(s, "PLATFORM_")
					},
				)
			},
		}),
		nil,
	); err != nil {
		panic(err)
	}

	// 6. Unmarshal into Config.
	if err := koans.Unmarshal("", &cfg); err != nil {
		panic(err)
	}

	return cfg
}
