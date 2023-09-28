//go:build debug

package log

import "github.com/sirupsen/logrus"

func InitLog() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
