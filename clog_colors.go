package clog

import (
	"fmt"
	"strconv"
)

const AnsiReset = "\033[0m"

type colorCodeEnum uint8

const (
	ColorBlack colorCodeEnum = iota
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)

const (
	ColorBrightBlack colorCodeEnum = iota + 60
	ColorBrightRed
	ColorBrightGreen
	ColorBrightYellow
	ColorBrightBlue
	ColorBrightMagenta
	ColorBrightCyan
	ColorBrightWhite
)

// AnsiColor creates an ansi escape sequence for setting color.
//
//   - Returns [AnsiReset] when no arguments.
//   - Allows 1-2 arguments (set foreground & background).
//   - Panics otherwise.
func AnsiColor(colors ...colorCodeEnum) string {
	switch len(colors) {
	case 0:
		return AnsiReset

	case 1:
		fg := strconv.Itoa(int(colors[0]) + 30)
		return "\x1b[" + fg + "m"

	case 2:
		fg := strconv.Itoa(int(colors[0]) + 30)
		bg := strconv.Itoa(int(colors[1]) + 30)
		return "\x1b[" + fg + ";" + bg + "m"

	default:
		panic(fmt.Sprintf("Invalid number of arguments in a call to AnsiColor(%#v)", colors))
	}
}

// AnsiColorFg creates a 24-bit true-color foreground color.
//
// rgb must be packed as 0xRRGGBB.
func AnsiColorFg(rgb uint32) string {
	r := (rgb >> 16) & 0xff
	g := (rgb >> 8) & 0xff
	b := rgb & 0xff

	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
}

// AnsiColorBg sets a 24-bit true-color background color.
//
// rgb must be packed as 0xRRGGBB.
func AnsiColorBg(rgb uint32) string {
	r := (rgb >> 16) & 0xff
	g := (rgb >> 8) & 0xff
	b := rgb & 0xff

	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
}

// AnsiColorFg8 sets an 8-bit foreground color.
func AnsiColorFg8(n uint8) string {
	return fmt.Sprintf("\x1b[38;5;%dm", n)
}

// AnsiColorBg8 sets an 8-bit background color.
func AnsiColorBg8(n uint8) string {
	return fmt.Sprintf("\x1b[48;5;%dm", n)
}

func httpStatusCodeColor(code int) colorCodeEnum {
	if code >= 500 {
		return ColorBrightRed
	}
	if code >= 400 {
		return ColorBrightYellow
	}
	if code >= 300 {
		return ColorBrightBlue
	}
	if code >= 200 {
		return ColorBrightGreen
	}
	return ColorBrightWhite
}
