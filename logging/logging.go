package logging

import (
	"github.com/cihub/seelog"
)

var Logger seelog.LoggerInterface

func init() {
	DisableLogging()
}

func DisableLogging() {
	UseLogger(seelog.Disabled)
}

func UseLogger(newLogger seelog.LoggerInterface) {
	Logger = newLogger
}
