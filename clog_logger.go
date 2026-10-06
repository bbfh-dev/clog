package clog

import (
	"fmt"
	"log/syslog"
	"os"

	"github.com/mattn/go-isatty"
)

type LevelEnum uint8

const (
	LOG_SILENT  LevelEnum = iota * 4 // Highest priority
	LOG_FATAL                        // Critical conditions (e.g. crash)
	LOG_ERROR                        // Error conditions (e.g. failed request)
	LOG_WARNING                      // Nothing has necessarily failed yet, but action may be needed to prevent failure
	LOG_NOTICE                       // A normal but significant event that may be worth recording
	LOG_INFO                         // Routine status information with no problem indicated
	LOG_DEBUG                        // Detailed diagnostic information intended for troubleshooting
	LOG_TRACE                        // Detailed execution flow with diagnostic information. NEVER PRINTS TO SYSLOG (may expose secret information)
)

func (level LevelEnum) String() string {
	switch level {
	case LOG_FATAL:
		return "FATAL"
	case LOG_ERROR:
		return "ERROR"
	case LOG_WARNING:
		return "WARN"
	case LOG_NOTICE:
		return "NOTICE"
	case LOG_INFO:
		return "INFO"
	case LOG_DEBUG:
		return "DEBUG"
	case LOG_TRACE:
		return "TRACE"
	default:
		return fmt.Sprintf("%d", level)
	}
}

func (level LevelEnum) AnsiColor() string {
	switch level {
	case LOG_FATAL:
		return AnsiColor(COLOR_RED)
	case LOG_ERROR:
		return AnsiColor(COLOR_BRIGHT_RED)
	case LOG_WARNING:
		return AnsiColor(COLOR_BRIGHT_YELLOW)
	case LOG_NOTICE:
		return AnsiColor(COLOR_BRIGHT_BLUE)
	case LOG_INFO:
		return AnsiColor(COLOR_BRIGHT_CYAN)
	default:
		return AnsiColor(COLOR_WHITE)
	}
}

type outputModeEnum uint8

const (
	OUTPUT_PLAIN outputModeEnum = iota
	OUTPUT_COLORED
	OUTPUT_AUTO
)

// Refer to [NewLogger]
type Logger struct {
	syslogWriter *syslog.Writer
	level        LevelEnum
	outputMode   outputModeEnum
	name         string
	out          *os.File
}

// NewLogger creates a new logger with default settings.
func NewLogger(name string, out *os.File) *Logger {
	return &Logger{
		syslogWriter: nil,
		level:        LOG_ERROR,
		outputMode:   OUTPUT_AUTO,
		name:         name,
		out:          out,
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
	logger.outputMode = mode
	return logger
}

// Fatal logs a message at the [LOG_FATAL] level.
func (logger *Logger) Fatal(format string, args ...any) {
	level := LOG_FATAL
	if logger.shouldPrint(level) {
		message := fmt.Sprintf(format, args...)
		if logger.syslogWriter != nil {
			logger.syslogWriter.Crit(message) // ignore error
		}

		logger.write(level, message)
	}
}

// Error logs a message at the [LOG_ERROR] level.
func (logger *Logger) Error(format string, args ...any) {
	level := LOG_ERROR
	if logger.shouldPrint(level) {
		message := fmt.Sprintf(format, args...)
		if logger.syslogWriter != nil {
			logger.syslogWriter.Err(message) // ignore error
		}

		logger.write(level, message)
	}
}

// Warn logs a message at the [LOG_WARNING] level.
func (logger *Logger) Warn(format string, args ...any) {
	level := LOG_WARNING
	if logger.shouldPrint(level) {
		message := fmt.Sprintf(format, args...)
		if logger.syslogWriter != nil {
			logger.syslogWriter.Warning(message) // ignore error
		}

		logger.write(level, message)
	}
}

// Notice logs a message at the [LOG_NOTICE] level.
func (logger *Logger) Notice(format string, args ...any) {
	level := LOG_NOTICE
	if logger.shouldPrint(level) {
		message := fmt.Sprintf(format, args...)
		if logger.syslogWriter != nil {
			logger.syslogWriter.Notice(message) // ignore error
		}

		logger.write(level, message)
	}
}

// Info logs a message at the [LOG_INFO] level.
func (logger *Logger) Info(format string, args ...any) {
	level := LOG_INFO
	if logger.shouldPrint(level) {
		message := fmt.Sprintf(format, args...)
		if logger.syslogWriter != nil {
			logger.syslogWriter.Info(message) // ignore error
		}

		logger.write(level, message)
	}
}

// Debug logs a message at the [LOG_DEBUG] level.
func (logger *Logger) Debug(format string, args ...any) {
	level := LOG_DEBUG
	if logger.shouldPrint(level) {
		message := fmt.Sprintf(format, args...)
		if logger.syslogWriter != nil {
			logger.syslogWriter.Debug(message) // ignore error
		}

		logger.write(level, message)
	}
}

// Trace logs a message at the [LOG_TRACE] level.
// Never prints to [Logger.syslogWriter].
func (logger *Logger) Trace(format string, args ...any) {
	level := LOG_TRACE
	if logger.shouldPrint(level) {
		logger.write(level, fmt.Sprintf(format, args...))
	}
}

func (logger *Logger) write(level LevelEnum, message string) {
	if logger.supportsColoredOutput() {
		fmt.Fprintf(
			logger.out,
			"%s%s %s%s %s| %s%s\n",
			AnsiColor(COLOR_BRIGHT_WHITE),
			logger.name,
			level.AnsiColor(),
			level,
			AnsiColor(COLOR_WHITE),
			AnsiReset,
			message,
		)
	} else {
		fmt.Fprintf(logger.out, "%s %s | %s\n", logger.name, level, message)
	}
}

func (logger *Logger) shouldPrint(level LevelEnum) bool {
	return logger.level >= level
}

func (logger *Logger) supportsColoredOutput() bool {
	return os.Getenv("TERM") != "dumb" &&
		(isatty.IsTerminal(logger.out.Fd()) || isatty.IsCygwinTerminal(logger.out.Fd()))
}
