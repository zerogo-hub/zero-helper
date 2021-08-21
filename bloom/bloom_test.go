package bloom_test

import (
	"testing"

	zerobloom "github.com/zerogo-hub/zero-helper/bloom"
)

func TestBloom(t *testing.T) {
	b := zerobloom.New(1000000, 0.01)

	b.AddString("hello")

	if !b.ContainsString("hello") {
		t.Error("test ContainsString failed")
	}
}
