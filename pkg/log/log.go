package log

/*
We define some logging functions here that act as a proxy to our
real implementation of the logger. We do it this way so that the
rest of the cli can use the logger without worrying the real implementation.
In this case, we use uber's super fast zap logger, but we replace with our
own logger if necessary.
*/

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger, _ = zap.NewProduction()
var log *zap.SugaredLogger
var atom zap.AtomicLevel

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	DPanic(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	DPanicf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	DPanicw(msg string, keysAndValues ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
}

func init() {
	// change the encoding to console from json
	cfg := zap.NewProductionConfig()
	cfg.Encoding = "console"
	atom = zap.NewAtomicLevelAt(zap.WarnLevel)
	cfg.Level = atom
	if paralusLogger, err := cfg.Build(); err == nil {
		// if there was no error, use the custom config, otherwise
		// use the default config
		logger = paralusLogger
	}

	// flush buffer, if any
	defer logger.Sync()
	log = logger.Sugar()
}

func GetLogger() *zap.SugaredLogger {
	return log
}

func SetLevel(level zapcore.Level) {
	atom.SetLevel(level)
}
