package zlib_test

import (
	"testing"

	zerozlib "github.com/zerogo-hub/zero-helper/compress/zlib"
)

func TestZlib(t *testing.T) {
	s := "test zlib"
	c := zerozlib.NewZlib(5)
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
