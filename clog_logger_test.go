package clog_test

import (
	"os"
	"testing"

	"codeberg.org/bbfh/clog"
)

func TestFormattingAuto(t *testing.T) {
	logger := clog.NewLogger("TEST", os.Stdout)
	logger.SetLevel(clog.LOG_DEBUG)
	print(logger)
	logger.SetOutputMode(clog.OUTPUT_COLORED)
	print(logger)
}

func print(logger *clog.Logger) {
	logger.Trace("Something happened with id=%d", 123)
	logger.Debug("Something happened with id=%d", 123)
	logger.Info("Something happened with id=%d", 123)
	logger.Notice("Something happened with id=%d", 123)
	logger.Warn("Something happened with id=%d", 123)
	logger.Error("Something happened with id=%d", 123)
	logger.Fatal("Something happened with id=%d", 123)
}
