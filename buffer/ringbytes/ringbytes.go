package ringbytes

import (
	"errors"
	"fmt"
)

var (
	// ErrTooManyToWrite 内容太长，无法写入
	ErrTooManyToWrite = errors.New("too many data to write")
	// ErrInvalidBuffer 用于读取的 buffer 无效
	ErrInvalidBuffer = errors.New("invalid buffer")
	// ErrNotEnough 数据不足以读取长度 n 的内容
	ErrNotEnough = errors.New("not enough data")
	// ErrInvalidSkipSize 无效的长度
	ErrInvalidSkipSize = errors.New("invalid skip size")
	// ErrIsFull 缓冲区已满
	ErrIsFull = errors.New("buffer is full")
	// ErrIsEmpty 缓冲区已满
	ErrIsEmpty = errors.New("buffer is empty")
	// ErrInvalidLength 无效长度
	ErrInvalidLength = errors.New("invalid length")
)

// RingBytes 环形缓冲区，无锁版本，线程不安全
type RingBytes struct {
	write int
	read  int
	size  int

	buffer     []byte
	copyBuffer []byte
	isEmpty    bool
}

// NewUnlock 创建环形结构，无锁版本，线程不安全
func New(size int) *RingBytes {
	if size <= 0 {
		size = 4
	}

	return &RingBytes{
		write:      0,
		read:       0,
		size:       size,
		buffer:     make([]byte, size),
		copyBuffer: make([]byte, size),
		isEmpty:    true,
	}
}

// Len 存储的有效数据长度
func (r *RingBytes) Len() int {
	// read 到 write 中间存储着有效的 buffer

	if r.IsFull() {
		return r.size
	}

	// |___read___>___write____>|
	if r.write >= r.read {
		return r.write - r.read
	}

	// |___write______read____|
	return r.size - (r.read - r.write)
}

// Cap 容量
func (r *RingBytes) Cap() int {
	return r.size
}

// Free 剩余空间
func (r *RingBytes) Free() int {
	return r.size - r.Len()
}

// Reset 重置
func (r *RingBytes) Reset() {
	r.write = 0
	r.read = 0
	r.isEmpty = true
}

func (r *RingBytes) String() string {
	return fmt.Sprintf("read: %d, write: %d, size: %d, len: %d, free: %d, isFull: %t, buffer: %v",
		r.read, r.write, r.size, r.Len(), r.Free(), r.IsFull(), r.buffer,
	)
}

// IsEmpty 缓冲区是否为空，没有任何数据可读
func (r *RingBytes) IsEmpty() bool {
	return r.isEmpty
}

// IsFull 缓冲区是否已满，不可以继续写入任何数据
func (r *RingBytes) IsFull() bool {
	return !r.isEmpty && r.write == r.read
}

// Write 写入数据
//
// 只有全部写入，或者都不写入两种情况，不存在只写入一部分的情形
func (r *RingBytes) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	if r.IsFull() {
		return 0, ErrIsFull
	}

	free := r.Free()
	lenp := len(p)
	if lenp > free {
		return 0, ErrTooManyToWrite
	}

	// 空闲部分足以写入数据
	// 开始写入数据

	// start|___read______write____|end
	if r.write >= r.read {
		lenWriteToEnd := r.size - r.write

		// 如果 write 到 end 能写入全部数据
		if lenWriteToEnd >= lenp {
			copy(r.buffer[r.write:], p)
			r.write += lenp
		} else {
			// 数据分两部分写入
			// 第一部分写入 write 到 end 的空间
			copy(r.buffer[r.write:], p[:lenWriteToEnd])
			// 第二部分写到 start 到 read 之间
			copy(r.buffer[:], p[lenWriteToEnd:])
			left := lenp - lenWriteToEnd

			// start|__write_read__________|end
			r.write = left
		}
	} else {
		// start|___write______read____|end
		copy(r.buffer[r.write:], p)
		r.write += lenp
	}

	r.isEmpty = false

	if r.write == r.size {
		r.write = 0
	}

	return lenp, nil
}

// WriteN 从 p 中取出长度 n 的数据写入
// 只有全部写入，或者都不写入两种情况，不存在只写入一部分的情形
// 当剩余空间不足时，写入失败
func (r *RingBytes) WriteN(p []byte, n int) error {
	if len(p) < n {
		return ErrInvalidLength
	}

	_, err := r.Write(p[:n])
	return err
}

// Read 获取缓冲中的数据
func (r *RingBytes) Read(n int) ([]byte, error) {
	if n <= 0 {
		return nil, ErrInvalidLength
	}

	// 没有内容
	if r.IsEmpty() {
		return nil, ErrIsEmpty
	}

	if n > r.Len() {
		return nil, ErrInvalidLength
	}

	defer r.tryReset()

	// start|___read______write____|end
	if r.write > r.read {
		r.read += n
		return r.buffer[r.read-n : r.read], nil
	}

	// r.write <= r.read
	// start|___write______read____|end

	if r.read+n <= r.size {
		r.read += n
		return r.buffer[r.read-n : r.read], nil
	}

	// 拷贝两部分
	copy(r.copyBuffer, r.buffer[r.read:])
	left := n - (r.size - r.read)
	copy(r.copyBuffer[r.size-r.read:], r.buffer[:left])

	r.read = left
	return r.copyBuffer[:n], nil
}

// Peek 读取 n 个长度的内容，但不会产生任何影响
// 不会改变 r.read 和 r.write
func (r *RingBytes) Peek(n int) ([]byte, error) {
	write, read := r.write, r.read

	defer func() {
		r.write, r.read = write, read
	}()

	p, err := r.Read(n)

	return p, err
}

// Skip 跳过 n 个字段，会改变 r.read
func (r *RingBytes) Skip(n int) error {
	if n <= 0 {
		return nil
	}

	if r.IsEmpty() {
		return ErrIsEmpty
	}

	// start|___read______write____|end
	if r.write > r.read {
		if r.read+n > r.write {
			return ErrInvalidSkipSize
		}

		r.read += n
	} else {
		// r.write <= r.read
		// start|___write______read____|end
		if r.read+n <= r.size {
			r.read += n
		} else {
			left := n - (r.size - r.read)
			if left > r.write {
				return ErrInvalidSkipSize
			}

			r.read = left
		}
	}

	r.tryReset()

	return nil
}

func (r *RingBytes) tryReset() {
	if r.read == r.write && r.write > 0 {
		r.Reset()
	}
}
