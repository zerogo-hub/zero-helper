package json_test

import (
	cjson "encoding/json"
	"testing"
	"time"

	zerojson "github.com/zerogo-hub/zero-helper/json"
)

type Item struct {
	ID   int64
	Name string
}

func newItem() *Item {
	id := time.Now().Unix() * 100000000
	item := &Item{
		ID:   id,
		Name: "ITEM_NAME",
	}

	return item
}

func newItems(n int64) []*Item {
	id := time.Now().Unix() * 100000000
	items := make([]*Item, 0, n)

	for i := int64(1); i <= n; i++ {
		items = append(items, &Item{
			ID:   id + i,
			Name: "ITEM_NAME",
		})
	}

	return items
}

func TestJSON(t *testing.T) {
	item := newItem()

	itemBytes, err := zerojson.Marshal(item)
	if err != nil {
		t.Fatalf("zerojson.Marshal failed: %s", err.Error())
	}

	var newItem *Item
	err = zerojson.Unmarshal(itemBytes, &newItem)
	if err != nil {
		t.Fatalf("zerojson.Unmarshal failed: %s", err.Error())
	}

	if newItem.ID != item.ID {
		t.Fatal("zerojson.Unmarshal invalid id")
	}
}

func BenchmarkMarshal_1_StdLib(b *testing.B) {
	items := newItems(1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cjson.Marshal(items)
	}
}

func BenchmarkMarshal_1_ZeroJSON(b *testing.B) {
	items := newItems(1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = zerojson.Marshal(items)
	}
}

func BenchmarkUnmarshal_1_StdLib(b *testing.B) {
	itemsStr, _ := cjson.Marshal(newItems(1))
	itemsBytes := []byte(itemsStr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cjson.Marshal(itemsBytes)
	}
}

func BenchmarkUnmarshal_1_ZeroJSON(b *testing.B) {
	itemsStr, _ := cjson.Marshal(newItems(1))
	itemsBytes := []byte(itemsStr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = zerojson.Marshal(itemsBytes)
	}
}

func BenchmarkMarshal_100_StdLib(b *testing.B) {
	items := newItems(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cjson.Marshal(items)
	}
}

func BenchmarkMarshal_100_ZeroJSON(b *testing.B) {
	items := newItems(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = zerojson.Marshal(items)
	}
}

func BenchmarkUnmarshal_100_StdLib(b *testing.B) {
	itemsStr, _ := cjson.Marshal(newItems(100))
	itemsBytes := []byte(itemsStr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cjson.Marshal(itemsBytes)
	}
}

func BenchmarkUnmarshal_100_ZeroJSON(b *testing.B) {
	itemsStr, _ := cjson.Marshal(newItems(100))
	itemsBytes := []byte(itemsStr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = zerojson.Marshal(itemsBytes)
	}
}

func BenchmarkMarsha_500_StdLib(b *testing.B) {
	items := newItems(500)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cjson.Marshal(items)
	}
}

func BenchmarkMarshal_500_ZeroJSON(b *testing.B) {
	items := newItems(500)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = zerojson.Marshal(items)
	}
}

func BenchmarkUnmarshal_500_StdLib(b *testing.B) {
	itemsStr, _ := cjson.Marshal(newItems(500))
	itemsBytes := []byte(itemsStr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cjson.Marshal(itemsBytes)
	}
}

func BenchmarkUnmarshal_500_ZeroJSON(b *testing.B) {
	itemsStr, _ := cjson.Marshal(newItems(500))
	itemsBytes := []byte(itemsStr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = zerojson.Marshal(itemsBytes)
	}
}
