package log

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	log *zap.Logger
}

type Config struct {
	File       FileConfig  `json:"enable_file"`
	Level      LoggerLevel `json:"level"`
	StackTrace bool        `json:"stack_trace"`
}

type FileConfig struct {
	Enable       bool   `json:"enable"`
	FileLocation string `json:"file_location"`
	MaxSize      int    `json:"max_size"`
	MaxBackups   int    `json:"max_backups"`
	MaxAge       int    `json:"max_age"`
	Compress     bool   `json:"compress"`
}

func (cfg Config) Build() *Logger {
	if !cfg.File.Enable {
		config := zap.NewProductionConfig()
		if !cfg.StackTrace {
			config.DisableStacktrace = true
		}
		logger, err := config.Build()
		if err != nil {
			log.Fatalf("failed to init logger: %s", err)
		}
		return &Logger{
			log: logger,
		}
	}

	opt := []zap.Option{
		zap.AddCaller(),
	}
	if cfg.StackTrace {
		opt = append(opt, zap.AddStacktrace(setLogLevel(cfg.Level)))
	}

	logger := &lumberjack.Logger{
		Filename:   cfg.File.FileLocation,
		MaxSize:    cfg.File.MaxSize,
		MaxBackups: cfg.File.MaxBackups,
		MaxAge:     cfg.File.MaxAge,
		Compress:   cfg.File.Compress,
	}
	writeSyncer := zapcore.AddSync(logger)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		setLogLevel(cfg.Level),
	)
	return &Logger{
		log: zap.New(core, opt...),
	}
}

func New(cfg *Config) *Logger {
	if cfg == nil {
		log.Fatalf("failed to init logger: %s", "config is nil")
	}
	return cfg.Build()
}

// Sync flushes any buffered log entries.
func (l *Logger) Sync() error {
	if l.log == nil {
		return nil
	}
	return l.log.Sync()
}

// Info logs a message at level Info.
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.log.Info(msg, fields...)
}

// Debug logs a message at level Debug.
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.log.Debug(msg, fields...)
}

// Warn logs a message at level Warn.
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.log.Warn(msg, fields...)
}

// Error logs a message at level Error.
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.log.Error(msg, fields...)
}

func setLogLevel(level LoggerLevel) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "info":
		fallthrough
	default:
		return zapcore.InfoLevel
	}
}
