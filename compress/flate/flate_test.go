package flate_test

import (
	"testing"

	zeroflate "github.com/zerogo-hub/zero-helper/compress/flate"
)

func TestFlate(t *testing.T) {
	s := "test flate"
	c := zeroflate.NewFlate(5)
	t.Log(c.Name())

	compressed, err := c.Compress([]byte(s))
	if err != nil {
		t.Fatal(err.Error())
		return
	}

	uncompressed, err := c.Uncompress(compressed)
	if err != nil {
		t.Fatal(err.Error())
	}

	s2 := string(uncompressed)
	if s != s2 {
		t.Fatalf("s: %s is not the same as s2: %s", s, s2)
	}
}
