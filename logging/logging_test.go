package logging

import (
	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigure_File(t *testing.T) {
	config := Config{
		File: "/tmp/test.log",
	}
	closer := Configure(config)
	defer closer()
	assert.Equal(t, logrus.DebugLevel, logrus.StandardLogger().Level)
}
