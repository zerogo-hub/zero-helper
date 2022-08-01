package zlib

import (
	"bytes"
	czlib "compress/zlib"
	"io/ioutil"

	zerocompress "github.com/zerogo-hub/zero-helper/compress"
)

type zlib struct {
	level int
}

// NewZlib ..
func NewZlib(level ...int) zerocompress.Compress {
	l := czlib.DefaultCompression
	if len(level) > 0 {
		l = level[0]
	}

	return &zlib{
		level: l,
	}
}

// Compress 压缩
func (zlib *zlib) Compress(in []byte) ([]byte, error) {
	var buffer bytes.Buffer

	writer, err := czlib.NewWriterLevel(&buffer, zlib.level)
	if err != nil {
		return nil, err
	}

	if _, err = writer.Write(in); err != nil {
		return nil, err
	}

	if err := writer.Flush(); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// Uncompress 解压缩
func (zlib *zlib) Uncompress(in []byte) ([]byte, error) {
	reader, err := czlib.NewReader(bytes.NewReader(in))
	if err != nil {
		return nil, err
	}

	defer reader.Close()

	return ioutil.ReadAll(reader)
}

// Name 获取名称
func (zlib *zlib) Name() string {
	return "zlib"
}
