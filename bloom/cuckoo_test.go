package bloom_test

import (
	"testing"

	zerobloom "github.com/zerogo-hub/zero-helper/bloom"
)

func TestCuckoo(t *testing.T) {
	c := zerobloom.NewCuckoo(10000)

	if c.Count() != 0 {
		t.Error("test Count failed")
	}

	c.AddString("hello")
	c.Add([]byte("world"))

	if c.Count() != 2 {
		t.Error("test Count failed")
	}

	if !c.ContainsString("hello") {
		t.Error("test ContainsString failed")
	}

	if !c.Contains([]byte("world")) {
		t.Error("test Contains failed")
	}

	c.DelString("world")

	if c.Contains([]byte("world")) {
		t.Error("test Contains failed")
	}

	c.ClearAll()

	if c.ContainsString("hello") {
		t.Error("test ContainsString failed after ClearAll")
	}

	c.AddString("abcdef")
	c.Del([]byte("abcdef"))
	if c.ContainsString("abcdef") {
		t.Error("test Del failed")
	}
}
