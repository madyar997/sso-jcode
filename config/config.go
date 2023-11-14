package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type (
	// Config -.
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
		PG   `yaml:"postgres"`
		Auth `yaml:"auth"`
		Jwt  `yaml:"jwt"`
		Grpc `yaml:"grpc"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Grpc struct {
		Port string `env-required:"true" yaml:"port" env:"GRPC_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level" mapstructure:"log_level"  env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"true" yaml:"url"      env:"PG_URL"`
	}

	Auth struct {
		Login    string `mapstructure:"login"`
		Password string `mapstructure:"pass"`
	}

	Jwt struct {
		SecretKey       string `mapstructure:"secret_key"`
		AccessTokenTTL  int64  `mapstructure:"access_token_ttl"`
		RefreshTokenTTL int64  `mapstructure:"refresh_token_ttl"`
	}
)

func NewViperConfig() (*Config, error) {
	cfg := Config{}

	viper.SetConfigName("config")    // name of config file (without extension)
	viper.SetConfigType("yml")       // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./config/") // path to look for the config file in
	err := viper.ReadInConfig()      // Find and read the config file
	if err != nil {                  // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
