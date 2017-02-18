// Package hexdump helps reading hexdump stream.
package hexdump

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"io"
	"regexp"
	"strings"
	"unicode"
)

// NewReader reads hexdump from r.
func NewReader(r io.Reader) io.Reader {
	return &dumpReader{
		scanner: bufio.NewScanner(r),
	}
}

// NewReaderString reads hexdump from s.
func NewReaderString(s string) io.Reader {
	return NewReader(strings.NewReader(s))
}

type dumpReader struct {
	scanner *bufio.Scanner
	cache   []byte
}

func (r *dumpReader) Read(b []byte) (int, error) {
	var buf []byte
	var err error

	if len(r.cache) > 0 {
		buf = r.cache
	} else {
		buf, err = r.decode()
	}

	n := copy(b, buf)
	r.cache = buf[n:]
	return n, err
}

var pattern = regexp.MustCompile(`(?i:\s+[0-9a-f]{2}){1,16}`)

func (r *dumpReader) decode() ([]byte, error) {
	var res []byte

	if !r.scanner.Scan() {
		return res, io.EOF
	}
	if err := r.scanner.Err(); err != nil {
		return res, err
	}

	// read line
	b := r.scanner.Bytes()

	// extract hexa chars
	b = pattern.Find(b)

	// strip spaces
	b = bytes.Map(isSpace, b)

	// decode
	res = make([]byte, hex.DecodedLen(len(b)))
	n, err := hex.Decode(res, b)
	return res[0:n], err
}

func isSpace(r rune) rune {
	if unicode.IsSpace(r) {
		return -1
	}
	return r
}
