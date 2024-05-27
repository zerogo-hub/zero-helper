package crypto_test

import (
	"bytes"
	"reflect"
	"testing"

	zerocrypto "github.com/zerogo-hub/zero-helper/crypto"
)

func TestMd5(t *testing.T) {
	input := "example"
	expected := "1a79a4d60de6718e8e5b326e338ae533"
	result := zerocrypto.Md5(input)
	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}

	input = ""
	expected = "d41d8cd98f00b204e9800998ecf8427e"
	result = zerocrypto.Md5(input)
	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestMd5Byte_NonASCII(t *testing.T) {
	input := []byte("你好，世界")                       // Chinese characters
	expected := "dbefd3ada018615b35588a01e216ae6e" // MD5 hash of the input
	result := zerocrypto.Md5Byte(input)

	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestMd5ByteToByte(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected []byte
	}{
		{[]byte("test"), []byte{9, 143, 107, 205, 70, 33, 211, 115, 202, 222, 78, 131, 38, 39, 180, 246}},
		{[]byte("hello"), []byte{93, 65, 64, 42, 188, 75, 42, 118, 185, 113, 157, 145, 16, 23, 197, 146}},
		{[]byte("world"), []byte{125, 121, 48, 55, 160, 118, 1, 134, 87, 75, 2, 130, 242, 244, 53, 231}},
	}

	for _, tc := range testCases {
		t.Run(string(tc.input), func(t *testing.T) {
			result := zerocrypto.Md5ByteToByte(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected: %v, but got: %v", tc.expected, result)
			}
		})
	}
}

func TestSha1(t *testing.T) {
	input := "test"
	expected := "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"
	result := zerocrypto.Sha1(input)
	if result != expected {
		t.Errorf("Sha1(%s) = %s; want %s", input, result, expected)
	}
}

func TestSha1_EmptyString(t *testing.T) {
	input := ""
	expected := "da39a3ee5e6b4b0d3255bfef95601890afd80709"
	result := zerocrypto.Sha1(input)
	if result != expected {
		t.Errorf("Sha1(%s) = %s; want %s", input, result, expected)
	}
}

func TestSha1_UnicodeString(t *testing.T) {
	input := "你好，世界"
	expected := "3becb03b015ed48050611c8d7afe4b88f70d5a20"
	result := zerocrypto.Sha1(input)
	if result != expected {
		t.Errorf("Sha1(%s) = %s; want %s", input, result, expected)
	}
}

func TestSha1_ByteInput(t *testing.T) {
	input := []byte("test")
	expected := "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"
	result := zerocrypto.Sha1Byte(input)
	if result != expected {
		t.Errorf("Sha1Byte(%s) = %s; want %s", input, result, expected)
	}
}

func TestSha1ByteToByte(t *testing.T) {
	input := []byte("test")
	expected := []byte{169, 74, 143, 229, 204, 177, 155, 166, 28, 76, 8, 115, 211, 145, 233, 135, 152, 47, 187, 211}
	result := zerocrypto.Sha1ByteToByte(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Sha1ByteToByte(%s) = %v; want %v", input, result, expected)
	}
}

func TestSha224(t *testing.T) {
	input := "test"
	expected := "90a3ed9e32b2aaf4c61c410eb925426119e1a9dc53d4286ade99a809"
	result := zerocrypto.Sha224(input)
	if result != expected {
		t.Errorf("Sha224(%s) = %s; want %s", input, result, expected)
	}
}

func TestSha224EmptyString(t *testing.T) {
	input := ""
	expected := "d14a028c2a3a2bc9476102bb288234c415a2b01f828ea62ac5b3e42f"
	result := zerocrypto.Sha224(input)
	if result != expected {
		t.Errorf("Sha224(%s) = %s; want %s", input, result, expected)
	}
}

func TestSha224LongString(t *testing.T) {
	input := "This is a long string for testing the Sha224 function with a long input to see if it produces the expected output"
	expected := "146a52b8b8596bc1eafe8ab47b35864cb110d202599f659b5a57cb87"
	result := zerocrypto.Sha224(input)
	if result != expected {
		t.Errorf("Sha224(%s) = %s; want %s", input, result, expected)
	}
}

