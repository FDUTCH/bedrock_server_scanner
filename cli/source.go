package cli

import (
	"bytes"
	"github.com/FDUTCH/bedrock_scanner/internal/util"
	"github.com/FDUTCH/bedrock_scanner/scanner"
	"io"
	"os"
	"strings"
)

func DeterminateSource(source string) (io.ReadCloser, error) {
	if strings.Count(source, ".") == 3 {
		return util.NewPrefixReader(source), nil
	}
	if strings.HasPrefix(source, "AS") {
		return ASSource(source), nil
	}
	return os.Open(source)
}

func ASSource(as string) io.ReadCloser {
	buff := bytes.NewBuffer(nil)
	for _, prefix := range scanner.GetRangesScraping(as) {
		buff.WriteString(prefix.String() + "\n")
	}
	return closer{buff}
}

type closer struct {
	io.Reader
}

func (c closer) Close() error {
	return nil
}
