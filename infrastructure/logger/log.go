package infrastructure_logger

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/mocha-bot/mochus/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type LoggerOptions struct {
	WithConsole  bool
	WithFile     bool
	LoggerConfig *config.LoggerConfig
}

type LoggerOption func(*LoggerOptions)

func WithConsole(isEnabled bool) LoggerOption {
	return func(o *LoggerOptions) {
		o.WithConsole = isEnabled
	}
}

func WithFile(isEnabled bool) LoggerOption {
	return func(o *LoggerOptions) {
		o.WithFile = isEnabled
	}
}

func WithLoggerConfig(conf *config.LoggerConfig) LoggerOption {
	return func(o *LoggerOptions) {
		o.LoggerConfig = conf
	}
}

func DefaultLoggerOptions() *LoggerOptions {
	return &LoggerOptions{
		WithConsole:  true,
		WithFile:     false,
		LoggerConfig: nil,
	}
}

var once = new(sync.Once)

func NewLogger(opts ...LoggerOption) zerolog.Logger {
	options := DefaultLoggerOptions()
	for _, opt := range opts {
		opt(options)
	}

	var log zerolog.Logger

	once.Do(func() {
		logWriters := make([]io.Writer, 0)

		if options.WithConsole {
			logWriters = append(logWriters, NewConsoleLogger())
		}

		if options.WithFile && options.LoggerConfig != nil {
			logWriters = append(logWriters, NewLumberjackFileLogger(options.LoggerConfig.ToLumberjackFileConfig()))
		}

		zerolog.TimestampFunc = time.Now
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

		mw := zerolog.MultiLevelWriter(logWriters...)

		log = zerolog.
			New(mw).
			With().
			Timestamp().
			Logger()
	})

	return log
}

func NewConsoleLogger() io.Writer {
	return zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    true,
		TimeFormat: time.RFC3339,
	}
}

func NewLumberjackFileLogger(conf config.LumberjackFileConfig) io.Writer {
	return &lumberjack.Logger{
		Filename:   conf.Filename,
		MaxBackups: conf.MaxBackups, // files
		MaxSize:    conf.MaxSize,    // megabytes
		MaxAge:     conf.MaxAge,     // days
		LocalTime:  conf.LocalTime,
		Compress:   conf.Compress,
	}
}
