// Package logger
package logger

import (
	"os"
	"strings"

	"github.com/ishi-o/transfmt/internal/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log     *zap.Logger
	Sugared *zap.SugaredLogger
)

func Init(cfg *config.ZapConfig) {
	encoder, level := zapEncoder(cfg), zapLevelEnabler(cfg)
	ws, errws := zapWriterSyncer(cfg)
	Log = zap.New(zapcore.NewTee(
		zapcore.NewCore(encoder, ws, level),
		zapcore.NewCore(encoder, errws, zap.ErrorLevel),
	))
	Sugared = Log.Sugar()
}

func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

func Infof(temp string, args ...any) {
	Sugared.Infof(temp, args)
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}

func zapEncoderConfig(cfg *config.ZapConfig) zapcore.EncoderConfig {
	var encoderCfg zapcore.EncoderConfig
	if cfg.Development {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}
	return encoderCfg
}

func zapEncoder(cfg *config.ZapConfig) zapcore.Encoder {
	encoderCfg := zapEncoderConfig(cfg)
	switch cfg.Encoding {
	case "json":
		return zapcore.NewJSONEncoder(encoderCfg)
	case "console":
		return zapcore.NewConsoleEncoder(encoderCfg)
	default:
		return zapcore.NewConsoleEncoder(encoderCfg)
	}
}

var levelMap = map[string]zapcore.Level{
	config.DebugLevel: zap.DebugLevel,
	config.InfoLevel:  zap.InfoLevel,
	config.WarnLevel:  zap.WarnLevel,
	config.ErrorLevel: zap.ErrorLevel,
	config.PanicLevel: zap.PanicLevel,
	config.FatalLevel: zap.FatalLevel,
}

func zapLevelEnabler(cfg *config.ZapConfig) zapcore.LevelEnabler {
	level, ok := levelMap[cfg.Level]
	if ok {
		return level
	} else {
		return levelMap[config.DefaultLevel]
	}
}

func newWriterSyncer(paths []string) zapcore.WriteSyncer {
	syncers := make([]zapcore.WriteSyncer, 0, len(paths))
	for _, path := range paths {
		var ws zapcore.WriteSyncer
		switch {
		case strings.ToLower(path) == "stdout":
			ws = zapcore.AddSync(os.Stdout)
		case strings.ToLower(path) == "stderr":
			ws = zapcore.AddSync(os.Stderr)
		default:
			lumberjackLogger := &lumberjack.Logger{
				Filename:   path,
				MaxSize:    100,
				MaxBackups: 30,
				MaxAge:     90,
				Compress:   true,
			}
			ws = zapcore.AddSync(lumberjackLogger)
		}
		syncers = append(syncers, ws)
	}
	return zapcore.NewMultiWriteSyncer(syncers...)
}

func zapWriterSyncer(cfg *config.ZapConfig) (zapcore.WriteSyncer, zapcore.WriteSyncer) {
	return newWriterSyncer(cfg.OutputPaths), newWriterSyncer(cfg.ErrorOutputPaths)
}
