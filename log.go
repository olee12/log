package log

import (
	"context"
	"os"
	"path"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type key int

var loggerKey = key(1)

type Level = zapcore.Level

// DefaultZapLogger is the default logger instance that should be used to log
// It's assigned a default value here for tests (which do not call log.Configure())
var DefaultZapLogger = newZapLogger(defaultConfig, os.Stdout, os.Stderr, true)

const (
	DebugLevel  = zapcore.DebugLevel
	InfoLevel   = zapcore.InfoLevel
	WarnLevel   = zapcore.WarnLevel
	ErrorLevel  = zapcore.ErrorLevel
	DPanicLevel = zapcore.DPanicLevel
	PanicLevel  = zapcore.PanicLevel
	FatalLevel  = zapcore.FatalLevel
)

// Config for logging
type Config struct {
	// Level set log level
	Level zapcore.Level
	// EncodeLogsAsJson makes the log framework log JSON
	EncodeLogsAsJson bool
	// FileLoggingEnabled makes the framework log to a file
	FileLoggingEnabled bool
	// ConsoleLoggingEnabled makes the framework log to console
	ConsoleLoggingEnabled bool
	// CallerEnabled makes the caller log to a file
	CallerEnabled bool
	// CallerSkip increases the number of callers skipped by caller
	CallerSkip int
	// Directory to log to when file logging is enabled
	Directory string
	// Filename is the name of the log file which will be placed inside the directory
	Filename string
	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int
	// MaxBackups the max number of rolled files to keep
	MaxBackups int
	// MaxAge the max age in days to keep a log file
	MaxAge int
	// ConsoleInfoStream
	ConsoleInfoStream *os.File
	// ConsoleErrorStream
	ConsoleErrorStream *os.File
	// ConsoleSeparator the separator of fields of the log record
	ConsoleSeparator string
	// LevelEncoder use lowercase or capital case encoder
	LevelEncoder zapcore.LevelEncoder
}

var (
	loglv zap.AtomicLevel
)

func SetLevel(l Level) {
	loglv.SetLevel(l)
}

func GetLevel() Level {
	return loglv.Level()
}

// ShortTimeEncoder serializes a time.Time to an short-formatted string
func ShortTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02T15:04:05.000"))
}

// ConsoleLogTimeEncoder serializes a time.Time to an short-formatted string
func ConsoleLogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// defaultConfig is used for DefaultZapLogger below only
var defaultConfig = Config{
	Level:            DebugLevel,
	EncodeLogsAsJson: false,
	CallerEnabled:    true,
	CallerSkip:       1,
	ConsoleSeparator: "|",
	LevelEncoder:     zapcore.LowercaseLevelEncoder,
}

// Configure sets up the logging framework
func Configure(config Config) error {
	infoWriters := []zapcore.WriteSyncer{}
	errWriters := []zapcore.WriteSyncer{}

	if config.FileLoggingEnabled {
		infoLog := newRollingFile(config.Directory, getNameByLogLevel(config.Filename, InfoLevel), config.MaxSize, config.MaxAge, config.MaxBackups)
		errLog := newRollingFile(config.Directory, getNameByLogLevel(config.Filename, ErrorLevel), config.MaxSize, config.MaxAge, config.MaxBackups)
		infoWriters = append(infoWriters, infoLog)
		errWriters = append(errWriters, errLog)
	} else {
		config.ConsoleLoggingEnabled = true
	}

	if config.ConsoleLoggingEnabled {
		if config.ConsoleInfoStream != nil {
			infoWriters = append(infoWriters, config.ConsoleInfoStream)
		} else {
			infoWriters = append(infoWriters, os.Stdout)
		}
		if config.ConsoleErrorStream != nil {
			errWriters = append(errWriters, config.ConsoleErrorStream)
		} else {
			errWriters = append(errWriters, os.Stderr)
		}
	}

	DefaultZapLogger = newZapLogger(config, zapcore.NewMultiWriteSyncer(infoWriters...), zapcore.NewMultiWriteSyncer(errWriters...), true)

	DeclareLogger(config, Infov)
	DeclareLogger(config, Errorv)

	return nil
}

// NewLogEntry create a new logentry instead of override defaultzaplogger
func NewLogEntry(config Config) *LogEntry {
	infoWriters := []zapcore.WriteSyncer{}
	errWriters := []zapcore.WriteSyncer{}

	if config.FileLoggingEnabled {
		infoLog := newRollingFile(config.Directory, getNameByLogLevel(config.Filename, InfoLevel), config.MaxSize, config.MaxAge, config.MaxBackups)
		errLog := newRollingFile(config.Directory, getNameByLogLevel(config.Filename, ErrorLevel), config.MaxSize, config.MaxAge, config.MaxBackups)
		infoWriters = append(infoWriters, infoLog)
		errWriters = append(errWriters, errLog)
	} else {
		config.ConsoleLoggingEnabled = true
		infoWriters = append(infoWriters, os.Stdout)
		errWriters = append(errWriters, os.Stderr)
	}

	logEntry := newZapLogger(config, zapcore.NewMultiWriteSyncer(infoWriters...), zapcore.NewMultiWriteSyncer(errWriters...), false)

	DeclareLogger(config, logEntry.Infov)
	DeclareLogger(config, logEntry.Errorv)
	return logEntry
}

