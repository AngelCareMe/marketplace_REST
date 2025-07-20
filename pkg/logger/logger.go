package logger

import "github.com/sirupsen/logrus"

func SetupLogger(logLevel, logFormat string) *logrus.Logger {
	logger := logrus.New()

	switch logLevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
		logrus.Warnf("Unknown log level: %s, using info", logLevel)
	}

	switch logFormat {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{})
	default:
		logger.SetFormatter(&logrus.JSONFormatter{})
		logrus.Warnf("Unknown log format: %s, using json", logFormat)
	}

	return logger
}
