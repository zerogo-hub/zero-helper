package random_test

import (
	"testing"

	zerorandom "github.com/zerogo-hub/zero-helper/random"
)

func TestLower(t *testing.T) {
	r := zerorandom.Lower(10)
	for _, c := range r {
		if c < 'a' || c > 'z' {
			t.Errorf("test lower failed, c: %v", c)
		}
	}
}

func TestLowerWithNumber(t *testing.T) {
	r := zerorandom.LowerWithNumber(10)
	for _, c := range r {
		if c > 'A' && c < 'Z' {
			t.Errorf("test lower with number failed, c: %v", c)
		}
	}
}

func TestUpper(t *testing.T) {
	r := zerorandom.Upper(10)
	for _, c := range r {
		if c < 'A' || c > 'Z' {
			t.Errorf("test upper failed, c: %v", c)
		}
	}
}

func TestUpperWithNumber(t *testing.T) {
	r := zerorandom.UpperWithNumber(10)
	for _, c := range r {
		if c > 'a' && c < 'z' {
			t.Errorf("test upper with number failed, c: %v", c)
		}
	}
}

func TestRandom(t *testing.T) {
	size := 10
	r1 := zerorandom.String(size)
	r2 := zerorandom.String(size)

	t.Logf("r1: %s", r1)
	t.Logf("r2: %s", r2)

	if len(r1) != size {
		t.Errorf("R1 length: %d is not size: %d", len(r1), size)
	}

	if len(r2) != size {
		t.Errorf("R2 length: %d is not size: %d", len(r2), size)
	}
}

func TestRangeInt(t *testing.T) {
	min := int64(1)
	max := int64(10)
	for i := 0; i < 1000; i++ {
		result := zerorandom.Int(min, max)
		if result < min || result >= max {
			t.Errorf("TestRangeInt failed, result: %d", result)
		}
	}

	max = min + 1
	for i := 0; i < 1000; i++ {
		result := zerorandom.Int(min, max)
		if result < min || result >= max {
			t.Errorf("TestRangeInt failed, result: %d", result)
		}
	}

	result := zerorandom.Int(5, 3)
	if result != 3 {
		t.Errorf("TestRangeInt failed, result: %d", result)
	}
}

func TestRangeUint32(t *testing.T) {
	max := ^uint32(0)
	t.Log(max)
	for i := 0; i < 1000; i++ {
		result := zerorandom.Uint32()
		if result > max {
			t.Errorf("test uint failed, result: %d", result)
		}
	}
}

func TestBytes(t *testing.T) {
	length := 10
	result := zerorandom.Bytes(length)
	if len(result) != length {
		t.Errorf("Expected length of result to be %d, but got %d", length, len(result))
	}
}

func TestLowerBytes(t *testing.T) {
	length := 10
	result := zerorandom.LowerBytes(length)
	if len(result) != length {
		t.Errorf("Expected length of result to be %d, but got %d", length, len(result))
	}
}

func TestLowerWithNumberBytes(t *testing.T) {
	length := 10
	result := zerorandom.LowerWithNumberBytes(length)
	if len(result) != length {
		t.Errorf("Expected length of result to be %d, but got %d", length, len(result))
	}
}

func TestUpperBytes(t *testing.T) {
	length := 10
	result := zerorandom.UpperBytes(length)
	if len(result) != length {
		t.Errorf("Expected length of result to be %d, but got %d", length, len(result))
	}
}

func TestUpperWithNumberBytes(t *testing.T) {
	length := 10
	result := zerorandom.UpperWithNumberBytes(length)
	if len(result) != length {
		t.Errorf("Expected length of result to be %d, but got %d", length, len(result))
	}
}
