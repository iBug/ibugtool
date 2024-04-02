package util

import (
	"io"
	"os"

	"golang.org/x/term"
)

const (
	CSI = "\x1B["

	ResetColor = CSI + "0m"

	ClearLine = CSI + "2K"
	ResetLine = ClearLine + "\r" + ResetColor
)

func IsWriterTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(int(f.Fd()))
}
