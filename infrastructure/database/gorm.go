package database

import (
	"fmt"

	"github.com/mocha-bot/mochus/config"
	zLog "github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GORMDatabase struct {
	DatabaseClient *gorm.DB
	config         *config.DatabaseConfig
}

func NewGORMDatabase(cfg *config.DatabaseConfig) (*GORMDatabase, error) {
	dialector := cfg.GetDSN()
	if dialector == nil {
		return nil, fmt.Errorf("unsupported dialect: %s", cfg.Dialect)
	}

	db, err := gorm.Open(dialector, &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	zLog.Info().Msgf("Connected to %s database at %s:%d", cfg.Dialect, cfg.Host, cfg.Port)

	if cfg.Debug {
		db = db.Debug()
		zLog.Info().Msg("Database debug mode enabled")
	}

	return &GORMDatabase{DatabaseClient: db, config: cfg}, nil
}

func (g *GORMDatabase) Close() error {
	sqlDB, err := g.DatabaseClient.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
