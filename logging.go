package kaiju
import (
	"github.com/cihub/seelog"
	"io"
	"errors"
)


var logger seelog.LoggerInterface

func init() {
	DisableLogging()
}

func DisableLogging() {
	UseLogger(seelog.Disabled)
}

func UseLogger(newLogger seelog.LoggerInterface) {
	logger = newLogger
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