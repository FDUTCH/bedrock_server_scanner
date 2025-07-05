package util

import "strings"

type SinglePrefixReader struct {
	reader *strings.Reader
}

func NewPrefixReader(str string) *SinglePrefixReader {
	return &SinglePrefixReader{reader: strings.NewReader(str)}
}

func (p *SinglePrefixReader) Read(data []byte) (n int, err error) {
	return p.reader.Read(data)
}

func (*SinglePrefixReader) Close() error {
	return nil
}
