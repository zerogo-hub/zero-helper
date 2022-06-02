package zerror

import "bytes"

// ZMultiError 多个错误
type ZMultiError struct {
	// Causes 保存多个原始错误
	Causes []error
}

func NewMultiError() *ZMultiError {
	return new(ZMultiError)
}

// Add 新增错误
func (e *ZMultiError) Add(errs ...error) {
	if e.Causes == nil {
		e.Causes = make([]error, 0, len(errs)*3)
	}

	e.Causes = append(e.Causes, errs...)
}

// String 输出错误信息
func (e *ZMultiError) Error() string {
	var buf bytes.Buffer

	for i, cause := range e.Causes {
		if i > 0 {
			buf.WriteByte('\n')
		}
		buf.WriteString(cause.Error())
	}

	return buf.String()
}
