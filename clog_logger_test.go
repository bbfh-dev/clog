package libclog_test

import (
	"net/http"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	libclog "codeberg.org/bbfh/lib-clog"
	"gotest.tools/assert"
)

func TestFormattingAuto(t *testing.T) {
	logger := libclog.NewLogger("TEST", os.Stdout)
	logger.SetLevel(libclog.LogDebug)
	print(t, logger)
	logger.SetOutputMode(libclog.OutputColored)
	print(t, logger)
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

func print(t *testing.T, logger *libclog.Logger) {
	assert.NilError(t, logger.Trace("Something happened with id=%d", 123))
	assert.NilError(t, logger.Debug("Something happened with id=%d", 123))
	assert.NilError(t, logger.Info("Something happened with id=%d", 123))
	assert.NilError(t, logger.Notice("Something happened with id=%d", 123))
	assert.NilError(t, logger.Warn("Something happened with id=%d", 123))
	assert.NilError(t, logger.Error("Something happened with id=%d", 123))
	assert.NilError(t, logger.Fatal("Something happened with id=%d", 123))
	uri, err := url.Parse("/index")
	assert.NilError(t, err)
	assert.NilError(t, logger.HttpRequest(&http.Request{
		Method: "GET",
		URL:    uri,
	}, time.Now().Add(-25*time.Second), 204))
}
