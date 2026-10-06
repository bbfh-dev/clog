package clog

import (
	"fmt"
	"strconv"
)

const AnsiReset = "\033[0m"

type colorCodeEnum uint8

const (
	COLOR_BLACK colorCodeEnum = iota
	COLOR_RED
	COLOR_GREEN
	COLOR_YELLOW
	COLOR_BLUE
	COLOR_MAGENTA
	COLOR_CYAN
	COLOR_WHITE
)

const (
	COLOR_BRIGHT_BLACK colorCodeEnum = iota + 60
	COLOR_BRIGHT_RED
	COLOR_BRIGHT_GREEN
	COLOR_BRIGHT_YELLOW
	COLOR_BRIGHT_BLUE
	COLOR_BRIGHT_MAGENTA
	COLOR_BRIGHT_CYAN
	COLOR_BRIGHT_WHITE
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
		return "\x033[" + fg + "m"

	case 2:
		fg := strconv.Itoa(int(colors[0]) + 30)
		bg := strconv.Itoa(int(colors[1]) + 30)
		return "\x033[" + fg + ";" + bg + "m"

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
		return COLOR_BRIGHT_RED
	}
	if code >= 400 {
		return COLOR_BRIGHT_YELLOW
	}
	if code >= 300 {
		return COLOR_BRIGHT_BLUE
	}
	if code >= 200 {
		return COLOR_BRIGHT_GREEN
	}
	return COLOR_BRIGHT_WHITE
}
