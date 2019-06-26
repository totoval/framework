package logs

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/totoval/framework/config"
	"github.com/totoval/framework/sentry"
)

var log *logrus.Logger
var logLevel Level

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.Out = os.Stdout
}

func Initialize() {
	levelStr := config.GetString("app.log_level")
	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		panic(err)
	}

	logLevel = level
	log.SetLevel(logLevel)
}

type Field = map[string]interface{}

func Println(level Level, msg string, fields Field) {
	if level >= logLevel {
		sentry.CaptureMsg(msg, fields)
	}

	if fields == nil {
		log.Log(level, msg)
	} else {
		log.WithFields(logrus.Fields(fields)).Log(level, msg)
	}
}