func DeclareLogger(config Config, logv func(msg string, fields ...zapcore.Field)) {
	logv("logging configured",
		zap.Bool("fileLogging", config.FileLoggingEnabled),
		zap.Bool("consoleLogging", config.ConsoleLoggingEnabled),
		zap.Bool("caller", config.CallerEnabled),
		zap.Int("callerSkip", config.CallerSkip),
		zap.Bool("jsonLogOutput", config.EncodeLogsAsJson),
		zap.String("logDirectory", config.Directory),
		zap.Int("maxSizeMB", config.MaxSize),
		zap.Int("maxBackups", config.MaxBackups),
		zap.Int("maxAgeInDays", config.MaxAge))
}

func getNameByLogLevel(filename string, level zapcore.Level) string {
	var name string
	if filename != "" {
		filename = strings.Replace(filename, ".log", "", -1)
		name = filename + "_"
	}
	switch level {
	case ErrorLevel:
		name += "error.log"
	default:
		name += "info.log"
	}
	return name
}

func newRollingFile(dir, filename string, maxSize, maxAge, maxBackups int) zapcore.WriteSyncer {
	if err := os.MkdirAll(dir, 0744); err != nil {
		Errorv("failed create log directory", zap.Error(err), zap.String("path", dir))
		return nil
	}

	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   path.Join(dir, filename),
		MaxSize:    maxSize,    //megabytes
		MaxAge:     maxAge,     //days
		MaxBackups: maxBackups, //files
		LocalTime:  true,
	})
}

func newZapLogger(config Config, infoOutput zapcore.WriteSyncer, errOutput zapcore.WriteSyncer, isDefaultLogger bool) *LogEntry {
	encCfg := zapcore.EncoderConfig{
		TimeKey:          "@t",
		LevelKey:         "lvl",
		NameKey:          "logger",
		CallerKey:        "caller",
		MessageKey:       "msg",
		StacktraceKey:    "stacktrace",
		ConsoleSeparator: config.ConsoleSeparator,
		EncodeLevel:      config.LevelEncoder,
		EncodeDuration:   zapcore.NanosDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
	}
	if encCfg.EncodeLevel == nil {
		encCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	if !config.EncodeLogsAsJson {
		encCfg.ConsoleSeparator = config.ConsoleSeparator
		encCfg.EncodeTime = ConsoleLogTimeEncoder
		config.EncodeLogsAsJson = false
	} else {
		encCfg.EncodeTime = ShortTimeEncoder
	}

	encoder := zapcore.NewConsoleEncoder(encCfg)
	if config.EncodeLogsAsJson {
		encoder = zapcore.NewJSONEncoder(encCfg)
	}

	// gloval var `loglv` is reserved for changing log level of defaultLogger
	localLoglv := zap.NewAtomicLevelAt(config.Level)
	if isDefaultLogger {
		loglv = localLoglv
	}

	if config.CallerEnabled {
		return getLogEntry(
			zap.New(zapcore.NewCore(encoder, infoOutput, localLoglv),
				zap.AddCaller(), zap.AddCallerSkip(config.CallerSkip)),

			zap.New(zapcore.NewCore(encoder, errOutput, localLoglv),
				zap.AddCaller(), zap.AddCallerSkip(config.CallerSkip)))
	}
	return getLogEntry(zap.New(zapcore.NewCore(encoder, infoOutput, localLoglv)),
		zap.New(zapcore.NewCore(encoder, errOutput, localLoglv)))
}

func newRotateWriter(dir, fileName string) *lumberjack.Logger {
	logFilePath := path.Join(dir, fileName+".log")
	return &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    500, // megabytes
		MaxBackups: 10,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}
}

// AtLevel logs the message at a specific log level
func AtLevel(level zapcore.Level, msg string, fields ...zapcore.Field) {
	switch level {
	case zapcore.DebugLevel:
		Debugv(msg, fields...)
	case zapcore.PanicLevel:
		Panicv(msg, fields...)
	case zapcore.ErrorLevel:
		Errorv(msg, fields...)
	case zapcore.WarnLevel:
		Warnv(msg, fields...)
	case zapcore.InfoLevel:
		Infov(msg, fields...)
	case zapcore.FatalLevel:
		Fatalv(msg, fields...)
	default:
		Warnv("Logging at unkown level", zap.Any("level", level))
		Warnv(msg, fields...)
	}
}

