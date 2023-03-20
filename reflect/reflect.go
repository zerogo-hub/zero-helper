package reflect

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// GetStructName 获取结构体名称
func GetStructName(v interface{}) string {
	return reflect.Indirect(reflect.ValueOf(v)).Type().Name()
}

// GetFuncName 获取函数名称
func GetFuncName(fn interface{}) string {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		panic(fmt.Sprintf("[fn = %v] is not func type.", fn))
	}

	// "github.com/zerogo-hub/zero-helper/reflect_test.TestFunction"
	fullName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	return ParseFuncName(fullName)
}

// ParseFuncName ..
func ParseFuncName(fullName string) string {
	if len(fullName) == 0 {
		return ""
	}

	begin := strings.LastIndex(fullName, ".")
	if begin == -1 {
		return fullName
	}

	end := len(fullName)

	return string(fullName[begin+1 : end])
}
