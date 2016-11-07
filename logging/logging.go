package logging

import (
	"github.com/Sirupsen/logrus"
	"os"
)

type Config struct {
	Level  string               `json:"level" yaml:"level"`
	File   string               `json:"file" yaml:"file"`
	Format logrus.TextFormatter `json:"format" yaml:"format"`
}

// Configure configures logrus's standard logger with the provided config object. A closer function is returned
// that should be used to close the logger and any open file descriptors before the application terminates.
func Configure(config Config) func() error {
	logger, closer := New(config)
	stdLogger := logrus.StandardLogger()
	*stdLogger = *logger
	return closer
}

// New creates a new *logrus.Logger instnace using the provided config object. A closer function is returned
// that should be used to close the logger and any open file descriptors before the application terminates.
func New(config Config) (*logrus.Logger, func() error) {
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.DebugLevel
	}
	out := os.Stdout
	closer := func() error {
		return nil
	}
	if config.File != "" {
		f, err := os.OpenFile(config.File, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		if err == nil {
			out = f
			closer = func() error {
				return f.Close()
			}
		}
	}
	return &logrus.Logger{
		Out:       out,
		Formatter: &config.Format,
		Hooks:     make(logrus.LevelHooks),
		Level:     level,
	}, closer
}
