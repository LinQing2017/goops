package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log *logrus.Logger

func InitLogger() {
	logrus.SetOutput(os.Stdout)
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.WithFields(logrus.Fields{
		"sdp-devops": "Kubernetes",
	})

}
