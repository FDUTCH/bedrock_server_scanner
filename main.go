package main

import (
	"github.com/FDUTCH/bedrock_scanner/cli"
	"github.com/FDUTCH/bedrock_scanner/gui"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		cli.Run()
		return
	}
	gui.Run()
}
