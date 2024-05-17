package codec_test

import (
	"testing"

	zerotest "github.com/zerogo-hub/zero-helper/test"

	libJSON "encoding/json"

	zerocjson "github.com/zerogo-hub/zero-helper/codec/json"
	zerocmsgpack "github.com/zerogo-hub/zero-helper/codec/msgpack"
	zerocprotobuf "github.com/zerogo-hub/zero-helper/codec/protobuf"
)

func Test_JSON(t *testing.T) {
	items := zerotest.NewItems(3)

	c := zerocjson.New()

	itemBytes, err := c.Marshal(items)
	if err != nil {
		t.Fatal("TestMarshal_JSON Marshal failed")
	}

	var newItems []*zerotest.Item
	if err := c.Unmarshal(itemBytes, &newItems); err != nil {
		t.Fatal("TestMarshal_JSON Unmarshal failed")
	}

	if len(newItems) != len(items) {
		t.Fatal("TestMarshal_JSON Unmarshal failed, invalid len")
	}

	if len(c.Name()) == 0 || len(c.MimeType()) == 0 {
		t.Fatal("TestMarshal_JSON invalid information")
	}
}

func Test_Msgpack(t *testing.T) {
	items := zerotest.NewItems(1)

	c := zerocmsgpack.New()

	itemBytes, err := c.Marshal(items)
	if err != nil {
		t.Fatal("Test_Msgpack Marshal failed")
	}

	var newItems []*zerotest.Item
	if err := c.Unmarshal(itemBytes, &newItems); err != nil {
		t.Fatal("Test_Msgpack Unmarshal failed")
	}

	if len(newItems) != len(items) {
		t.Fatal("Test_Msgpack Unmarshal failed, invalid len")
	}

	if len(c.Name()) == 0 || len(c.MimeType()) == 0 {
		t.Fatal("Test_Msgpack invalid information")
	}
}

func Test_Protobuf(t *testing.T) {
	item := zerotest.NewItem()

	c := zerocprotobuf.New()

	itemBytes, err := c.Marshal(item)
	if err != nil {
		t.Fatal("Test_Protobuf Marshal failed")
	}

	newItems := &zerotest.Item{}
	if err := c.Unmarshal(itemBytes, newItems); err != nil {
		t.Fatal("Test_Protobuf Unmarshal failed")
	}

	if len(c.Name()) == 0 || len(c.MimeType()) == 0 {
		t.Fatal("Test_Protobuf invalid information")
	}
}

func BenchmarkMarshal_LibJSON(b *testing.B) {
	item := zerotest.NewItem()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = libJSON.Marshal(item)
	}
}

func BenchmarkMarshal_ZeroJSON(b *testing.B) {
	item := zerotest.NewItem()
	b.ResetTimer()

	c := zerocjson.New()
	for i := 0; i < b.N; i++ {
		_, _ = c.Marshal(item)
	}
}

func BenchmarkMarshal_MsgPack(b *testing.B) {
	item := zerotest.NewItem()
	b.ResetTimer()

	c := zerocmsgpack.New()
	for i := 0; i < b.N; i++ {
		_, _ = c.Marshal(item)
	}
}

func BenchmarkMarshal_Protobuf(b *testing.B) {
	item := zerotest.NewItem()
	b.ResetTimer()

	c := zerocprotobuf.New()
	for i := 0; i < b.N; i++ {
		_, _ = c.Marshal(item)
	}
}

func BenchmarkUnmarshal_LibJSON(b *testing.B) {
	itemBytes, _ := libJSON.Marshal(zerotest.NewItem())
	b.ResetTimer()

	var item zerotest.Item
	for i := 0; i < b.N; i++ {
		_ = libJSON.Unmarshal(itemBytes, &item)
	}
}

func BenchmarkUnmarshal_ZeroJSON(b *testing.B) {
	itemBytes, _ := libJSON.Marshal(zerotest.NewItem())
	b.ResetTimer()

	c := zerocjson.New()
	var item zerotest.Item
	for i := 0; i < b.N; i++ {
		_ = c.Unmarshal(itemBytes, &item)
	}
}

func BenchmarkUnmarshal_Msgpack(b *testing.B) {
	itemBytes, _ := libJSON.Marshal(zerotest.NewItem())
	b.ResetTimer()

	c := zerocmsgpack.New()
	var item zerotest.Item
	for i := 0; i < b.N; i++ {
		_ = c.Unmarshal(itemBytes, &item)
	}
}

func BenchmarkUnmarshal_Protobuf(b *testing.B) {
	itemBytes, _ := libJSON.Marshal(zerotest.NewItem())
	b.ResetTimer()

	c := zerocprotobuf.New()
	var item zerotest.Item
	for i := 0; i < b.N; i++ {
		_ = c.Unmarshal(itemBytes, &item)
	}
}
