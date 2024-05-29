package entity

// Results 查询结果
type Result struct {
	Vals [][]byte
	Errs []error
}

func (r *Result) Len() int {
	return len(r.Vals)
}

// Index 根据索引获取结果
func (r *Result) Index(idx int) ([]byte, error) {
	if idx < 0 || idx >= r.Len() {
		return nil, ErrResultIndexInvalid
	}

	if r.Errs[idx] != nil {
		return nil, r.Errs[idx]
	}

	return r.Vals[idx], nil
}
