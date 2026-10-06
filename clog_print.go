package clog

import (
	"fmt"
	"os"
)

type Writer struct {
	out  *os.File
	mode outputModeEnum
}

func NewWriter(out *os.File, mode outputModeEnum) *Writer {
	return &Writer{
		out:  out,
		mode: mode,
	}
}

func (writer *Writer) Write(text string) *Writer {
	writer.out.Write([]byte(text))
	return writer
}

func (writer *Writer) Writef(format string, args ...any) *Writer {
	fmt.Fprintf(writer.out, format, args...)
	return writer
}

// FormattedWrite applies ANSI escape sequence on the text if the output supports it
func (writer *Writer) FormattedWrite(ansi string, text string) *Writer {
	if supportsColoredOutput(writer.out, writer.mode) {
		fmt.Fprint(writer.out, ansi+text+AnsiReset)
	} else {
		fmt.Fprint(writer.out, text)
	}
	return writer
}

// FormattedWritef applies ANSI escape sequence on the text if the output supports it
func (writer *Writer) FormattedWritef(ansi string, format string, args ...any) *Writer {
	if supportsColoredOutput(writer.out, writer.mode) {
		fmt.Fprint(writer.out, ansi+fmt.Sprintf(format, args...)+AnsiReset)
	} else {
		fmt.Fprintf(writer.out, format, args...)
	}
	return writer
}

func (write *Writer) ColoredWrite(color colorCodeEnum, text string) *Writer {
	return write.FormattedWrite(AnsiColor(color), text)
}

func (write *Writer) ColoredWritef(color colorCodeEnum, format string, args ...any) *Writer {
	return write.FormattedWritef(AnsiColor(color), format, args...)
}