func TestSha224Byte(t *testing.T) {
	input := []byte("example")
	expected := "312b3e578a63c0a34ed3f359263f01259e5cda07df73771d26928be5"
	result := zerocrypto.Sha224Byte(input)

	if result != expected {
		t.Errorf("Sha224Byte(%v) = %v; want %v", input, result, expected)
	}
}

func TestSha224Byte_EmptyInput(t *testing.T) {
	input := []byte("")
	expected := "d14a028c2a3a2bc9476102bb288234c415a2b01f828ea62ac5b3e42f"
	result := zerocrypto.Sha224Byte(input)

	if result != expected {
		t.Errorf("Sha224Byte(%v) = %v; want %v", input, result, expected)
	}
}

func TestSha224Byte_LongInput(t *testing.T) {
	input := []byte("This is a very long input string to test the Sha224Byte function with a long input")
	expected := "7fd3f4d26706ed33f8dc98ec1052095a5cc7be4d402f346d4469bde9"
	result := zerocrypto.Sha224Byte(input)

	if result != expected {
		t.Errorf("Sha224Byte(%v) = %v; want %v", input, result, expected)
	}
}

func TestSha224ByteToByte(t *testing.T) {
	input := []byte("example")
	expected := []byte{49, 43, 62, 87, 138, 99, 192, 163, 78, 211, 243, 89, 38, 63, 1, 37, 158, 92, 218, 7, 223, 115, 119, 29, 38, 146, 139, 229}
	result := zerocrypto.Sha224ByteToByte(input)
	if !bytes.Equal(result, expected) {
		t.Errorf("Sha224ByteToByte(%v) = %v; want %v", input, result, expected)
	}
}

func TestSha256(t *testing.T) {
	input := "test"
	expected := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	result := zerocrypto.Sha256(input)
	if result != expected {
		t.Errorf("Sha256(%s) = %s; want %s", input, result, expected)
	}
}

func TestSha256_EmptyString(t *testing.T) {
	input := ""
	expected := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	result := zerocrypto.Sha256(input)
	if result != expected {
		t.Errorf("Sha256(%s) = %s; want %s", input, result, expected)
	}
}

func TestSha256_UnicodeString(t *testing.T) {
	input := "你好，世界"
	expected := "46932f1e6ea5216e77f58b1908d72ec9322ed129318c6d4bd4450b5eaab9d7e7"
	result := zerocrypto.Sha256(input)
	if result != expected {
		t.Errorf("Sha256(%s) = %s; want %s", input, result, expected)
	}
}

func TestSha256_ByteInput(t *testing.T) {
	input := []byte("test")
	expected := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	result := zerocrypto.Sha256Byte(input)
	if result != expected {
		t.Errorf("Sha256Byte(%s) = %s; want %s", input, result, expected)
	}
}

func TestSha256ByteToByte(t *testing.T) {
	input := []byte("test")
	expected := []byte{159, 134, 208, 129, 136, 76, 125, 101, 154, 47, 234, 160, 197, 90, 208, 21, 163, 191, 79, 27, 43, 11, 130, 44, 209, 93, 108, 21, 176, 240, 10, 8}
	result := zerocrypto.Sha256ByteToByte(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Sha256ByteToByte(%s) = %v; want %v", input, result, expected)
	}
}

