package bytes_test

import (
	"testing"

	zerobytes "github.com/zerogo-hub/zero-helper/bytes"
)

func TestStringToBytes(t *testing.T) {
	s := "hello world"

	b1 := zerobytes.StringToBytes(s)
	b2 := []byte(s)

	if len(b1) != len(b2) {
		t.Error("TestStringToBytes error")
		return
	}

	for idx, b10 := range b1 {
		b20 := b2[idx]
		if b10 != b20 {
			t.Error("TestStringToBytes not the same")
			break
		}
	}
}

func TestSliceByteToString_InvalidCharacters(t *testing.T) {
	input := []byte{72, 101, 108, 108, 111, 255} // Contains an invalid character (255)
	expected := "Hello\xFF"                      // Expected string representation with the invalid character
	result := zerobytes.SliceByteToString(input)

	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestCharLower(t *testing.T) {
	if zerobytes.CharLower('A') != 'a' {
		t.Error("CharLower A error")
	}

	if zerobytes.CharLower('Z') != 'z' {
		t.Error("CharLower Z error")
	}

	if zerobytes.CharLower(1) != 1 {
		t.Error("CharLower 1 error")
	}
}

func TestCharUpper(t *testing.T) {
	if zerobytes.CharUpper('a') != 'A' {
		t.Error("CharLower a error")
	}

	if zerobytes.CharUpper('z') != 'Z' {
		t.Error("CharLower z error")
	}

	if zerobytes.CharUpper(1) != 1 {
		t.Error("CharLower 1 error")
	}
}

func TestFirstLower(t *testing.T) {
	if zerobytes.FirstLower("Hello world") != "hello world" {
		t.Error("FirstLower error")
	}

	if zerobytes.FirstLower("") != "" {
		t.Error("FirstLower empty error")
	}
}

func TestFirstUpper(t *testing.T) {
	if zerobytes.FirstUpper("hello world") != "Hello world" {
		t.Error("FirstUpper error")
	}

	if zerobytes.FirstUpper("") != "" {
		t.Error("FirstUpper empty error")
	}
}
