package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func New(logLevel string) *logrus.Logger {
	l := logrus.New()
	l.Out = os.Stdout
	switch logLevel {
	case "debug":
		l.SetLevel(logrus.DebugLevel)
	case "info":
		l.SetLevel(logrus.InfoLevel)
	case "warn":
		l.SetLevel(logrus.WarnLevel)
	default:
		l.SetLevel(logrus.InfoLevel)
	}
	l.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	return l
}
