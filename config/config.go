package config

import "github.com/caarlos0/env/v6"

type Config struct {
	Database DatabaseConfig
	Discord  DiscordConfig
	Redis    RedisConfig
	Logger   LoggerConfig
	App      AppConfig
}

func NewConfig() (Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	return cfg, err
}
