package scanner

import "io"

type Settings struct {
	Source           io.ReadCloser
	Sockets          int
	PacketsPerSecond int
	NoPPSLimit       bool

	Out io.StringWriter
}
