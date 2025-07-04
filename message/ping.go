package message

import (
	"encoding/binary"
	"io"
	"math/rand/v2"
	"time"
)

func NewPingSeq() []byte {
	unconnectedPing := unconnectedPing{
		PingTime:   timestamp(),
		ClientGUID: int64(id),
	}
	data, err := unconnectedPing.Marshal()
	if err != nil {
		panic(err)
	}
	return data
}

var id = rand.Uint64()

// timestamp returns a timestamp in milliseconds.
func timestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

var unconnectedMessageSequence = [16]byte{0x00, 0xff, 0xff, 0x00, 0xfe, 0xfe, 0xfe, 0xfe, 0xfd, 0xfd, 0xfd, 0xfd, 0x12, 0x34, 0x56, 0x78}

type unconnectedPing struct {
	PingTime   int64
	ClientGUID int64
}

func (pk *unconnectedPing) Marshal() (data []byte, err error) {
	b := make([]byte, 33)
	b[0] = 1
	binary.BigEndian.PutUint64(b[1:], uint64(pk.PingTime))
	copy(b[9:], unconnectedMessageSequence[:])
	binary.BigEndian.PutUint64(b[25:], uint64(pk.ClientGUID))
	return b, nil
}

type Pong struct {
	// PingTime is filled out using unconnectedPing.PingTime.
	PingTime   int64
	ServerGUID int64
	Data       []byte
}

func (pk *Pong) UnmarshalBinary(data []byte) error {
	if len(data) < 34 || len(data) < 34+int(binary.BigEndian.Uint16(data[32:])) {
		return io.ErrUnexpectedEOF
	}
	pk.PingTime = int64(binary.BigEndian.Uint64(data))
	pk.ServerGUID = int64(binary.BigEndian.Uint64(data[8:]))
	// Magic: 16 bytes.
	n := binary.BigEndian.Uint16(data[32:])
	pk.Data = append([]byte(nil), data[34:34+n]...)
	return nil
}

func CorrectPingTime(data []byte) {
	binary.BigEndian.PutUint64(data[1:], uint64(timestamp()))
}
