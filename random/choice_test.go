package random_test

import (
	"testing"

	zerorandom "github.com/zerogo-hub/zero-helper/random"
	zerotest "github.com/zerogo-hub/zero-helper/test"
)

func TestChoicePick(t *testing.T) {
	_, err := zerorandom.NewChooser([]*zerotest.Item{}, func(item *zerotest.Item) int64 {
		return item.Weight
	})
	if err == nil {
		t.Fatal("err cant be nil")
	}

	items := zerotest.NewItems(10)

	c, err := zerorandom.NewChooser(items, func(item *zerotest.Item) int64 {
		return item.Weight
	})
	if err != nil {
		t.Fatalf("build chooser failed, err: %s", err.Error())
	}

	for i := 0; i < len(items); i++ {
		_ = c.Pick()
	}
}

func TestChoicePickOnce(t *testing.T) {
	items := zerotest.NewItems(10)
	_ = zerorandom.PickOnce(items, func(item *zerotest.Item) int64 { return item.Weight })
}
