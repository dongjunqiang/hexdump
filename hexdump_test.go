package hexdump_test

import (
	"encoding/hex"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/schorlet/hexdump"
)

func Example() {
	b := []byte("Emma et Dexter se rencontrent pour la première fois")
	r := strings.NewReader(hex.Dump(b))
	xd := hexdump.NewReader(r)
	io.Copy(os.Stdout, xd)
	// Output:
	// Emma et Dexter se rencontrent pour la première fois
}

func TestShortRead(t *testing.T) {
	var want [26]byte
	for i := range want {
		want[i] = byte(i + 65)
	}

	dump := hex.Dump(want[:])
	reader := hexdump.NewReaderString(dump)

	got := make([]byte, 1)
	count := 0

	for i := range want {
		_, err := reader.Read(got)
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(err)
		}

		if got[0] != want[i] {
			t.Fatalf("got: %c want: %c", got[0], want[i])
		}
		count++

	}

	if count != 26 {
		t.Fatalf("got: %d want: 26", count)
	}
}
