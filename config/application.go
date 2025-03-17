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

	// TLS
	TLSCertFile string `env:"APP_CERT_FILE" envDefault:""`
	TLSKeyFile  string `env:"APP_KEY_FILE" envDefault:""`
	TLS         bool   `env:"APP_TLS" envDefault:"false"`

	// CORS
	CORSAllowOrigins     []string `env:"APP_CORS_ALLOW_ORIGINS" envSeparator:"," envDefault:"*.mocha-bot.xyz,mocha-bot.xyz"`
	CORSAllowMethods     []string `env:"APP_CORS_ALLOW_METHODS" envSeparator:"," envDefault:"GET,POST,PUT,DELETE,OPTIONS"`
	CORSAllowCredentials bool     `env:"APP_CORS_ALLOW_CREDENTIALS" envDefault:"true"`
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

func (a AppConfig) IsTLS() bool {
	return a.TLS
}
