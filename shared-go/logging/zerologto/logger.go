package zerologto

import (
	"github.com/rs/zerolog"

	"github.com/stackus/edat/log"
)

type zerologLogger struct {
	l zerolog.Logger
}

func Logger(logger zerolog.Logger) log.Logger {
	zLog := logger.With().CallerWithSkipFrameCount(3).Logger()
	return &zerologLogger{l: zLog}
}

func (l *zerologLogger) Trace(msg string, fields ...log.Field) {
	if l.l.GetLevel() > zerolog.TraceLevel {
		return
	}
	logger := l.fields(l.l.With(), fields).Logger()
	logger.Trace().Msg(msg)
}

func (l *zerologLogger) Debug(msg string, fields ...log.Field) {
	if l.l.GetLevel() > zerolog.DebugLevel {
		return
	}
	logger := l.fields(l.l.With(), fields).Logger()
	logger.Debug().Msg(msg)
}

func (l *zerologLogger) Info(msg string, fields ...log.Field) {
	if l.l.GetLevel() > zerolog.InfoLevel {
		return
	}
	logger := l.fields(l.l.With(), fields).Logger()
	logger.Info().Msg(msg)
}

func (l *zerologLogger) Warn(msg string, fields ...log.Field) {
	if l.l.GetLevel() > zerolog.WarnLevel {
		return
	}
	logger := l.fields(l.l.With(), fields).Logger()
	logger.Warn().Msg(msg)
}

func (l *zerologLogger) Error(msg string, fields ...log.Field) {
	if l.l.GetLevel() > zerolog.ErrorLevel {
		return
	}
	logger := l.fields(l.l.With(), fields).Logger()
	logger.Error().Msg(msg)
}

func (l *zerologLogger) Sub(fields ...log.Field) log.Logger {
	return &zerologLogger{l: l.fields(l.l.With(), fields).Logger()}
}

func (l *zerologLogger) fields(ctx zerolog.Context, fields []log.Field) zerolog.Context {
	for _, field := range fields {
		switch field.Type {
		case log.StringType:
			ctx = ctx.Str(field.Key, field.String)
		case log.IntType:
			ctx = ctx.Int(field.Key, field.Int)
		case log.DurationType:
			ctx = ctx.Dur(field.Key, field.Duration)
		case log.ErrorType:
			ctx = ctx.Stack().Err(field.Error)
		}
	}
	return ctx
}
