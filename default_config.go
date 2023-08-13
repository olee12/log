package log

import (
	"os"
	"path/filepath"

	"go.uber.org/zap/zapcore"
)

var DefaultRotateLoggerConfig = &Config{
	Level:              DebugLevel,
	EncodeLogsAsJson:   true,
	CallerEnabled:      true,
	CallerSkip:         1,
	ConsoleSeparator:   "",
	LevelEncoder:       zapcore.LowercaseLevelEncoder,
	Directory:          "logs",
	Filename:           filepath.Base(os.Args[0]),
	FileLoggingEnabled: true,
	MaxSize:            128,
	MaxBackups:         10,
	MaxAge:             90,
}