func Debugv(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.infoLogger.Debug(msg, fields...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	DefaultZapLogger.Debugw(msg, keysAndValues...)
}

// Debugf Log a format message at the debug level
func Debugf(template string, args ...interface{}) {
	DefaultZapLogger.infoSugared.Debugf(template, args...)
}

// Debug Log a message at the debug level
func Debug(msg string) {
	DefaultZapLogger.infoLogger.Debug(msg)
}

// DebugWith Log a message with fields at the debug level
func DebugWith(msg string, fields Fields) {
	if len(fields) > 0 {
		DefaultZapLogger.infoLogger.Debug(msg, convertFields(fields)...)
	} else {
		DefaultZapLogger.infoLogger.Debug(msg)
	}
}

func Infov(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.infoLogger.Info(msg, fields...)
}

func Infof(template string, args ...interface{}) {
	DefaultZapLogger.infoSugared.Infof(template, args...)
}

func Info(msg string) {
	DefaultZapLogger.infoLogger.Info(msg)
}

func InfoWith(msg string, fields Fields) {
	if len(fields) > 0 {
		DefaultZapLogger.infoLogger.Info(msg, convertFields(fields)...)
	} else {
		DefaultZapLogger.infoLogger.Info(msg)
	}
}

func Infow(msg string, keysAndValues ...interface{}) {
	DefaultZapLogger.Infow(msg, keysAndValues...)
}

func Warnv(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.errorLogger.Warn(msg, fields...)
}

func Warnf(template string, args ...interface{}) {
	DefaultZapLogger.errorSugared.Warnf(template, args...)
}

func Warn(msg string) {
	DefaultZapLogger.errorLogger.Warn(msg)
}
func Warnw(msg string, keysAndValues ...interface{}) {
	DefaultZapLogger.Warnw(msg, keysAndValues...)
}

func WarnWith(msg string, fields Fields) {
	if len(fields) > 0 {
		DefaultZapLogger.errorLogger.Warn(msg, convertFields(fields)...)
	} else {
		DefaultZapLogger.errorLogger.Warn(msg)
	}
}

func Errorv(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.errorLogger.Error(msg, fields...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	DefaultZapLogger.Errorw(msg, keysAndValues...)
}

func Errorf(template string, args ...interface{}) {
	DefaultZapLogger.errorSugared.Errorf(template, args...)
}

func Error(msg string) {
	DefaultZapLogger.errorLogger.Error(msg)
}

func ErrorWith(msg string, fields Fields) {
	if len(fields) > 0 {
		DefaultZapLogger.errorLogger.Error(msg, convertFields(fields)...)
	} else {
		DefaultZapLogger.errorLogger.Error(msg)
	}
}

func Panicv(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.errorLogger.Panic(msg, fields...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	DefaultZapLogger.Panicw(msg, keysAndValues...)
}

func Panicf(template string, args ...interface{}) {
	DefaultZapLogger.errorSugared.Panicf(template, args...)
}

func Panic(msg string) {
	DefaultZapLogger.errorLogger.Panic(msg)
}

func PanicWith(msg string, fields Fields) {
	if len(fields) > 0 {
		DefaultZapLogger.errorLogger.Panic(msg, convertFields(fields)...)
	} else {
		DefaultZapLogger.errorLogger.Panic(msg)
	}
}

func Fatalv(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.errorLogger.Fatal(msg, fields...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	DefaultZapLogger.Fatalw(msg, keysAndValues...)
}

func Fatalf(template string, args ...interface{}) {
	DefaultZapLogger.errorSugared.Fatalf(template, args...)
}

func Fatal(msg string) {
	DefaultZapLogger.errorLogger.Fatal(msg)
}

func FatalWith(msg string, fields Fields) {
	if len(fields) > 0 {
		DefaultZapLogger.errorLogger.Fatal(msg, convertFields(fields)...)
	} else {
		DefaultZapLogger.errorLogger.Fatal(msg)
	}
}

func DPanicv(msg string, fields ...zapcore.Field) {
	DefaultZapLogger.errorLogger.DPanic(msg, fields...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	DefaultZapLogger.DPanicw(msg, keysAndValues...)
}

func DPanicf(template string, args ...interface{}) {
	DefaultZapLogger.errorSugared.DPanicf(template, args...)
}

func DPanic(msg string) {
	DefaultZapLogger.errorLogger.DPanic(msg)
}

func DPanicWith(msg string, fields Fields) {
	if len(fields) > 0 {
		DefaultZapLogger.errorLogger.DPanic(msg, convertFields(fields)...)
	} else {
		DefaultZapLogger.errorLogger.DPanic(msg)
	}
}

func WithFields(fields Fields) *LogEntry {
	return newLogEntry(DefaultZapLogger, fields)
}

func With(data string) *LogEntry {
	return WithField(DefaultFieldName, data)
}

func WithField(k, v string) *LogEntry {
	return newLogEntry(DefaultZapLogger, Fields{k: v})
}

func FromContext(ctx context.Context) *LogEntry {
	logger, ok := ctx.Value(loggerKey).(*LogEntry)
	if !ok {
		return DefaultZapLogger
	}
	return logger
}

func ContextWithLogger(ctx context.Context) context.Context {
	return DefaultZapLogger.ContextWithLogger(ctx)
}

func ContextWithCustomizedLogger(ctx context.Context, logEntry *LogEntry) context.Context {
	return logEntry.ContextWithLogger(ctx)
}
