package scanner

import (
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
	"unsafe"
)

type address struct {
	host uint32
	port uint16
	buff *strings.Builder
}

func (a *address) Network() string {
	return "udp"
}

func (a *address) String() string {
	arr := *(*[4]uint8)(unsafe.Pointer(&a.host))
	return fmt.Sprintf("%d.%d.%d.%d:%d", arr[3], arr[2], arr[1], arr[0], a.port)
}

func (a *address) Udp() *net.UDPAddr {
	var host [4]uint8
	host = *(*[4]byte)(unsafe.Pointer(&a.host))

	return &net.UDPAddr{
		IP:   net.IP{host[3], host[2], host[1], host[0]},
		Port: int(a.port),
	}

}

func parseIp(addr string) (result *address, ok bool) {
	defer func() {
		if err := recover(); err != nil {
			ok = false
		}
	}()
	portAddr := strings.Split(addr, ":")

	arr := strings.Split(portAddr[0], ".")
	var port int

	if len(portAddr) > 1 {
		port, _ = strconv.Atoi(portAddr[1])
	}

	var a [4]uint8
	for i := range a {
		num := parseByte(arr[i])
		a[3-i] = num
	}

	var ip = *(*uint32)(unsafe.Pointer(&a))

	return &address{
		host: ip,
		port: uint16(port),
		buff: &strings.Builder{},
	}, true
}

func subIPs(mask int) int {
	return int(math.Pow(2, float64(32-mask)))
}

func parseByte(str string) byte {
	var result byte
	var multiplier byte = 1
	for i := range str {
		index := len(str) - (i + 1)
		val := str[index] - '0'
		result += val * multiplier
		multiplier *= 10
	}
	return result
}
