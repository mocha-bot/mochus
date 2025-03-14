package config

import "fmt"

type AppConfig struct {
	Name     string `env:"APP_NAME" envDefault:"mochus"`
	Host     string `env:"APP_HOST" envDefault:"localhost"`
	Port     string `env:"APP_PORT" envDefault:"8083"`
	Timezone string `env:"APP_TIMEZONE" envDefault:"Asia/Jakarta"`
	Debug    bool   `env:"APP_DEBUG" envDefault:"true"`
	Gateway  string `env:"APP_GATEWAY_URL" envDefault:"http://localhost:3000"`
}

func (a AppConfig) GetAddress() string {
	return fmt.Sprintf("%s:%s", a.Host, a.Port)
}
