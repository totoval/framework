package logs

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/totoval/framework/config"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.Out = os.Stdout
}

func Initialize() {
	levelStr := config.GetString("app.log_level")
	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		panic(err)
	}
	log.SetLevel(level)
}

type Field = map[string]interface{}

func Println(level Level, msg string, fields Field) {
	if fields == nil {
		log.Log(level, msg)
	} else {
		log.WithFields(logrus.Fields(fields)).Log(level, msg)
	}
}
