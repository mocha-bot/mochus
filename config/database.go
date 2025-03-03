package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     uint   `env:"DB_PORT" envDefault:"3306"`
	Username string `env:"DB_USERNAME" envDefault:"root"`
	Password string `env:"DB_PASSWORD" envDefault:""`
	Schema   string `env:"DB_SCHEMA" envDefault:""`
	Dialect  string `env:"DB_DIALECT" envDefault:"mysql"`
}

const (
	MySQLDialect    = "mysql"
	PostgresDialect = "postgres"
)

func (d DatabaseConfig) GetDSN() gorm.Dialector {
	switch d.Dialect {
	case string(MySQLDialect):
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			d.Username,
			d.Password,
			d.Host,
			d.Port,
			d.Schema,
		)
		return mysql.Open(dsn)
	case string(PostgresDialect):
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
			d.Host,
			d.Username,
			d.Password,
			d.Schema,
			d.Port,
		)
		return postgres.Open(dsn)
	default:
		return nil
	}
}
