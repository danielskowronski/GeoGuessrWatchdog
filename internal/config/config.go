package config

import (
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Database     DatabaseConfig     `koanf:"database"`
	Temporal     TemporalConfig     `koanf:"temporal"`
	GeoguessrAPI GeoguessrAPIConfig `koanf:"geoguessrApi"`
	NotifierAPI  NotifierAPIConfig  `koanf:"notifierApi"`
	Watchdogs    WatchdogsConfig    `koanf:"watchdogs"`
}

type DatabaseConfig struct {
	URL string `koanf:"url"`
}

type TemporalConfig struct {
	Address   string `koanf:"address"`
	Namespace string `koanf:"namespace"`
	TaskQueue string `koanf:"taskQueue"`
}

type GeoguessrAPIConfig struct {
	BaseURL        string `koanf:"baseUrl"`
	Cookie         string `koanf:"cookie"`
	Proxy          string `koanf:"proxy"`
	TimeoutSeconds int    `koanf:"timeoutSeconds"`
}

type NotifierAPIConfig struct {
	Target  *string       `koanf:"target"`
	Discord DiscordConfig `koanf:"discord"`
}

type DiscordConfig struct {
	BotToken  string   `koanf:"botToken"`
	Receivers []string `koanf:"receivers"`
}

type WatchdogsConfig struct {
	CompetitiveMaps WatchdogCompetitiveMaps `koanf:"CompetitiveMaps"`
	UserStats       WatchdogUserStats       `koanf:"UserStats"`
}

type WatchdogCompetitiveMaps struct {
	Enabled                bool                            `koanf:"enabled"`
	ScheduleFrequencyHours int                             `koanf:"scheduleFrequencyHours"`
	NotifyAbout            []NotifyAboutDivisionModeConfig `koanf:"notifyAbout"`
	Temporal               TemporalAdvanced                `koanf:"temporalAdvanced"`
}

type ScheduleDailyConfig struct {
	Hour   int `koanf:"hour"`
	Minute int `koanf:"minute"`
}

type NotifyAboutDivisionModeConfig struct {
	DivisionName string `koanf:"divisionName"`
	GameMode     string `koanf:"gameMode"`
}

type TemporalAdvanced struct {
	WorkflowName   string         `koanf:"workflowName"`
	FanoutActivity ActivityConfig `koanf:"fanoutActivity"`
	Parallelism    int            `koanf:"parallelism"`
}

type WatchdogUserStats struct {
	Enabled       bool                `koanf:"enabled"`
	ScheduleDaily ScheduleDailyConfig `koanf:"scheduleDaily"`
	ObserveUsers  []string            `koanf:"observeUsers"`
	Temporal      TemporalAdvanced    `koanf:"temporalAdvanced"`
}

type ActivityConfig struct {
	StartToCloseTimeoutSeconds  int     `koanf:"startToCloseTimeoutSeconds"`
	RetryInitialIntervalSeconds int     `koanf:"retryInitialIntervalSeconds"`
	RetryBackoffCoefficient     float64 `koanf:"retryBackoffCoefficient"`
	RetryMaximumIntervalSeconds int     `koanf:"retryMaximumIntervalSeconds"`
	RetryMaximumAttempts        int     `koanf:"retryMaximumAttempts"`
}

const (
	ENV_VAR_PREFIX = "GGWD_"

	DefaultTemporalOptions_StartToCloseTimeoutSeconds  = 60
	DefaultTemporalOptions_RetryInitialIntervalSeconds = 5
	DefaultTemporalOptions_RetryBackoffCoefficient     = 2.0
	DefaultTemporalOptions_RetryMaximumIntervalSeconds = 60
	DefaultTemporalOptions_RetryMaximumAttempts        = 3
	DefaultTemporalOptions_Parallelism                 = 2
)

func Load(path string) (*Config, error) {
	k := koanf.New(".")

	if err := k.Load(confmap.Provider(map[string]any{
		"temporal.address":   "localhost:7233",
		"temporal.namespace": "default",
		"temporal.taskQueue": "ggwd-task-queue",

		"geoguessrApi.baseUrl":        "https://www.geoguessr.com/api",
		"geoguessrApi.proxy":          "",
		"geoguessrApi.timeoutSeconds": 30,

		"notifierApi.target": "",

		"watchdogs.CompetitiveMaps.enabled":                       true,
		"watchdogs.CompetitiveMaps.scheduleFrequencyHours":        6,
		"watchdogs.CompetitiveMaps.temporalAdvanced.workflowName": "GGWDCompetitiveMapsWorkflow",
		"watchdogs.CompetitiveMaps.temporalAdvanced.parallelism":  DefaultTemporalOptions_Parallelism,

		"watchdogs.CompetitiveMaps.temporalAdvanced.fanoutActivity.startToCloseTimeoutSeconds":  DefaultTemporalOptions_StartToCloseTimeoutSeconds,
		"watchdogs.CompetitiveMaps.temporalAdvanced.fanoutActivity.retryInitialIntervalSeconds": DefaultTemporalOptions_RetryInitialIntervalSeconds,
		"watchdogs.CompetitiveMaps.temporalAdvanced.fanoutActivity.retryBackoffCoefficient":     DefaultTemporalOptions_RetryBackoffCoefficient,
		"watchdogs.CompetitiveMaps.temporalAdvanced.fanoutActivity.retryMaximumIntervalSeconds": DefaultTemporalOptions_RetryMaximumIntervalSeconds,
		"watchdogs.CompetitiveMaps.temporalAdvanced.fanoutActivity.retryMaximumAttempts":        DefaultTemporalOptions_RetryMaximumAttempts,

		"watchdogs.UserStats.enabled":                       true,
		"watchdogs.UserStats.scheduleDaily.hour":            0,
		"watchdogs.UserStats.scheduleDaily.minute":          0,
		"watchdogs.UserStats.temporalAdvanced.workflowName": "GGWDUserStatsWorkflow",
		"watchdogs.UserStats.temporalAdvanced.parallelism":  DefaultTemporalOptions_Parallelism,

		"watchdogs.UserStats.temporalAdvanced.fanoutActivity.startToCloseTimeoutSeconds":  DefaultTemporalOptions_StartToCloseTimeoutSeconds,
		"watchdogs.UserStats.temporalAdvanced.fanoutActivity.retryInitialIntervalSeconds": DefaultTemporalOptions_RetryInitialIntervalSeconds,
		"watchdogs.UserStats.temporalAdvanced.fanoutActivity.retryBackoffCoefficient":     DefaultTemporalOptions_RetryBackoffCoefficient,
		"watchdogs.UserStats.temporalAdvanced.fanoutActivity.retryMaximumIntervalSeconds": DefaultTemporalOptions_RetryMaximumIntervalSeconds,
		"watchdogs.UserStats.temporalAdvanced.fanoutActivity.retryMaximumAttempts":        DefaultTemporalOptions_RetryMaximumAttempts,
	}, "."), nil); err != nil {
		return nil, err
	}

	if path != "" {
		if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
			return nil, err
		}
	}

	if err := k.Load(env.Provider(".", env.Opt{
		Prefix: ENV_VAR_PREFIX,
		TransformFunc: func(key string, value string) (string, any) {
			key = strings.TrimPrefix(key, ENV_VAR_PREFIX)
			key = strings.ReplaceAll(key, "_", ".")

			if key == "watchdogs.CompetitiveMaps.notifyAbout" || key == "watchdogs.UserStats.observeUsers" || key == "notifierApi.discord.receivers" {
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

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
