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
	logger.SetFormatter(&CustomFormatter{
        TextFormatter: logrus.TextFormatter{
            FullTimestamp: true,
        },
    })
	// Enable reporting the caller information
	logger.SetReportCaller(true)

	return logger
}

type CustomFormatter struct {
    logrus.TextFormatter
}
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
    _, file, line, ok := runtime.Caller(8)
    if !ok {
        file = "unknown"
        line = 0
    }

    return []byte(fmt.Sprintf("[%s] %s %s: %s:%d\n",
        entry.Level.String(),
        timestamp,
        entry.Message,
        file,
        line,
    )), nil
}