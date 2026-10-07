package clog

import (
	"fmt"
	"os"
	"sync"

	"github.com/mattn/go-isatty"
)

// Refer to [NewPrinter].
type Printer struct {
	file                        *os.File
	mode                        outputModeEnum
	mutex                       sync.Mutex
	cachedSupportsColoredOutput bool
}

// NewPrinter creates a new wrapper around [os.File] with specific output mode.
func NewPrinter(file *os.File, mode outputModeEnum) *Printer {
	printer := &Printer{
		file:  file,
		mode:  mode,
		mutex: sync.Mutex{},
	}
	printer.cachedSupportsColoredOutput = printer.supportsColoredOutput()
	return printer
}

// Lock MUST be called in the beginning of the write method chain
// assuming that the writer might be called from a goroutine.
//
// Close the method chain with [Printer.Unlock].
func (printer *Printer) Lock() *Printer {
	printer.mutex.Lock()
	return printer
}

// Refer to [Printer.Lock].
func (printer *Printer) Unlock() *Printer {
	printer.mutex.Unlock()
	return printer
}

func (printer *Printer) SetOutputMode(mode outputModeEnum) *Printer {
	printer.mode = mode
	printer.cachedSupportsColoredOutput = printer.supportsColoredOutput()
	return printer
}

func (printer *Printer) Write(text string) *Printer {
	printer.file.Write([]byte(text))
	return printer
}

func (printer *Printer) Writef(format string, args ...any) *Printer {
	fmt.Fprintf(printer.file, format, args...)
	return printer
}

func (printer *Printer) Writeln(line string) *Printer {
	printer.file.Write([]byte(line + "\n"))
	return printer
}

// Styled applies ANSI escape sequence on the text if the output supports it
func (printer *Printer) Styled(ansi string, text string) *Printer {
	if printer.cachedSupportsColoredOutput {
		printer.file.WriteString(ansi)
		printer.file.WriteString(text)
		printer.file.WriteString(AnsiReset)
	} else {
		fmt.Fprint(printer.file, text)
	}
	return printer
}

// Styledf applies ANSI escape sequence on the text if the output supports it
func (printer *Printer) Styledf(ansi string, format string, args ...any) *Printer {
	if printer.cachedSupportsColoredOutput {
		printer.file.WriteString(ansi)
		fmt.Fprintf(printer.file, format, args...)
		printer.file.WriteString(AnsiReset)
	} else {
		fmt.Fprintf(printer.file, format, args...)
	}
	return printer
}

// Colored is sugar code for [Printer.Styled]
func (printer *Printer) Colored(color colorCodeEnum, text string) *Printer {
	return printer.Styled(AnsiColor(color), text)
}

// Coloredf is sugar code for [Printer.Styledf]
func (printer *Printer) Coloredf(color colorCodeEnum, format string, args ...any) *Printer {
	return printer.Styledf(AnsiColor(color), format, args...)
}

func (printer *Printer) supportsColoredOutput() bool {
	switch printer.mode {
	case OutputPlain:
		return false
	case OutputColored:
		return true
	default:
		return os.Getenv("TERM") != "dumb" &&
			(isatty.IsTerminal(printer.file.Fd()) || isatty.IsCygwinTerminal(printer.file.Fd()))
	}
}
