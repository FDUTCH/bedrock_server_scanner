package cli

import "os"

type output struct{}

func (o output) WriteString(s string) (n int, err error) {
	return os.Stdout.WriteString(s + "\n")
}
