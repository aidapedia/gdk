package log

import (
	"context"
	"fmt"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	gdkCtx "github.com/aidapedia/gdk/context"
)

var Log *Logger

type Logger struct {
	*zap.Logger
	defaultTags map[string]interface{}
}

type Config struct {
	File        FileConfig  `json:"enable_file"`
	Level       LoggerLevel `json:"level"`
	Caller      bool        `json:"caller"`
	DefaultTags map[string]interface{}
}

type FileConfig struct {
	Enable       bool   `json:"enable"`
	FileLocation string `json:"file_location"`
	MaxSize      int    `json:"max_size"`
	MaxBackups   int    `json:"max_backups"`
	MaxAge       int    `json:"max_age"`
	Compress     bool   `json:"compress"`
}

func (cfg Config) build() *Logger {
	if !cfg.File.Enable {
		config := zap.NewProductionConfig()
		if !cfg.Caller {
			config.DisableCaller = true
		}
		logger, err := config.Build()
		if err != nil {
			log.Fatalf("failed to init logger: %s", err)
		}
		return &Logger{
			Logger:      logger,
			defaultTags: cfg.DefaultTags,
		}
	}

	opt := []zap.Option{}
	if cfg.Caller {
		opt = append(opt, zap.AddCaller())
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
		Logger:      zap.New(core, opt...),
		defaultTags: cfg.DefaultTags,
	}
}

func New(cfg *Config) {
	if cfg == nil {
		panic("failed to init logger: config is nil")
	}
	Log = cfg.build()
}

// Sync flushes any buffered log entries.
func Sync() error {
	if Log == nil {
		return nil
	}
	return Log.Sync()
}

// InfoCtx logs a message at level Info with log id.
func InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, fieldCheck(ctx)...)
	Log.Info(msg, fields...)
}

// DebugCtx logs a message at level Debug with log id.
func DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, fieldCheck(ctx)...)
	Log.Debug(msg, fields...)
}

// WarnCtx logs a message at level Warn with log id.
func WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, fieldCheck(ctx)...)
	Log.Warn(msg, fields...)
}

// ErrorCtx logs a message at level Error with log id.
func ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	fields = append(fields, fieldCheck(ctx)...)
	Log.Error(msg, fields...)
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

func fieldCheck(ctx context.Context) []zap.Field {
	var fields = make([]zap.Field, 0)
	logID := ctx.Value(gdkCtx.ContextKeyLogID)
	if logID != nil {
		fields = append(fields, zap.String("log_id", fmt.Sprintf("%s", logID)))
	}
	for k, v := range Log.defaultTags {
		fields = append(fields, zap.Any(k, v))
	}
	return fields
}
