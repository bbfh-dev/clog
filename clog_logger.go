package clog

import (
	"fmt"
	"log/syslog"
	"net/http"
	"os"
	"strconv"
	"time"
)

type LevelEnum uint8

const (
	LogSilent  LevelEnum = iota * 4 // Highest priority
	LogFatal                        // Critical conditions (e.g. crash)
	LogError                        // Error conditions (e.g. failed request)
	LogWarning                      // Nothing has necessarily failed yet, but action may be needed to prevent failure
	LogNotice                       // A normal but significant event that may be worth recording
	LogInfo                         // Routine status information with no problem indicated
	LogDebug                        // Detailed diagnostic information intended for troubleshooting
	LogTrace                        // Detailed execution flow with diagnostic information. NEVER PRINTS TO SYSLOG (may expose secret information)
)

func (level LevelEnum) String() string {
	switch level {
	case LogFatal:
		return "FATAL"
	case LogError:
		return "ERROR"
	case LogWarning:
		return "WARN"
	case LogNotice:
		return "NOTICE"
	case LogInfo:
		return "INFO"
	case LogDebug:
		return "DEBUG"
	case LogTrace:
		return "TRACE"
	default:
		return strconv.Itoa(int(level))
	}
}

// [AnsiColor] sequence associated with the level.
func (level LevelEnum) AnsiColor() string {
	switch level {
	case LogFatal:
		return AnsiColor(ColorRed)
	case LogError:
		return AnsiColor(ColorBrightRed)
	case LogWarning:
		return AnsiColor(ColorBrightYellow)
	case LogNotice:
		return AnsiColor(ColorBrightBlue)
	case LogInfo:
		return AnsiColor(ColorBrightCyan)
	default:
		return AnsiColor(ColorWhite)
	}
}

type outputModeEnum uint8

const (
	// Do not insert any ANSI escape sequences
	OutputPlain outputModeEnum = iota
	// Use ANSI escape sequences for colored output
	OutputColored
	// Use ANSI escape sequences if the output is a terminal that supports ANSI
	OutputAuto
)

// Refer to [NewLogger]
type Logger struct {
	syslogWriter *syslog.Writer
	level        LevelEnum
	name         string
	printer      *Printer
}

// NewLogger creates a new logger with default settings.
func NewLogger(name string, out *os.File) *Logger {
	return &Logger{
		syslogWriter: nil,
		level:        LogError,
		name:         name,
		printer:      NewPrinter(out, OutputAuto),
	}
}

// Duplicate messages into [syslog.Writer].
func (logger *Logger) UseSyslog() error {
	syslogWriter, err := syslog.New(syslog.LOG_INFO, logger.name)
	logger.syslogWriter = syslogWriter
	return err
}

// SetLevel changes the desired level of printed messages.
func (logger *Logger) SetLevel(level LevelEnum) *Logger {
	logger.level = level
	return logger
}

// SetOutputMode changes the output color mode.
func (logger *Logger) SetOutputMode(mode outputModeEnum) *Logger {
	logger.printer.SetOutputMode(mode)
	return logger
}

// SetOutputFile changes the output file (could be [os.Stdout] or similar).
func (logger *Logger) SetOutputFile(file *os.File) *Logger {
	logger.printer.file = file
	return logger
}

// Fatal logs a message at the [LogFatal] level.
func (logger *Logger) Fatal(format string, args ...any) {
	logger.genericLog(LogFatal, format, args)
}

// Error logs a message at the [LogError] level.
func (logger *Logger) Error(format string, args ...any) {
	logger.genericLog(LogError, format, args)
}

// Warn logs a message at the [LogWarning] level.
func (logger *Logger) Warn(format string, args ...any) {
	logger.genericLog(LogWarning, format, args)
}

// Notice logs a message at the [LogNotice] level.
func (logger *Logger) Notice(format string, args ...any) {
	logger.genericLog(LogNotice, format, args)
}

// Info logs a message at the [LogInfo] level.
func (logger *Logger) Info(format string, args ...any) {
	logger.genericLog(LogInfo, format, args)
}

// Debug logs a message at the [LogDebug] level.
func (logger *Logger) Debug(format string, args ...any) {
	logger.genericLog(LogDebug, format, args)
}

// Trace logs a message at the [LogTrace] level.
// Never prints to syslog.
func (logger *Logger) Trace(format string, args ...any) {
	logger.genericLog(LogTrace, format, args)
}

func (logger *Logger) genericLog(level LevelEnum, format string, args []any) {
	if logger.shouldPrint(level) {
		message := fmt.Sprintf(format, args...)
		if logger.syslogWriter != nil {
			switch level {
			case LogFatal:
				logger.syslogWriter.Crit(message)
			case LogError:
				logger.syslogWriter.Err(message)
			case LogWarning:
				logger.syslogWriter.Warning(message)
			case LogNotice:
				logger.syslogWriter.Notice(message)
			case LogInfo:
				logger.syslogWriter.Info(message)
			case LogDebug:
				logger.syslogWriter.Debug(message)
			}
		}

		logger.write(level, message)
	}
}

// Logs an HTTP request regardless of log level.
// Never prints to syslog.
func (logger *Logger) HttpRequest(
	request *http.Request,
	request_start time.Time,
	response_code int,
) {
	logger.printer.
		Lock().
		Colored(ColorBrightWhite, logger.name).
		Colored(ColorBrightBlack, time.Now().Format(" 2006/01/02 15:04:05 ")).
		Coloredf(httpStatusCodeColor(response_code), "%6d ", response_code).
		Coloredf(ColorBrightBlack, "| %-6s | ", time.Since(request_start).Round(100*time.Millisecond)).
		Writef("%s %q\n", request.Method, request.URL).
		Unlock()
}

func (logger *Logger) write(level LevelEnum, message string) {
	logger.printer.
		Lock().
		Colored(ColorBrightWhite, logger.name).
		Colored(ColorBrightBlack, time.Now().Format(" 2006/01/02 15:04:05 ")).
		Styledf(level.AnsiColor(), "% 6s", level.String()).
		Colored(ColorBrightBlack, " | ").
		Write(message + "\n").
		Unlock()
}

func (logger *Logger) shouldPrint(level LevelEnum) bool {
	return logger.level >= level
}
