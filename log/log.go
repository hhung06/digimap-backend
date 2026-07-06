package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/hhung06/digimap-backend/config"
)

// Logger defines methods for structured application logging.
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Panicln(args ...interface{})
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})
	WithFields(fields Fields) Logger
}

// Fields is a map of structured log fields.
type Fields map[string]interface{}

type zapLogger struct {
	s *zap.SugaredLogger
}

func (z *zapLogger) Debug(args ...interface{})            { z.s.Debug(args...) }
func (z *zapLogger) Debugf(f string, args ...interface{}) { z.s.Debugf(f, args...) }
func (z *zapLogger) Debugln(args ...interface{})          { z.s.Debug(fmt.Sprint(args...)) }
func (z *zapLogger) Error(args ...interface{})            { z.s.Error(args...) }
func (z *zapLogger) Errorf(f string, args ...interface{}) { z.s.Errorf(f, args...) }
func (z *zapLogger) Errorln(args ...interface{})          { z.s.Error(fmt.Sprint(args...)) }
func (z *zapLogger) Fatal(args ...interface{})            { z.s.Fatal(args...) }
func (z *zapLogger) Fatalf(f string, args ...interface{}) { z.s.Fatalf(f, args...) }
func (z *zapLogger) Fatalln(args ...interface{})          { z.s.Fatal(fmt.Sprint(args...)) }
func (z *zapLogger) Info(args ...interface{})             { z.s.Info(args...) }
func (z *zapLogger) Infof(f string, args ...interface{})  { z.s.Infof(f, args...) }
func (z *zapLogger) Infoln(args ...interface{})           { z.s.Info(fmt.Sprint(args...)) }
func (z *zapLogger) Panic(args ...interface{})            { z.s.Panic(args...) }
func (z *zapLogger) Panicf(f string, args ...interface{}) { z.s.Panicf(f, args...) }
func (z *zapLogger) Panicln(args ...interface{})          { z.s.Panic(fmt.Sprint(args...)) }
func (z *zapLogger) Print(args ...interface{})            { z.s.Info(args...) }
func (z *zapLogger) Printf(f string, args ...interface{}) { z.s.Infof(f, args...) }
func (z *zapLogger) Println(args ...interface{})          { z.s.Info(fmt.Sprint(args...)) }
func (z *zapLogger) Warn(args ...interface{})             { z.s.Warn(args...) }
func (z *zapLogger) Warnf(f string, args ...interface{})  { z.s.Warnf(f, args...) }
func (z *zapLogger) Warnln(args ...interface{})           { z.s.Warn(fmt.Sprint(args...)) }

func (z *zapLogger) WithFields(fields Fields) Logger {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &zapLogger{s: z.s.With(args...)}
}

// Sync flushes buffered log entries. Not on the Logger interface.
func (z *zapLogger) Sync() error { return z.s.Sync() }

// Sync flushes the logger if it supports it. Call with defer after NewLogger.
func Sync(l Logger) {
	if s, ok := l.(interface{ Sync() error }); ok {
		_ = s.Sync() // discard error — ENOTTY/EINVAL on TTY/pipe stderr is harmless
	}
}

// NewLogger creates a Logger from the typed Config.
func NewLogger(cfg *config.Config) Logger {
	var level zapcore.Level
	switch cfg.App.LogLevel {
	case "warn", "warning":
		level = zapcore.WarnLevel
	case "info":
		level = zapcore.InfoLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.DebugLevel
	}

	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = "time"
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	var encoder zapcore.Encoder
	if cfg.App.LogFormat == "json" {
		encoder = zapcore.NewJSONEncoder(encCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(encCfg)
	}

	writers := []zapcore.WriteSyncer{zapcore.AddSync(os.Stderr)}
	if cfg.App.LogDir != "" {
		rotator, err := rotatelogs.New(
			filepath.Join(cfg.App.LogDir, "app-%Y-%m-%d.log"),
			rotatelogs.WithMaxAge(30*24*time.Hour),
			rotatelogs.WithRotationTime(24*time.Hour),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to init log rotation: %v", err))
		}
		writers = append(writers, zapcore.AddSync(rotator))
	}

	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writers...), level)
	return &zapLogger{s: zap.New(core).Sugar()}
}
