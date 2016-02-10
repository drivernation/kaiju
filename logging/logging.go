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

func SetLogWriter(writer io.Writer, minLevel seelog.LogLevel) error {
	if writer == nil {
		return errors.New("Nil writer.")
	}

	newLogger, err := seelog.LoggerFromWriterWithMinLevel(writer, minLevel)
	if err != nil {
		return err
	}

	UseLogger(newLogger)
	return nil
}
