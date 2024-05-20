package __

import "time"

func NewItem() *Item {
	id := time.Now().Unix() * 100000000
	item := &Item{
		ID:     id,
		Name:   "ITEM_NAME",
		Weight: 100,
	}

	return item
}

func NewItems(n int64) []*Item {
	id := time.Now().Unix() * 100000000
	items := make([]*Item, 0, n)

	for i := int64(1); i <= n; i++ {
		items = append(items, &Item{
			ID:     id + i,
			Name:   "ITEM_NAME",
			Weight: i * 100,
		})
	}

	return items
}
