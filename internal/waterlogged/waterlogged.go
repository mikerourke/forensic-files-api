package waterlogged

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/orandin/lumberjackrus"
	"github.com/sirupsen/logrus"
)

// ServiceLogger creates a new logger instance with the specified service name
// and creates a corresponding Lumberjack hook for file logging.
func ServiceLogger(serviceName string) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})

	addLumberjackHook(logger, serviceName)
	return logger
}

func addLumberjackHook(logger *logrus.Logger, serviceName string) {
	pwd, err := os.Getwd()
	logDirPath := filepath.Join(pwd, "logs")
	err = os.MkdirAll(logDirPath, os.ModePerm)

	hook, err := lumberjackrus.NewHook(
		&lumberjackrus.LogFile{
			Filename:  filepath.Join(logDirPath, serviceName+".json"),
			MaxSize:   100,
			MaxAge:    60,
			Compress:  false,
			LocalTime: false,
		},
		logrus.InfoLevel,
		&logrus.JSONFormatter{},
		&lumberjackrus.LogFileOpts{},
	)

	if err != nil {
		fmt.Printf("Error adding Lumberjack hook: %v\n", err)
	}

	logger.AddHook(hook)
}
