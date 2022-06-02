package gogoprotobuf

import (
	"errors"

	"github.com/gogo/protobuf/proto"
	zerocodec "github.com/zerogo-hub/zero-helper/codec"
)

var (
	// ErrInvalidPBMessage 无效的 google protobuf 消息
	ErrInvalidPBMessage = errors.New("invalid pb message")
)

type protobufCodec struct{}

// NewProtobufCodec 创建默认的编码与解码器
func NewProtobufCodec() zerocodec.Codec {
	return &protobufCodec{}
}

// Marshal 编码
func (c *protobufCodec) Marshal(in interface{}) ([]byte, error) {
	m, ok := in.(proto.Message)
	if !ok {
		return nil, ErrInvalidPBMessage
	}

	return proto.Marshal(m)
}

// Unmarshal 解码
func (c *protobufCodec) Unmarshal(in []byte, out interface{}) error {
	m, ok := out.(proto.Message)
	if !ok {
		return ErrInvalidPBMessage
	}

	return proto.Unmarshal(in, m)
}

// Name 名称
func (*protobufCodec) Name() string {
	return "protobuf"
}

// MimeType 媒体类型
func (*protobufCodec) MimeType() string {
	return "application/binary"
}
