package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const DefaultFieldName = "-"

type Fields map[string]interface{}

type LogEntry struct {
	infoSugared  *zap.SugaredLogger
	errorSugared *zap.SugaredLogger
	infoLogger   *zap.Logger
	errorLogger  *zap.Logger
}

func (le *LogEntry) ContextWithLogger(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey, le)
}

func getLogEntry(infoLogger *zap.Logger, errorLogger *zap.Logger) *LogEntry {
	return &LogEntry{
		infoLogger:   infoLogger,
		errorLogger:  errorLogger,
		infoSugared:  infoLogger.Sugar(),
		errorSugared: errorLogger.Sugar(),
	}
}

func newLogEntry(logEntry *LogEntry, fields Fields) *LogEntry {
	args := convertFields(fields)

	le := &LogEntry{
		infoLogger:  logEntry.infoLogger.With(args...),
		errorLogger: logEntry.errorLogger.With(args...),
	}

	le.infoSugared = le.infoLogger.Sugar()
	le.errorSugared = le.errorLogger.Sugar()

	return le
}

func convertFields(fields Fields) []zapcore.Field {
	zfields := make([]zapcore.Field, 0, len(fields))
	for k, v := range fields {
		zfields = append(zfields, zap.Any(k, v))
	}
	return zfields
}

func (le *LogEntry) WithFields(f Fields) *LogEntry {
	args := convertFields(f)
	l := &LogEntry{
		infoLogger:  le.infoLogger.With(args...),
		errorLogger: le.errorLogger.With(args...),
	}
	l.infoSugared = l.infoLogger.Sugar()
	l.errorSugared = l.errorLogger.Sugar()
	return l
}

func (le *LogEntry) DebugWith(msg string, fields Fields) {
	le.infoLogger.Debug(msg, convertFields(fields)...)
}

func (le *LogEntry) Debug(msg string) {
	le.infoLogger.Debug(msg)

}

func (le *LogEntry) Debugf(template string, args ...interface{}) {
	le.infoSugared.Debugf(template, args)
}

func (le *LogEntry) Debugln(args ...interface{}) {
	le.infoSugared.Debugln(args)
}

func (le *LogEntry) Debugw(msg string, keysAndValues ...interface{}) {
	le.infoSugared.Debugw(msg, keysAndValues...)
}

func (le *LogEntry) Debugv(msg string, fields ...zapcore.Field) {
	le.infoLogger.Debug(msg, fields...)
}

func (le *LogEntry) InfoWith(msg string, fields Fields) {
	le.infoLogger.Info(msg, convertFields(fields)...)
}

func (le *LogEntry) Info(msg string) {
	le.infoLogger.Info(msg)

}

func (le *LogEntry) Infof(template string, args ...interface{}) {
	le.infoSugared.Infof(template, args)
}

func (le *LogEntry) Infoln(args ...interface{}) {
	le.infoSugared.Infoln(args)
}

func (le *LogEntry) Infow(msg string, keysAndValues ...interface{}) {
	le.infoSugared.Infow(msg, keysAndValues...)
}

func (le *LogEntry) Infov(msg string, fields ...zapcore.Field) {
	le.infoLogger.Info(msg, fields...)
}

func (le *LogEntry) WarnWith(msg string, fields Fields) {
	le.errorLogger.Warn(msg, convertFields(fields)...)
}

func (le *LogEntry) Warn(msg string) {
	le.errorLogger.Warn(msg)

}

func (le *LogEntry) Warnf(template string, args ...interface{}) {
	le.errorSugared.Warnf(template, args)
}

func (le *LogEntry) Warnln(args ...interface{}) {
	le.errorSugared.Warnln(args)
}

func (le *LogEntry) Warnw(msg string, keysAndValues ...interface{}) {
	le.errorSugared.Warnw(msg, keysAndValues...)
}

func (le *LogEntry) Warnv(msg string, fields ...zapcore.Field) {
	le.errorLogger.Warn(msg, fields...)
}

func (le *LogEntry) ErrorWith(msg string, fields Fields) {
	le.errorLogger.Error(msg, convertFields(fields)...)
}

func (le *LogEntry) Error(msg string) {
	le.errorLogger.Error(msg)

}

func (le *LogEntry) Errorf(template string, args ...interface{}) {
	le.errorSugared.Errorf(template, args)
}

func (le *LogEntry) Errorln(args ...interface{}) {
	le.errorSugared.Errorln(args)
}

func (le *LogEntry) Errorw(msg string, keysAndValues ...interface{}) {
	le.errorSugared.Errorw(msg, keysAndValues...)
}

func (le *LogEntry) Errorv(msg string, fields ...zapcore.Field) {
	le.errorLogger.Error(msg, fields...)
}

func (le *LogEntry) FatalWith(msg string, fields Fields) {
	le.errorLogger.Fatal(msg, convertFields(fields)...)
}

func (le *LogEntry) Fatal(msg string) {
	le.errorLogger.Fatal(msg)

}

func (le *LogEntry) Fatalf(template string, args ...interface{}) {
	le.errorSugared.Fatalf(template, args)
}

func (le *LogEntry) Fatalln(args ...interface{}) {
	le.errorSugared.Fatalln(args)
}

func (le *LogEntry) Fatalw(msg string, keysAndValues ...interface{}) {
	le.errorSugared.Fatalw(msg, keysAndValues...)
}

func (le *LogEntry) Fatalv(msg string, fields ...zapcore.Field) {
	le.errorLogger.Fatal(msg, fields...)
}

func (le *LogEntry) PanicWith(msg string, fields Fields) {
	le.errorLogger.Panic(msg, convertFields(fields)...)
}

func (le *LogEntry) Panic(msg string) {
	le.errorLogger.Panic(msg)

}

func (le *LogEntry) Panicf(template string, args ...interface{}) {
	le.errorSugared.Panicf(template, args)
}

func (le *LogEntry) Panicln(args ...interface{}) {
	le.errorSugared.Panicln(args)
}

func (le *LogEntry) Panicw(msg string, keysAndValues ...interface{}) {
	le.errorSugared.Panicw(msg, keysAndValues...)
}

func (le *LogEntry) Panicv(msg string, fields ...zapcore.Field) {
	le.errorLogger.Panic(msg, fields...)
}

func (le *LogEntry) DPanicWith(msg string, fields Fields) {
	le.errorLogger.DPanic(msg, convertFields(fields)...)
}

func (le *LogEntry) DPanic(msg string) {
	le.errorLogger.DPanic(msg)

}

func (le *LogEntry) DPanicf(template string, args ...interface{}) {
	le.errorSugared.DPanicf(template, args)
}

func (le *LogEntry) DPanicln(args ...interface{}) {
	le.errorSugared.DPanicln(args)
}

func (le *LogEntry) DPanicw(msg string, keysAndValues ...interface{}) {
	le.errorSugared.DPanicw(msg, keysAndValues...)
}

func (le *LogEntry) DPanicv(msg string, fields ...zapcore.Field) {
	le.errorLogger.DPanic(msg, fields...)
}
