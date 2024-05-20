package random_test

import (
	"strconv"
	"testing"

	zerorandom "github.com/zerogo-hub/zero-helper/random"
)

type Item struct {
	id     int
	name   string
	weight int
}

func TestChoicePick(t *testing.T) {
	n := 10
	items := make([]Item, 0, n)
	for i := 1; i <= n; i++ {
		items = append(items, Item{id: i, name: "NAME_" + strconv.Itoa(i), weight: i * 100})
	}

	c, err := zerorandom.NewChooser(items, func(item Item) int {
		return item.weight
	})
	if err != nil {
		t.Fatalf("build chooser failed, err: %s", err.Error())
	}

	for i := 0; i < len(items); i++ {
		_ = c.Pick()
	}
}

func TestChoicePickOnce(t *testing.T) {
	n := 10
	items := make([]Item, 0, n)
	for i := 1; i <= n; i++ {
		items = append(items, Item{id: i, name: "NAME_" + strconv.Itoa(i), weight: i * 100})
	}

	_ = zerorandom.PickOnce(items, func(item Item) int { return item.weight })
}
