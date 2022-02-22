package parse

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"time"
	"unicode"
)

func FFScanner(r io.Reader) *bufio.Scanner {
	sc := bufio.NewScanner(r)
	sc.Split(func(data []byte, atEOF bool) (advance int, token []byte, spliterror error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, '\n'); i >= 0 {
			// We have a full newline-terminated line.
			return i + 1, data[0:i], nil
		}
		if i := bytes.IndexByte(data, '\r'); i >= 0 {
			// We have a cr terminated line
			return i + 1, data[0:i], nil
		}
		if atEOF {
			return len(data), data, nil
		}

		return 0, nil, nil
	})
	return sc
}

func scanLine(line string, key string, fieldFunc func(rune) bool) time.Duration {
	fields := strings.FieldsFunc(line, fieldFunc)

	var s string
	for i, f := range fields {
		if f == key && i+1 <= len(fields) {
			s = fields[i+1]
		}
	}

	d, err := duration(s)
	if err != nil {
		return 0
	}
	return d
}

func ScanDuration(line string) time.Duration {
	return scanLine(line, "Duration:", func(r rune) bool {
		return unicode.IsSpace(r) || r == ','
	})
}

func ScanTime(line string) time.Duration {
	return scanLine(line, "time", func(r rune) bool {
		return unicode.IsSpace(r) || r == '='
	})
}
