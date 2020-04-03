package waterlogged

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/orandin/lumberjackrus"
	"github.com/sirupsen/logrus"
)

// Waterlogged represents the logger instance with service name details.
type Waterlogged struct {
	*logrus.Logger
	serviceName string
}

// New creates a new logger instance with the specified service name
// and creates a corresponding Lumberjack hook for file logging.
func New(serviceName string) *Waterlogged {
	wl := &Waterlogged{
		Logger:      logrus.New(),
		serviceName: serviceName,
	}
	wl.SetFormatter(&logrus.TextFormatter{})
	return wl
}

func (wl *Waterlogged) addLumberjackHook() {
	pwd, err := os.Getwd()
	logDirPath := filepath.Join(pwd, "logs")
	err = os.MkdirAll(logDirPath, os.ModePerm)

	hook, err := lumberjackrus.NewHook(
		&lumberjackrus.LogFile{
			Filename:  filepath.Join(logDirPath, wl.serviceName+".log"),
			MaxSize:   100,
			MaxAge:    60,
			Compress:  false,
			LocalTime: false,
		},
		logrus.InfoLevel,
		&logrus.TextFormatter{},
		&lumberjackrus.LogFileOpts{},
	)

	if err != nil {
		fmt.Printf("Error adding Lumberjack hook: %v\n", err)
	}

	wl.AddHook(hook)
}
