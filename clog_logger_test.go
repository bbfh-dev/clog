package libclog_test

import (
	"net/http"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"codeberg.org/bbfh/libclog"
)

func TestFormattingAuto(t *testing.T) {
	logger := libclog.NewLogger("TEST", os.Stdout)
	logger.SetLevel(libclog.LogDebug)
	print(logger)
	logger.SetOutputMode(libclog.OutputColored)
	print(logger)
}

func TestConcurrency(t *testing.T) {
	var wg sync.WaitGroup
	logger := libclog.NewLogger("TEST", os.Stdout).SetLevel(libclog.LogInfo)
	for i := range 100 {
		wg.Go(func() {
			logger.Info("Call from routine %d", i)
		})
	}
	wg.Wait()
}

func print(logger *libclog.Logger) {
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
