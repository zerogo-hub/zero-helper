package bloom_test

import (
	"testing"

	zerobloom "github.com/zerogo-hub/zero-helper/bloom"
)

func TestBloom(t *testing.T) {
	b := zerobloom.New(1000000, 0.01)

	if b.Cap() == 0 {
		t.Error("test Cap failed")
	}
	if b.K() == 0 {
		t.Error("test K failed")
	}

	b.AddString("hello")
	b.Add([]byte("world"))

	if !b.ContainsString("hello") {
		t.Error("test ContainsString failed")
	}

	if !b.Contains([]byte("world")) {
		t.Error("test Contains failed")
	}

	b.ClearAll()

	if b.ContainsString("hello") {
		t.Error("test ContainsString failed after ClearAll")
	}
}
