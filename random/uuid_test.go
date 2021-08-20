package random_test

import (
	"testing"

	zerorandom "github.com/zerogo-hub/zero-helper/random"
)

func TestUUID(t *testing.T) {
	id1 := zerorandom.NewUUID()
	id2 := zerorandom.NewUUID()

	t.Logf("id1: %s\n", id1)
	t.Logf("id2: %s\n", id2)

	if id1 == id2 {
		t.Error("Id1 and id2 cannot be the same")
	}
}
