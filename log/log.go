package log

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
	k.Info(v)
}

func init() {
	Logger = &KLogger{logging.MustGetLogger("kahinah")}
	logFormat := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`)
	logBackend := logging.NewLogBackend(os.Stdout, "", 0)
	logFormattedBackend := logging.NewBackendFormatter(logBackend, logFormat)
	logging.SetBackend(logFormattedBackend)
}
