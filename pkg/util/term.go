package util

const (
	CSI = "\x1B["

	ResetColor = CSI + "0m"

	ClearLine = CSI + "2K"
	ResetLine = ClearLine + "\r" + ResetColor
)
