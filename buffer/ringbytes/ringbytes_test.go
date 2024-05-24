package ringbytes

import (
	"reflect"
	"testing"
)

func TestNewRingBytesWithPositiveSize(t *testing.T) {
	size := 10
	rb := New(size)

	if rb.Cap() != size {
		t.Errorf("Expected size to be %d, but got %d", size, rb.size)
	}
	if rb.write != 0 {
		t.Error("Expected write to be 0")
	}
	if rb.read != 0 {
		t.Error("Expected read to be 0")
	}
	if len(rb.buffer) != size {
		t.Errorf("Expected buffer length to be %d, but got %d", size, len(rb.buffer))
	}
	if len(rb.copyBuffer) != size {
		t.Errorf("Expected copyBuffer length to be %d, but got %d", size, len(rb.copyBuffer))
	}
	if !rb.isEmpty {
		t.Error("Expected isEmpty to be true")
	}

	_ = rb.String()
}

func TestNewRingBytesWithZeroSize(t *testing.T) {
	size := 0
	rb := New(size)

	if rb.size != 4 {
		t.Errorf("Expected size to be 4, but got %d", rb.size)
	}
	if rb.write != 0 {
		t.Error("Expected write to be 0")
	}
	if rb.read != 0 {
		t.Error("Expected read to be 0")
	}
	if len(rb.buffer) != 4 {
		t.Errorf("Expected buffer length to be 4, but got %d", len(rb.buffer))
	}
	if len(rb.copyBuffer) != 4 {
		t.Errorf("Expected copyBuffer length to be 4, but got %d", len(rb.copyBuffer))
	}
	if !rb.isEmpty {
		t.Error("Expected isEmpty to be true")
	}
}

func TestRingBytesUnlock_Write(t *testing.T) {
	r := New(5)

	// Test case 1: Write data within the available space
	data := []byte{1, 2, 3}
	n, err := r.Write(data)
	if n != 3 || err != nil {
		t.Errorf("Expected to write 3 bytes without error, got %d bytes and error: %v", n, err)
	}
	if r.write != 3 || r.read != 0 || r.isEmpty != false {
		t.Errorf("Write did not update the write and read pointers or set isEmpty to false")
	}

	r.Skip(n)

	// Test case 2: Write data that spans the buffer boundary
	data = []byte{4, 5, 6, 7}
	n, err = r.Write(data)
	if n != 4 || err != nil {
		t.Errorf("Expected to write 4 bytes without error, got %d bytes and error: %v", n, err)
	}
	if r.Len() != 4 || r.Free() != 1 {
		t.Error("Unexcepted len and free")
	}
}

func TestRingBytes_Len_BufferFull(t *testing.T) {
	r := New(5)
	r.Write([]byte{1, 2, 3, 4, 5})
	expected := 5
	if result := r.Len(); result != expected {
		t.Errorf("Expected length: %d, but got: %d", expected, result)
	}
}

func TestRingBytes_Len_WriteGreaterThanRead(t *testing.T) {
	r := New(5)
	r.write = 4
	r.read = 2
	expected := 2
	if result := r.Len(); result != expected {
		t.Errorf("Expected length: %d, but got: %d", expected, result)
	}
}

func TestRingBytes_Len_ReadGreaterThanWrite(t *testing.T) {
	r := New(5)
	r.write = 2
	r.read = 4
	expected := 3
	if result := r.Len(); result != expected {
		t.Errorf("Expected length: %d, but got: %d", expected, result)
	}
}
func TestRingBytes_WriteN(t *testing.T) {
	r := New(10)
	data := []byte{1, 2, 3, 4, 5}
	n := 3
	err := r.WriteN(data, n)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := []byte{1, 2, 3}
	readData, _ := r.Read(n)
	if !reflect.DeepEqual(readData, expected) {
		t.Errorf("Expected %v, got %v", expected, readData)
	}
}

func TestRingBytes_WriteN_InvalidLength(t *testing.T) {
	r := New(10)
	data := []byte{1, 2, 3, 4, 5}
	n := 6
	err := r.WriteN(data, n)
	if err != ErrInvalidLength {
		t.Errorf("Expected ErrInvalidLength, got %v", err)
	}
}

