package entity

// Results 查询结果
type Result struct {
	IDIndexs map[uint64]int
	Vals     map[int][]byte
	Errs     map[int]error
}

func (r *Result) Len() int {
	return len(r.IDIndexs)
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

// Get 根据 ID 获取结果
func (r *Result) Get(id uint64) ([]byte, error) {
	idx, ok := r.IDIndexs[id]
	if !ok {
		return nil, ErrResultIdNotFound
	}

	return r.Index(idx)
}
