package log

import (
	"os"

	"github.com/hhung06/digimap-backend/config"
	"github.com/sirupsen/logrus"
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

// wrappedEntry wraps a logrus.Entry to satisfy the Logger interface.
type wrappedEntry struct {
	entry *logrus.Entry
}

func (w *wrappedEntry) Debug(args ...interface{})                 { w.entry.Debug(args...) }
func (w *wrappedEntry) Debugf(f string, args ...interface{})      { w.entry.Debugf(f, args...) }
func (w *wrappedEntry) Debugln(args ...interface{})               { w.entry.Debugln(args...) }
func (w *wrappedEntry) Error(args ...interface{})                 { w.entry.Error(args...) }
func (w *wrappedEntry) Errorf(f string, args ...interface{})      { w.entry.Errorf(f, args...) }
func (w *wrappedEntry) Errorln(args ...interface{})               { w.entry.Errorln(args...) }
func (w *wrappedEntry) Fatal(args ...interface{})                 { w.entry.Fatal(args...) }
func (w *wrappedEntry) Fatalf(f string, args ...interface{})      { w.entry.Fatalf(f, args...) }
func (w *wrappedEntry) Fatalln(args ...interface{})               { w.entry.Fatalln(args...) }
func (w *wrappedEntry) Info(args ...interface{})                  { w.entry.Info(args...) }
func (w *wrappedEntry) Infof(f string, args ...interface{})       { w.entry.Infof(f, args...) }
func (w *wrappedEntry) Infoln(args ...interface{})                { w.entry.Infoln(args...) }
func (w *wrappedEntry) Panic(args ...interface{})                 { w.entry.Panic(args...) }
func (w *wrappedEntry) Panicf(f string, args ...interface{})      { w.entry.Panicf(f, args...) }
func (w *wrappedEntry) Panicln(args ...interface{})               { w.entry.Panicln(args...) }
func (w *wrappedEntry) Print(args ...interface{})                 { w.entry.Print(args...) }
func (w *wrappedEntry) Printf(f string, args ...interface{})      { w.entry.Printf(f, args...) }
func (w *wrappedEntry) Println(args ...interface{})               { w.entry.Println(args...) }
func (w *wrappedEntry) Warn(args ...interface{})                  { w.entry.Warn(args...) }
func (w *wrappedEntry) Warnf(f string, args ...interface{})       { w.entry.Warnf(f, args...) }
func (w *wrappedEntry) Warnln(args ...interface{})                { w.entry.Warnln(args...) }
func (w *wrappedEntry) WithFields(fields Fields) Logger {
	return &wrappedEntry{entry: w.entry.WithFields(logrus.Fields(fields))}
}

// logrusLogger wraps *logrus.Logger and implements Logger.
type logrusLogger struct {
	l *logrus.Logger
}

func (r *logrusLogger) Debug(args ...interface{})                 { r.l.Debug(args...) }
func (r *logrusLogger) Debugf(f string, args ...interface{})      { r.l.Debugf(f, args...) }
func (r *logrusLogger) Debugln(args ...interface{})               { r.l.Debugln(args...) }
func (r *logrusLogger) Error(args ...interface{})                 { r.l.Error(args...) }
func (r *logrusLogger) Errorf(f string, args ...interface{})      { r.l.Errorf(f, args...) }
func (r *logrusLogger) Errorln(args ...interface{})               { r.l.Errorln(args...) }
func (r *logrusLogger) Fatal(args ...interface{})                 { r.l.Fatal(args...) }
func (r *logrusLogger) Fatalf(f string, args ...interface{})      { r.l.Fatalf(f, args...) }
func (r *logrusLogger) Fatalln(args ...interface{})               { r.l.Fatalln(args...) }
func (r *logrusLogger) Info(args ...interface{})                  { r.l.Info(args...) }
func (r *logrusLogger) Infof(f string, args ...interface{})       { r.l.Infof(f, args...) }
func (r *logrusLogger) Infoln(args ...interface{})                { r.l.Infoln(args...) }
func (r *logrusLogger) Panic(args ...interface{})                 { r.l.Panic(args...) }
func (r *logrusLogger) Panicf(f string, args ...interface{})      { r.l.Panicf(f, args...) }
func (r *logrusLogger) Panicln(args ...interface{})               { r.l.Panicln(args...) }
func (r *logrusLogger) Print(args ...interface{})                 { r.l.Print(args...) }
func (r *logrusLogger) Printf(f string, args ...interface{})      { r.l.Printf(f, args...) }
func (r *logrusLogger) Println(args ...interface{})               { r.l.Println(args...) }
func (r *logrusLogger) Warn(args ...interface{})                  { r.l.Warn(args...) }
func (r *logrusLogger) Warnf(f string, args ...interface{})       { r.l.Warnf(f, args...) }
func (r *logrusLogger) Warnln(args ...interface{})                { r.l.Warnln(args...) }
func (r *logrusLogger) WithFields(fields Fields) Logger {
	return &wrappedEntry{entry: r.l.WithFields(logrus.Fields(fields))}
}

// NewLogger creates a Logger from the typed Config.
func NewLogger(cfg *config.Config) Logger {
	l := logrus.New()
	l.Out = os.Stderr

	if cfg.App.LogFormat == "json" {
		l.Formatter = &logrus.JSONFormatter{}
	} else {
		l.Formatter = &logrus.TextFormatter{FullTimestamp: true}
	}

	switch cfg.App.LogLevel {
	case "warn", "warning":
		l.Level = logrus.WarnLevel
	case "info":
		l.Level = logrus.InfoLevel
	case "error":
		l.Level = logrus.ErrorLevel
	default:
		l.Level = logrus.DebugLevel
	}

	return &logrusLogger{l: l}
}
