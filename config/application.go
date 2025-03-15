package config

import "fmt"

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
)

type AppConfig struct {
	Name        string `env:"APP_NAME" envDefault:"mochus"`
	Host        string `env:"APP_HOST" envDefault:"localhost"`
	Port        string `env:"APP_PORT" envDefault:"8083"`
	Timezone    string `env:"APP_TIMEZONE" envDefault:"Asia/Jakarta"`
	Debug       bool   `env:"APP_DEBUG" envDefault:"true"`
	Gateway     string `env:"APP_GATEWAY_URL" envDefault:"http://127.0.0.1:3000"`
	Environment string `env:"APP_ENV" envDefault:"development"`
}

func (a AppConfig) GetAddress() string {
	return fmt.Sprintf("%s:%s", a.Host, a.Port)
}

func (a AppConfig) IsProduction() bool {
	return a.Environment == string(Production)
}

func (a AppConfig) IsLocalhost() bool {
	return a.Host == "localhost" || a.Host == "127.0.0.1"
}
