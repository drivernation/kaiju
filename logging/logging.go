package logging

import (
	"errors"
	"github.com/cihub/seelog"
	"io"
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

