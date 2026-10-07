package clog_test

import (
	"net/http"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"codeberg.org/bbfh/clog"
)

func TestFormattingAuto(t *testing.T) {
	logger := clog.NewLogger("TEST", os.Stdout)
	logger.SetLevel(clog.LogDebug)
	print(logger)
	logger.SetOutputMode(clog.OutputColored)
	print(logger)
}

func TestConcurrency(t *testing.T) {
	var wg sync.WaitGroup
	logger := clog.NewLogger("TEST", os.Stdout).SetLevel(clog.LogInfo)
	for i := range 100 {
		wg.Go(func() {
			logger.Info("Call from routine %d", i)
		})
	}
	wg.Wait()
}

func print(logger *clog.Logger) {
	logger.Trace("Something happened with id=%d", 123)
	logger.Debug("Something happened with id=%d", 123)
	logger.Info("Something happened with id=%d", 123)
	logger.Notice("Something happened with id=%d", 123)
	logger.Warn("Something happened with id=%d", 123)
	logger.Error("Something happened with id=%d", 123)
	logger.Fatal("Something happened with id=%d", 123)
	uri, _ := url.Parse("/index")
	logger.HttpRequest(&http.Request{
		Method: "GET",
		URL:    uri,
	}, time.Now().Add(-25*time.Second), 204)
}
