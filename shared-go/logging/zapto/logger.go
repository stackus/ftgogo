package zapto

import (
	"go.uber.org/zap"

	"github.com/stackus/edat/log"
)

type zapLogger struct {
	l *zap.Logger
}

func Logger(logger *zap.Logger) log.Logger {
	return &zapLogger{l: logger.WithOptions(zap.AddCallerSkip(1))}
}

func (l *zapLogger) Trace(msg string, fields ...log.Field) {
	if !l.l.Core().Enabled(zap.DebugLevel) {
		return
	}
	l.l.Debug(msg, l.fields(fields)...)
}

func (l *zapLogger) Debug(msg string, fields ...log.Field) {
	if !l.l.Core().Enabled(zap.DebugLevel) {
		return
	}
	l.l.Debug(msg, l.fields(fields)...)
}

func (l *zapLogger) Info(msg string, fields ...log.Field) {
	if !l.l.Core().Enabled(zap.InfoLevel) {
		return
	}
	l.l.Info(msg, l.fields(fields)...)
}

func (l *zapLogger) Warn(msg string, fields ...log.Field) {
	if !l.l.Core().Enabled(zap.WarnLevel) {
		return
	}
	l.l.Warn(msg, l.fields(fields)...)
}

func (l *zapLogger) Error(msg string, fields ...log.Field) {
	if !l.l.Core().Enabled(zap.ErrorLevel) {
		return
	}
	l.l.Error(msg, l.fields(fields)...)
}

func (l *zapLogger) Sub(fields ...log.Field) log.Logger {
	return &zapLogger{l: l.l.With(l.fields(fields)...)}
}

func (l *zapLogger) fields(fields []log.Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))

	for _, field := range fields {
		switch field.Type {
		case log.StringType:
			zapFields = append(zapFields, zap.String(field.Key, field.String))
		case log.IntType:
			zapFields = append(zapFields, zap.Int(field.Key, field.Int))
		case log.DurationType:
			zapFields = append(zapFields, zap.Duration(field.Key, field.Duration))
		case log.ErrorType:
			zapFields = append(zapFields, zap.NamedError(field.Key, field.Error))
		}
	}

	return zapFields
}
