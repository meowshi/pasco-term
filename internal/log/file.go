//go:build !debug

package log

import (
	"fmt"
	"os"

	"github.com/meowshi/pasco/internal/utils"
	"github.com/sirupsen/logrus"
)

func InitLog() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logFile, err := os.OpenFile(utils.AppFolderPath()+"/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("ошибка при работе с лог-файлом: %s", err))
	}

	logrus.SetOutput(logFile)
}
