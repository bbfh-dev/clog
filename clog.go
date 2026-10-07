package clog

import (
	"os"

	"github.com/mattn/go-isatty"
)

func supportsColoredOutput(out *os.File, mode outputModeEnum) bool {
	switch mode {
	case OutputPlain:
		return false
	case OutputColored:
		return true
	default:
		return os.Getenv("TERM") != "dumb" &&
			(isatty.IsTerminal(out.Fd()) || isatty.IsCygwinTerminal(out.Fd()))
	}
}
