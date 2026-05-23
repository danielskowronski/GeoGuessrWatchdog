package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const (
	CONF_PATH_ENV_VAR = "GGWD_CONFIG_PATH"
)

type DatabaseConfig struct {
	URL string `koanf:"url"`
}

func LoadConfig[T any](defaultPath string, defaults map[string]any) (T, error) {
	configPath := os.Getenv(CONF_PATH_ENV_VAR)

	cfg, err := Load[T](configPath, defaults)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to load config: %w", err)
	}
	return *cfg, nil
}

func Load[T any](pathsCommaSeparated string, defaults map[string]any) (*T, error) {
	k := koanf.New(".")

	if err := k.Load(confmap.Provider(defaults, "."), nil); err != nil {
		return nil, err
	}

	if err := k.Load(confmap.Provider(map[string]any{}, "."), nil); err != nil {
		return nil, err
	}

	if pathsCommaSeparated != "" {
		paths := strings.Split(pathsCommaSeparated, ",")
		for _, path := range paths {
			path = strings.TrimSpace(path)
			fmt.Printf("loading config from %s ...\n", path)
			if path != "" {
				if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
					return nil, err
				}
			}
		}
	}

	if err := k.Load(env.Provider(".", env.Opt{
		Prefix: ENV_VAR_PREFIX,
		TransformFunc: func(key string, value string) (string, any) {
			key = strings.TrimPrefix(key, ENV_VAR_PREFIX)
			key = strings.ReplaceAll(key, "_", ".")

			if key == "watchdogs.CompetitiveMaps.notifyAbout" || key == "watchdogs.UserStats.observeUsers" || key == "notifierApi.shoutrrr" {
				// split by comma for these list values
				parts := strings.Split(value, ",")
				values := make([]string, 0, len(parts))
				for _, part := range parts {
					part = strings.TrimSpace(part)
					if part != "" {
						values = append(values, part)
					}
				}
				return key, values
			}

			return key, value
		},
	}), nil); err != nil {
		return nil, err
	}

	var cfg T
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
