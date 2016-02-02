package log

import (
	"os"

	"github.com/op/go-logging"
)

var (
	Logger *logging.Logger
)

func init() {
	Logger = logging.MustGetLogger("kahinah")
	logFormat := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`)
	logBackend := logging.NewLogBackend(os.Stdout, "", 0)
	logFormattedBackend := logging.NewBackendFormatter(logBackend, logFormat)
	logging.SetBackend(logFormattedBackend)
}
