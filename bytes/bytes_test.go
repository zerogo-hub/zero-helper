package bytes_test

import (
	"testing"

	zerobytes "github.com/zerogo-hub/zero-helper/bytes"
)

func TestInt(t *testing.T) {
	if int8(55) != zerobytes.ToInt8(zerobytes.PutInt8(int8(55))) {
		t.Error("int8 error")
	}

	if uint8(55) != zerobytes.ToUint8(zerobytes.PutUint8(uint8(55))) {
		t.Error("uint8 error")
	}

	if int16(55) != zerobytes.ToInt16(zerobytes.PutInt16(int16(55))) {
		t.Error("int16 error")
	}

	if uint16(55) != zerobytes.ToUint16(zerobytes.PutUint16(uint16(55))) {
		t.Error("uint16 error")
	}

	if int32(55) != zerobytes.ToInt32(zerobytes.PutInt32(int32(55))) {
		t.Error("int32 error")
	}

	if uint32(55) != zerobytes.ToUint32(zerobytes.PutUint32(uint32(55))) {
		t.Error("uint32 error")
	}

	if int64(55) != zerobytes.ToInt64(zerobytes.PutInt64(int64(55))) {
		t.Error("int64 error")
	}

	if uint64(55) != zerobytes.ToUint64(zerobytes.PutUint64(uint64(55))) {
		t.Error("uint64 error")
	}
}
