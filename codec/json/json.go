package json

import (
	"encoding/json"

	zerocodec "github.com/zerogo-hub/zero-helper/codec"
)

type jsonCodec struct{}

// NewJSONCodec JSON
func NewJSONCodec() zerocodec.Codec {
	return &jsonCodec{}
}

// Marshal 编码
func (*jsonCodec) Marshal(in interface{}) ([]byte, error) {
	return json.Marshal(in)
}

// Unmarshal 解码
func (*jsonCodec) Unmarshal(in []byte, out interface{}) error {
	return json.Unmarshal(in, out)
}

// Name 名称
func (*jsonCodec) Name() string {
	return "json"
}

// MimeType 媒体类型
func (*jsonCodec) MimeType() string {
	return "application/json"
}
