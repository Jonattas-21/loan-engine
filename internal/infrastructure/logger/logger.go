package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
)

func LogSetup() *logrus.Logger {
	// Create a new instance of the logger
	logger := logrus.New()

	// Set the output to standard error
	logger.SetOutput(os.Stderr)

	// Set the log level (e.g., Info, Warn, Error, Debug)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetLevel(logrus.ErrorLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", f.File, f.Line)
		},
	})

	// Enable reporting the caller information
	logger.SetReportCaller(true)

	return logger
}
