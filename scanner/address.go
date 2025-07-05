package scanner

import (
	"fmt"
	"math"
	"net"
	"net/netip"
	"strconv"
	"strings"
	"unsafe"
)

type Address struct {
	Host uint32
	Port uint16
}

func (a *Address) Network() string {
	return "udp"
}

func (a *Address) String() string {
	arr := *(*[4]byte)(unsafe.Pointer(&a.Host))
	return fmt.Sprintf("%d.%d.%d.%d:%d", arr[3], arr[2], arr[1], arr[0], a.Port)
}

func (a *Address) Udp() *net.UDPAddr {
	var host [4]uint8
	host = *(*[4]byte)(unsafe.Pointer(&a.Host))

	return &net.UDPAddr{
		IP:   net.IP{host[3], host[2], host[1], host[0]},
		Port: int(a.Port),
	}

}

func parsePrefix(r netip.Prefix, port uint16) (Address, int) {
	startAddress, ok := parseIp(r.Addr().String())
	if !ok {
		return Address{}, 0
	}
	startAddress.Port = port

	mask, _ := strconv.Atoi(strings.Split(r.String(), "/")[1])
	return startAddress, subIPs(mask)
}

func parseIp(addr string) (result Address, ok bool) {
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

	return Address{
		Host: ip,
		Port: uint16(port),
	}, true
}

// subIPs returns ip range of the mask.
func subIPs(mask int) int {
	return int(math.Pow(2, float64(32-mask)))
}

// parseByte is parsing string into byte, way faster than strconv.Atoi.
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
