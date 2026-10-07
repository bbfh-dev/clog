package clog

import (
	"fmt"
	"os"
)

type Printer struct {
	file *os.File
	mode outputModeEnum
}

func NewPrinter(file *os.File, mode outputModeEnum) *Printer {
	return &Printer{
		file: file,
		mode: mode,
	}
}

func (printer *Printer) Write(text string) *Printer {
	printer.file.Write([]byte(text))
	return printer
}

func (printer *Printer) Writef(format string, args ...any) *Printer {
	fmt.Fprintf(printer.file, format, args...)
	return printer
}

// FormattedWrite applies ANSI escape sequence on the text if the output supports it
func (printer *Printer) FormattedWrite(ansi string, text string) *Printer {
	if supportsColoredOutput(printer.file, printer.mode) {
		printer.file.WriteString(ansi)
		printer.file.WriteString(text)
		printer.file.WriteString(AnsiReset)
	} else {
		fmt.Fprint(printer.file, text)
	}
	return printer
}

// FormattedWritef applies ANSI escape sequence on the text if the output supports it
func (printer *Printer) FormattedWritef(ansi string, format string, args ...any) *Printer {
	if supportsColoredOutput(printer.file, printer.mode) {
		printer.file.WriteString(ansi)
		fmt.Fprintf(printer.file, format, args...)
		printer.file.WriteString(AnsiReset)
	} else {
		fmt.Fprintf(printer.file, format, args...)
	}
	return printer
}

func (printer *Printer) ColoredWrite(color colorCodeEnum, text string) *Printer {
	return printer.FormattedWrite(AnsiColor(color), text)
}

func (printer *Printer) ColoredWritef(color colorCodeEnum, format string, args ...any) *Printer {
	return printer.FormattedWritef(AnsiColor(color), format, args...)
}
