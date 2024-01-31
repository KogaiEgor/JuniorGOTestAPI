package initializers

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func InitLogrus() {
	Log.SetLevel(logrus.DebugLevel)

	Log.SetFormatter(&logrus.JSONFormatter{})

	Log.SetOutput(os.Stdout)
}
