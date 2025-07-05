package scanner

import (
	"bufio"
	"io"
	"net/netip"
	"sync"
)

type Settings struct {
	Source           io.ReadCloser
	Sockets          int
	PacketsPerSecond int
	NoPPSLimit       bool
	Port             int

	Out io.StringWriter
}

func (s Settings) Scan() {
	defer s.Source.Close()
	scanners := make([]*Scanner, 0, s.Sockets)
	for range s.Sockets {
		var limiter Limiter
		if !s.NoPPSLimit {
			limiter = NewLimiter(uint64(s.PacketsPerSecond / s.Sockets))
		}
		scanners = append(scanners, NewScanner(limiter, s.Out))
	}

	sourceScanner := bufio.NewScanner(s.Source)

	for sourceScanner.Scan() {
		text := sourceScanner.Text()
		prefix, err := netip.ParsePrefix(text)
		if err != nil {
			continue
		}
		handlePrefix(scanners, prefix, s.Port)
	}
}

func handlePrefix(scanners []*Scanner, prefix netip.Prefix, port int) {
	var wg sync.WaitGroup
	baseAddr, r := parsePrefix(prefix, uint16(port))
	wg.Add(len(scanners))
	taskCount := r / len(scanners)
	rest := r % len(scanners)
	for i, scanner := range scanners {
		if i == len(scanners)-1 {
			taskCount += rest
		}
		scanner.ScanSync(baseAddr, taskCount, &wg)
		baseAddr.Host += uint32(taskCount)
	}
	wg.Wait()
}
