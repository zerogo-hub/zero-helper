package json_test

import (
	cjson "encoding/json"
	"testing"

	zerojson "github.com/zerogo-hub/zero-helper/json"
	zerotest "github.com/zerogo-hub/zero-helper/test"
)

func TestJSON(t *testing.T) {
	item := zerotest.NewItem()

	itemBytes, err := zerojson.Marshal(item)
	if err != nil {
		t.Fatalf("zerojson.Marshal failed: %s", err.Error())
	}

	var newItem *zerotest.Item
	err = zerojson.Unmarshal(itemBytes, &newItem)
	if err != nil {
		t.Fatalf("zerojson.Unmarshal failed: %s", err.Error())
	}

	if newItem.ID != item.ID {
		t.Fatal("zerojson.Unmarshal invalid id")
	}
}

func BenchmarkMarshal_1_StdLib(b *testing.B) {
	items := zerotest.NewItems(1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cjson.Marshal(items)
	}
}

func BenchmarkMarshal_1_ZeroJSON(b *testing.B) {
	items := zerotest.NewItems(1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = zerojson.Marshal(items)
	}
}

func BenchmarkUnmarshal_1_StdLib(b *testing.B) {
	itemsStr, _ := cjson.Marshal(zerotest.NewItems(1))
	itemsBytes := []byte(itemsStr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cjson.Marshal(itemsBytes)
	}
}

func BenchmarkUnmarshal_1_ZeroJSON(b *testing.B) {
	itemsStr, _ := cjson.Marshal(zerotest.NewItems(1))
	itemsBytes := []byte(itemsStr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = zerojson.Marshal(itemsBytes)
	}
}

func BenchmarkMarshal_100_StdLib(b *testing.B) {
	items := zerotest.NewItems(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cjson.Marshal(items)
	}
}

func BenchmarkMarshal_100_ZeroJSON(b *testing.B) {
	items := zerotest.NewItems(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = zerojson.Marshal(items)
	}
}

func BenchmarkUnmarshal_100_StdLib(b *testing.B) {
	itemsStr, _ := cjson.Marshal(zerotest.NewItems(100))
	itemsBytes := []byte(itemsStr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cjson.Marshal(itemsBytes)
	}
}

func BenchmarkUnmarshal_100_ZeroJSON(b *testing.B) {
	itemsStr, _ := cjson.Marshal(zerotest.NewItems(100))
	itemsBytes := []byte(itemsStr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = zerojson.Marshal(itemsBytes)
	}
}

func BenchmarkMarsha_500_StdLib(b *testing.B) {
	items := zerotest.NewItems(500)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cjson.Marshal(items)
	}
}

func BenchmarkMarshal_500_ZeroJSON(b *testing.B) {
	items := zerotest.NewItems(500)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = zerojson.Marshal(items)
	}
}

func BenchmarkUnmarshal_500_StdLib(b *testing.B) {
	itemsStr, _ := cjson.Marshal(zerotest.NewItems(500))
	itemsBytes := []byte(itemsStr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cjson.Marshal(itemsBytes)
	}
}

func BenchmarkUnmarshal_500_ZeroJSON(b *testing.B) {
	itemsStr, _ := cjson.Marshal(zerotest.NewItems(500))
	itemsBytes := []byte(itemsStr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = zerojson.Marshal(itemsBytes)
	}
}
