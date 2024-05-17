package msgpack

import (
	libMsgpack "github.com/vmihailenco/msgpack/v5"
	zerocodec "github.com/zerogo-hub/zero-helper/codec"
)

type msgpackCodec struct{}

// New ..
func New() zerocodec.Codec {
	return &msgpackCodec{}
}

// Marshal 编码
func (*msgpackCodec) Marshal(in interface{}) ([]byte, error) {
	return libMsgpack.Marshal(in)
}

// Unmarshal 解码
func (*msgpackCodec) Unmarshal(in []byte, out interface{}) error {
	return libMsgpack.Unmarshal(in, out)
}

// Name 名称
func (*msgpackCodec) Name() string {
	return "msgpack"
}

// MimeType 媒体类型
func (*msgpackCodec) MimeType() string {
	return "application/msgpack"
}
