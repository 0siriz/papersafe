package paper

import (
	"encoding/base32"
	"strings"
)

const (
	segmentSize    = 5
	segmentSpacing = 1
)

var zbase32 = base32.NewEncoding("ybndrfg8ejkmcpqxot1uwisza345h769").WithPadding(base32.NoPadding)

func segment(text string) string {
	var result strings.Builder

	for i := 0; i < len(text); i += segmentSize {
		end := min(i+segmentSize, len(text))

		result.WriteString(text[i:end])
		result.WriteByte(0x20)
	}

	return result.String()
}
