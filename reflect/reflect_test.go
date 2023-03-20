package reflect_test

import (
	"testing"

	zeroreflect "github.com/zerogo-hub/zero-helper/reflect"
)

type HelloStruct struct {
}

func TestStruct(t *testing.T) {
	r := zeroreflect.GetStructName(HelloStruct{})
	if r != "HelloStruct" {
		t.Error("test GetStructName failed")
	}

	r = zeroreflect.GetStructName(&HelloStruct{})
	if r != "HelloStruct" {
		t.Error("test GetStructName failed 2")
	}
}

func TestParseFuncName(t *testing.T) {
	if r := zeroreflect.ParseFuncName("hello.test"); r != "test" {
		t.Error("test TestParseFuncName failed")
	}

	if r := zeroreflect.ParseFuncName("test"); r != "test" {
		t.Error("test TestParseFuncName failed 2")
	}

	if r := zeroreflect.ParseFuncName("test."); r != "" {
		t.Error("test TestParseFuncName failed 3")
	}
}

func TestFunction(t *testing.T) {
	r := zeroreflect.GetFuncName(TestFunction)
	if r != "TestFunction" {
		t.Error("test GetStructName failed")
	}
}