func TestSha384(t *testing.T) {
	result := zerocrypto.Sha384("test")
	expected := "768412320f7b0aa5812fce428dc4706b3cae50e02a64caa16a782249bfe8efc4b7ef1ccb126255d196047dfedf17a0a9"
	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestSha384EmptyString(t *testing.T) {
	result := zerocrypto.Sha384("")
	expected := "38b060a751ac96384cd9327eb1b1e36a21fdb71114be07434c0cc7bf63f6e1da274edebfe76f65fbd51ad2f14898b95b"
	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestSha384LongString(t *testing.T) {
	longStr := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
	result := zerocrypto.Sha384(longStr)
	expected := "8ac70afc3c5cb8daf3c332fbe88b7b9c2a5104217693b065a42db8de1bf40aca0bdaad7dbc6c405e97cf0809d1afcf52"
	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestSha384Byte_NonSliceInput(t *testing.T) {
	input := []byte("123456")
	want := "0a989ebc4a77b56a6e2bb7b19d995d185ce44090c13e2984b7ecc6d446d4b61ea9991b76a4c2f04b1b4d244841449454"
	if result := zerocrypto.Sha384Byte(input); result != want {
		t.Errorf("Sha384Byte(%v) = %v; want %v", input, result, want)
	}
}

func TestSha384ByteToByte(t *testing.T) {
	input := []byte("example")
	expected := []byte{254, 238, 191, 136, 79, 109, 171, 230, 236, 168, 214, 142, 55, 61, 107, 228, 136, 205, 170, 94, 183, 100, 232, 149, 41, 3, 54, 255, 233, 255, 150, 150, 134, 242, 169, 211, 98, 233, 168, 187, 221, 246, 231, 178, 225, 69, 95, 45}

	result := zerocrypto.Sha384ByteToByte(input)

	if !bytes.Equal(result, expected) {
		t.Errorf("Sha384ByteToByte(%v) = %v; want %v", input, result, expected)
	}
}

func TestSha512(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
		},
		{
			name:     "Short string",
			input:    "hello",
			expected: "9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043",
		},
		{
			name:     "Long string",
			input:    "a very long string that exceeds typical lengths but should still be hashed correctly",
			expected: "20cc1396a83144490e486e08e7ed61e48afa3594e2961ebcc24884a4f9152786721b1dcb327c792b0c62d7194c40a1dca938e2a69d79a48034a5003b6e4e1db1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := zerocrypto.Sha512(tt.input)
			if got != tt.expected {
				t.Errorf("Sha512() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestURLEncode_NonASCII(t *testing.T) {
	input := "www.keylala.cn?name=alex&age=18&say=你好"
	expected := "www.keylala.cn%3Fname%3Dalex%26age%3D18%26say%3D%E4%BD%A0%E5%A5%BD"
	result := zerocrypto.URLEncode(input)
	if result != expected {
		t.Errorf("URLEncode did not handle non-ASCII characters correctly, got: %s, want: %s", result, expected)
	}
}

func TestURLDecode_NormalCase(t *testing.T) {
	input := "www.keylala.cn%3Fname%3Dalex%26age%3D18%26say%3D%E4%BD%A0%E5%A5%BD"
	expected := "www.keylala.cn?name=alex&age=18&say=你好"
	result := zerocrypto.URLDecode(input)
	if result != expected {
		t.Errorf("URLDecode(%s) = %s; want %s", input, result, expected)
	}
}

func TestURLDecode_EmptyInput(t *testing.T) {
	input := ""
	expected := ""
	result := zerocrypto.URLDecode(input)
	if result != expected {
		t.Errorf("URLDecode(%s) = %s; want %s", input, result, expected)
	}
}

func TestURLDecode_ErrorCase(t *testing.T) {
	input := "www.keylala.cn%3Fname%3Dalex%26age%3D18%26say%3D%E4%BD%A0%E5%A5%BD%"
	expected := ""
	result := zerocrypto.URLDecode(input)
	if result != expected {
		t.Errorf("URLDecode(%s) = %s; want %s", input, result, expected)
	}
}

func TestBase64Encode(t *testing.T) {
	input := "hello world"
	expected := "aGVsbG8gd29ybGQ="
	result := zerocrypto.Base64Encode(input)
	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestBase64Encode_EmptyString(t *testing.T) {
	input := ""
	expected := ""
	result := zerocrypto.Base64Encode(input)
	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestBase64Encode_UnicodeString(t *testing.T) {
	input := "你好，世界"
	expected := "5L2g5aW977yM5LiW55WM"
	result := zerocrypto.Base64Encode(input)
	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestBase64Encode_WithSpecialCharacters(t *testing.T) {
	input := "a@b#c$d%e^f&g*h"
	expected := "YUBiI2MkZCVlXmYmZypo"
	result := zerocrypto.Base64Encode(input)
	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestBase64Encode_WithPadding(t *testing.T) {
	input := "12345"
	expected := "MTIzNDU="
	result := zerocrypto.Base64Encode(input)
	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestBase64EncodeByte(t *testing.T) {
	input := []byte("hello world")
	expected := "aGVsbG8gd29ybGQ="
	result := zerocrypto.Base64EncodeByte(input)
	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestBase64EncodeByte_EmptyInput(t *testing.T) {
	input := []byte("")
	expected := ""
	result := zerocrypto.Base64EncodeByte(input)
	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestBase64EncodeByte_SpecialCharacters(t *testing.T) {
	input := []byte("你好，世界！Hello, World!")
	expected := "5L2g5aW977yM5LiW55WM77yBSGVsbG8sIFdvcmxkIQ=="
	result := zerocrypto.Base64EncodeByte(input)
	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestBase64Decode_ValidInput(t *testing.T) {
	input := "SGVsbG8gV29ybGQh" // Base64 encoded "Hello World!"
	expected := []byte("Hello World!")

	result, err := zerocrypto.Base64Decode(input)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %v, but got: %v", expected, result)
	}
}

func TestBase64Decode_InvalidInput(t *testing.T) {
	input := "InvalidBase64String"

	_, err := zerocrypto.Base64Decode(input)

	if err == nil {
		t.Error("Expected an error, but got nil")
	}
}

func TestAesCBCEncodeBase64_Success(t *testing.T) {
	key := []byte("1234567890123456")
	iv := []byte("1234567890123456")
	plaintext := []byte("hello world")

	expected := "bAx40eFUVf/hIxbaV8/GaQ=="
	result, err := zerocrypto.AesCBCEncodeBase64(key, iv, plaintext)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expected {
		t.Errorf("Expected: %s, but got: %s", expected, result)
	}
}

func TestAesCBCEncodeBase64_Error(t *testing.T) {
	key := []byte("invalidkey1234567")
	iv := []byte("1234567890123456")
	plaintext := []byte("hello world")

	// Injecting an incorrect key to force an error

	_, err := zerocrypto.AesCBCEncodeBase64(key, iv, plaintext)
	if err == nil {
		t.Error("Expected an error but got nil")
	}
}

func TestAesCBCEncodeHex_Success(t *testing.T) {
	key := []byte("1234567890123456")
	iv := []byte("1234567890123456")
	plaintext := []byte("abcdefg")

	expectedHex := "ae5d9a1e7e4260832cba80647b1e788d"

	result, err := zerocrypto.AesCBCEncodeHex(key, iv, plaintext)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result != expectedHex {
		t.Errorf("Expected: %s, but got: %s", expectedHex, result)
	}
}

func TestAesCBCEncodeHex_InvalidKey(t *testing.T) {
	key := []byte("shortkey") // Invalid key length
	iv := []byte("1234567890123456")
	plaintext := []byte("abcdefg")

	_, err := zerocrypto.AesCBCEncodeHex(key, iv, plaintext)
	if err == nil {
		t.Error("Expected error for invalid key length, but got nil")
	}
}

func TestAesCBCEncodeHex_InvalidIV(t *testing.T) {
	key := []byte("1234567890123456")
	iv := []byte("shortiv") // Invalid IV length
	plaintext := []byte("abcdefg")

	_, err := zerocrypto.AesCBCEncodeHex(key, iv, plaintext)
	if err == nil {
		t.Error("Expected error for invalid IV length, but got nil")
	}
}

func TestAesCBCDecodeBase64_DecodeFailure(t *testing.T) {
	key := []byte("examplekey123456")
	iv := []byte("exampleiv12345678")
	invalidBase64 := "invalidBase64String"

	_, err := zerocrypto.AesCBCDecodeBase64(key, iv, invalidBase64)

	if err == nil {
		t.Error("Expected an error for decoding invalid base64, but got nil")
	}
}
