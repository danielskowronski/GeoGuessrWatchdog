package config

type ApiConfig struct {
	Database    DatabaseConfig    `koanf:"database"`
	Server      ApiServerConfig   `koanf:"server"`
	UserAliases UserAliasesConfig `koanf:"userAliases"`
}

type ApiServerConfig struct {
	Bind string `koanf:"bind"`
}

type UserAliasesConfig map[string]string

func ApiConfigDefaults() map[string]any {
	return map[string]any{
		"server.bind": ":8090",
	}
}
