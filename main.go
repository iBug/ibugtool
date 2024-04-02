package main

import (
	"os"

	"github.com/iBug/ibugtool/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
