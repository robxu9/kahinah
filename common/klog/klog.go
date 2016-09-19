package klog

import (
	"os"

	"github.com/op/go-logging"
)

var (
	// Logger is the main Kahinah logger
	Logger *KLogger
)

// KLogger is a wrapper around the logging Logger
type KLogger struct {
	*logging.Logger
}

// Println is a compatibility method for Info(v)
func (k *KLogger) Println(v ...interface{}) {
	k.Info(v...)
}

func init() {
	Logger = &KLogger{logging.MustGetLogger("kahinah")}
	logFormat := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`)
	logBackend := logging.NewLogBackend(os.Stdout, "", 0)
	logFormattedBackend := logging.NewBackendFormatter(logBackend, logFormat)
	logging.SetBackend(logFormattedBackend)
}

// Critical writes a log message to the Critical layer.
func Critical(args ...interface{}) {
	Logger.Critical(args...)
}

// Criticalf writes a log message to the Critical layer.
func Criticalf(format string, args ...interface{}) {
	Logger.Criticalf(format, args...)
}

// Debug writes a log message to the Debug layer.
func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

// Debugf writes a log message to the Debug layer.
func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

// Error writes a log message to the Error layer.
func Error(args ...interface{}) {
	Logger.Error(args...)
}

// Errorf writes a log message to the Error layer.
func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

// Fatal writes a log message to the Fatal layer.
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

// Fatalf writes a log message to the Fatal layer.
func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}

// Info writes a log message to the Info layer.
func Info(args ...interface{}) {
	Logger.Info(args...)
}

// Infof writes a log message to the Info layer.
func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

// Notice writes a log message to the Notice layer.
func Notice(args ...interface{}) {
	Logger.Notice(args...)
}

// Noticef writes a log message to the Notice layer.
func Noticef(format string, args ...interface{}) {
	Logger.Noticef(format, args...)
}

// Panic writes a log message to the Panic layer.
func Panic(args ...interface{}) {
	Logger.Panic(args...)
}

// Panicf writes a log message to the Panic layer.
func Panicf(format string, args ...interface{}) {
	Logger.Panicf(format, args...)
}

// Warning writes a log message to the Warning layer.
func Warning(args ...interface{}) {
	Logger.Warning(args...)
}

// Warningf writes a log message to the Warning layer.
func Warningf(format string, args ...interface{}) {
	Logger.Warningf(format, args...)
}
