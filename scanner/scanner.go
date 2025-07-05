package scanner

import (
	"fmt"
	"github.com/FDUTCH/bedrock_scanner/message"
	"io"
	"log/slog"
	"math"
	"net"
	"net/netip"
	"sync"
)

type Scanner struct {
	net.PacketConn
	limiter Limiter
	out     io.StringWriter
}

func NewScanner(limiter Limiter, out io.StringWriter) *Scanner {
	if limiter == nil {
		limiter = NonLimiter{}
	}

	conn, err := net.ListenPacket("udp", ":0")
	if err != nil {
		panic(err)
	}

	s := &Scanner{
		limiter:    limiter,
		PacketConn: conn,
		out:        out,
	}

	go s.listen()

	return s
}

func (s *Scanner) ScanRange(r netip.Prefix, port uint16) {
	s.Scan(parsePrefix(r, port))
}

func (s *Scanner) Scan(addr Address, r int) {
	if r > math.MaxUint32 {
		r = math.MaxUint32
	}

	msg := message.NewPingSeq()
	for range uint32(r) {
		_, err := s.WriteTo(msg, addr.Udp())
		if err != nil {
			fmt.Println(err)
			return
		}

		s.limiter.Limit()
		addr.Host += 1

		//correcting ping time.
		if r%1000 == 0 {
			message.CorrectPingTime(msg)
		}
	}
}

func (s *Scanner) ScanSync(addr Address, r int, wg *sync.WaitGroup) {
	go func() {
		s.Scan(addr, r)
		wg.Done()
	}()
}

func (s *Scanner) listen() {
	var (
		buff = make([]byte, 1492)
		pong = new(message.Pong)
	)

	for {
		n, addr, err := s.ReadFrom(buff)
		if err != nil {
			return
		}

		if n == 0 || buff[0] != 0x1c {
			s.out.WriteString(fmt.Sprintf("non-pong packet found ID=%d address=%s", buff[0], addr.String()))
			continue
		}

		if err := pong.UnmarshalBinary(buff[1:n]); err != nil {
			slog.Error(err.Error())
		}
		s.out.WriteString(fmt.Sprintf(string(pong.Data)+" address=%s", addr))
	}
}
