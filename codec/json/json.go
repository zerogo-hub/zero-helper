package json

import (
	zerocodec "github.com/zerogo-hub/zero-helper/codec"
	zerojson "github.com/zerogo-hub/zero-helper/json"
)

type jsonCodec struct{}

// New JSON
func New() zerocodec.Codec {
	return &jsonCodec{}
}

// Marshal 编码
func (*jsonCodec) Marshal(in interface{}) ([]byte, error) {
	return zerojson.Marshal(in)
}

// Unmarshal 解码
func (*jsonCodec) Unmarshal(in []byte, out interface{}) error {
	return zerojson.Unmarshal(in, out)
}

// Name 名称
func (*jsonCodec) Name() string {
	return "json"
}

// MimeType 媒体类型
func (*jsonCodec) MimeType() string {
	return "application/json"
}
