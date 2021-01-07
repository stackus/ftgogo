package logging

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level string

const (
	TRACE Level = "TRACE"
	DEBUG Level = "DEBUG"
	INFO  Level = "INFO"
	WARN  Level = "WARN"
	ERROR Level = "ERROR"
	PANIC Level = "PANIC"
)

//type EnvironmentConfig = func(options []zap.Option) zap.Config
//
type Config struct {
	Environment string
	LogLevel    Level
	//	Options     []zap.Option
}

func NewZapLogger(cfg Config) (*zap.Logger, error) {
	var config zap.Config
	var options []zap.Option

	switch cfg.Environment {
	case "production":
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	default:
		config = zap.NewDevelopmentConfig()
		config.OutputPaths = []string{"stdout"}

		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	config.Level = zap.NewAtomicLevelAt(logLevelToZap(cfg.LogLevel))

	return config.Build(options...)
}

func NewZeroLogger(cfg Config) (zerolog.Logger, error) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	switch cfg.Environment {
	case "production":
		return zerolog.New(os.Stdout).
			Level(logLevelToZero(cfg.LogLevel)).
			With().
			Timestamp().
			Caller().
			Logger(), nil
	default:
		return zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = "03:04:05.000PM"
		})).
			Level(logLevelToZero(cfg.LogLevel)).
			With().
			Timestamp().
			Caller().
			Logger(), nil
	}
}

func logLevelToZap(logLevel Level) zapcore.Level {
	switch logLevel {
	case PANIC:
		return zapcore.PanicLevel
	case ERROR:
		return zapcore.ErrorLevel
	case WARN:
		return zapcore.WarnLevel
	case INFO:
		return zapcore.InfoLevel
	case DEBUG, TRACE:
		return zapcore.DebugLevel
	default:
		return zapcore.InfoLevel
	}
}

func logLevelToZero(level Level) zerolog.Level {
	switch level {
	case PANIC:
		return zerolog.PanicLevel
	case ERROR:
		return zerolog.ErrorLevel
	case WARN:
		return zerolog.WarnLevel
	case INFO:
		return zerolog.InfoLevel
	case DEBUG:
		return zerolog.DebugLevel
	case TRACE:
		return zerolog.TraceLevel
	default:
		return zerolog.InfoLevel
	}
}
