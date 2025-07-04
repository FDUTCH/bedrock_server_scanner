package scanner

import (
	"github.com/FDUTCH/bedrock_scanner/message"
	"log/slog"
	"net"
	"net/netip"
	"strconv"
	"strings"
)

type Scanner struct {
	net.PacketConn
	limiter Limiter
	logger  *slog.Logger
}

func NewScanner(limiter Limiter, logger *slog.Logger) *Scanner {
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
		logger:     logger,
	}

	go s.listen()

	return s
}

func (s *Scanner) ScanRange(r netip.Prefix, port uint16) {
	startAddress, ok := parseIp(r.Addr().String())
	if !ok {
		return
	}
	startAddress.port = port

	mask, _ := strconv.Atoi(strings.Split(r.String(), "/")[1])

	s.scan(startAddress, subIPs(mask))
}

func (s *Scanner) scan(addr *address, r int) {
	msg := message.NewPingSeq()
	for i := range uint32(r) {
		s.limiter.Limit()
		addr.host += i
		_, _ = s.WriteTo(msg, addr.Udp())

		//correcting ping time.
		if r%1000 == 0 {
			message.CorrectPingTime(msg)
		}
	}
}

func (s *Scanner) listen() {

	buff := make([]byte, 1492)

	pong := &message.Pong{}

	for {
		n, addr, err := s.ReadFrom(buff)
		if err != nil {
			return
		}
		if n == 0 || buff[0] != 0x1c {
			s.logger.Error("non-pong packet found", "ID", buff[0], "address", addr)
			continue
		}
		if err := pong.UnmarshalBinary(buff[1:n]); err != nil {
			slog.Error(err.Error())
		}
		s.logger.Info(string(pong.Data), "address", addr)
	}
}
