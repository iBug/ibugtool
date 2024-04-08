package version

import _ "unsafe"

//go:linkname Version main.version
var Version string = "<unknown>"
