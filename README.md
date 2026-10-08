# Clog

a Go library for clogging the output: syslog, stdout, stderr; with optional ansi-based formatting.

[Go Documentation](https://pkg.go.dev/codeberg.org/bbfh/lib-clog)

## Example usage

```go
// Create a new logger with name "TEST"
logger := clog.NewLogger("TEST", os.Stderr)
// (Default) automatically detect whether the output supports ansi escapes
logger.SetOutputMode(clog.OUTPUT_AUTO).
    SetLevel(clog.LOG_INFO).
    // (Returns an error that may be ignored) creates a syslog logger
    UseSyslog()

logger.Notice("Something happened with id=%d", 123)

start_time := time.Now()
// ...
logger.HttpRequest(request, start_time, 204)
```
