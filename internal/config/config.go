package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	EnvKeyPrefix = "OBS_PLATFORM"
)

type ServiceSettings struct {
	Host string       `mapstructure:"host"`
	Port int          `mapstructure:"port"`
	Auth AuthSettings `mapstructure:"auth"`
}

type AuthSettings struct {
	Enabled bool     `mapstructure:"enabled"`
	APIKeys []string `mapstructure:"api_keys"`
}

// Load reads configuration from environment variables (prefix: OBS_PLATFORM_)
// with defaults of 0.0.0.0:8080 and auth disabled.
func Load() (ServiceSettings, error) {
	defaults := map[string]any{
		"host":         "0.0.0.0",
		"port":         8080,
		"auth.enabled": false,
	}

	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetEnvPrefix(EnvKeyPrefix)
	v.AutomaticEnv()

	var cfg ServiceSettings
	if err := v.Unmarshal(&cfg); err != nil {
		return ServiceSettings{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}