func TestRingBytes_WriteN_TooManyToWrite(t *testing.T) {
	r := New(3)
	data := []byte{1, 2, 3, 4, 5}
	n := 5
	err := r.WriteN(data, n)
	if err != ErrTooManyToWrite {
		t.Errorf("Expected ErrTooManyToWrite, got %v", err)
	}
}

func TestRingBytes_WriteN_WriteError(t *testing.T) {
	r := New(3)
	data := []byte{1, 2, 3}
	n := 3
	r.Write(data)
	err := r.WriteN(data, n)
	if err != ErrIsFull {
		t.Errorf("Expected ErrIsFull, got %v", err)
	}
}

func TestRingBytes_Read_EmptyBuffer(t *testing.T) {
	r := New(5)
	data, err := r.Read(3)
	if err != ErrIsEmpty {
		t.Errorf("Expected ErrIsEmpty, got %v", err)
	}
	if data != nil {
		t.Errorf("Expected nil data, got %v", data)
	}
}

func TestRingBytes_Read_InvalidLength(t *testing.T) {
	r := New(5)
	r.Write([]byte{1, 2, 3, 4})
	data, err := r.Read(6)
	if err != ErrInvalidLength {
		t.Errorf("Expected ErrInvalidLength, got %v", err)
	}
	if data != nil {
		t.Errorf("Expected nil data, got %v", data)
	}
}

func TestRingBytes_Read_SingleRead(t *testing.T) {
	r := New(5)
	r.Write([]byte{1, 2, 3, 4})
	data, err := r.Read(3)
	expected := []byte{1, 2, 3}
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("Expected %v, got %v", expected, data)
	}
}

func TestRingBytes_Read_WrapAroundRead(t *testing.T) {
	r := New(5)
	r.Write([]byte{1, 2, 3, 4})

	if r.Len() != 4 {
		t.Errorf("Expected len to be 4, got %d", r.Len())
	}

	data, err := r.Read(4)
	expected := []byte{1, 2, 3, 4}
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("Expected %v, got %v", expected, data)
	}

	if n, err := r.Write([]byte{11, 12}); n != 2 || err != nil {
		t.Errorf("Expected no error, n: %d, got %v", n, err)
	}

	if r.Len() != 2 {
		t.Errorf("Expected len to be 2, got %d", r.Len())
	}

	// 尾部 11，头部 12
	data, err = r.Read(2)
	expected = []byte{11, 12}
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("Expected %v, got %v", expected, data)
	}
}

func TestRingBytes_Peek(t *testing.T) {
	r := New(10)
	data := []byte{1, 2, 3, 4, 5}
	r.Write(data)

	t.Run("Peek with valid n", func(t *testing.T) {
		p, err := r.Peek(3)
		if err != nil {
			t.Errorf("Expected no error, but got %v", err)
		}
		expected := []byte{1, 2, 3}
		if !reflect.DeepEqual(p, expected) {
			t.Errorf("Expected %v, but got %v", expected, p)
		}
	})

	t.Run("Peek with n greater than available data", func(t *testing.T) {
		p, err := r.Peek(6)
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		if p != nil {
			t.Errorf("Expected nil, but got %v", p)
		}
	})
}

func TestRingBytes_Skip_NegativeValue(t *testing.T) {
	r := New(10)
	err := r.Skip(-1)
	if err != nil {
		t.Errorf("Expected no error for negative value, but got: %v", err)
	}
}

func TestSkip_LargerThanRemainingData(t *testing.T) {
	tests := []struct {
		name      string
		initial   []byte
		skipSize  int
		expectErr error
	}{
		{
			name:      "Skip larger than remaining data when write > read",
			initial:   []byte{1, 2, 3, 4, 5},
			skipSize:  6,
			expectErr: ErrInvalidSkipSize,
		},
		{
			name:      "Skip exactly the remaining data when write > read",
			initial:   []byte{1, 2, 3, 4, 5},
			skipSize:  5,
			expectErr: nil,
		},
		{
			name:      "Skip exactly the remaining data when write <= read",
			initial:   []byte{1, 2, 3, 4, 5},
			skipSize:  5,
			expectErr: nil,
		},
		{
			name:      "Skip with empty buffer",
			initial:   []byte{},
			skipSize:  1,
			expectErr: ErrIsEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := New(len(tt.initial))
			rb.Write(tt.initial)
			err := rb.Skip(tt.skipSize)
			if err != tt.expectErr {
				t.Errorf("expected error %v, got %v", tt.expectErr, err)
			}
		})
	}
}
