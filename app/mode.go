package app

import (
	"github.com/gin-gonic/gin"
)

var mode Mode

type Mode byte

const (
	ModeProduction Mode = iota
	ModeDevelop
	ModeTest
)

func init() {
	SetMode(ModeDevelop)
}

func SetMode(m Mode) {

	gin.SetMode(gin.ReleaseMode)

	switch m {
	case ModeTest:
		//.SetMode(gin.TestMode)
	case ModeDevelop:
		// gin.SetMode(gin.DebugMode)
	default:
		fallthrough
	case ModeProduction:
		//gin.SetMode(gin.ReleaseMode)
	}

	mode = m
}

func GetMode() Mode {
	return mode
}
