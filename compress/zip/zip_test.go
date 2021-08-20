package zip_test

import (
	"testing"

	zerozip "github.com/zerogo-hub/zero-helper/compress/zip"
)

func TestZip(t *testing.T) {
	s := "test zip"
	c := zerozip.NewZip()

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
