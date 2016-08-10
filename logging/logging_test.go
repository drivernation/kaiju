package logging

import (
	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigure(t *testing.T) {
	config := Config{
		Level: "error",
	}
	closer := Configure(config)
	defer closer()
	assert.Equal(t, logrus.ErrorLevel, logrus.StandardLogger().Level)
}
