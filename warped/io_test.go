package warped

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	// TODO: Real tests
	b := make([]byte, 3)
	var n int
	var err error

	r := Reader(strings.NewReader("ABCDE"))

	for err == nil {
		n, err = r.Read(b)
		t.Logf("%q %d %v", b, n, err)
	}

	r = Reader(strings.NewReader("ABCDE"))
	all, err := ioutil.ReadAll(r)
	t.Logf("ALL: %q %v", all, err)
}
