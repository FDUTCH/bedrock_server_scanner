package cli

import (
	"flag"
	"github.com/FDUTCH/bedrock_scanner/gui"
	"github.com/FDUTCH/bedrock_scanner/scanner"
	"log/slog"
	"os"
)

var (
	useGui   bool
	source   string
	sockets  int
	pps      int
	ppsLimit bool
	port     int
)

func Run() {
	flag.BoolVar(&useGui, "g", false, "run program in cli")
	flag.StringVar(&source, "s", "0.0.0.0/0", "input source")
	flag.IntVar(&sockets, "t", 1, "number of sockets/threads")
	flag.IntVar(&pps, "pps", 10000, "number of packets per second")
	flag.BoolVar(&ppsLimit, "l", true, "packets per second limit")
	flag.IntVar(&port, "p", 19132, "port")
	flag.Parse()

	if useGui {
		gui.Run()
		return
	}
	src, err := DeterminateSource(source)
	if err != nil {
		slog.Error("unable to scan", "err", err, "source", source)
		os.Exit(1)
	}
	settings := scanner.Settings{
		Source:           src,
		Sockets:          sockets,
		PacketsPerSecond: pps,
		NoPPSLimit:       !ppsLimit,
		Port:             port,
		Out:              output{},
	}
	settings.Scan()
}
