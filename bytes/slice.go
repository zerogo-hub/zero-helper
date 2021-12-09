package bytes

import (
	"strconv"
)

// SliceStringToInt8 []string 转 []int8
func SliceStringToInt8(from []string) ([]int8, error) {
	if len(from) == 0 {
		return []int8{}, nil
	}

	to := make([]int8, 0, len(from))

	for _, s := range from {
		i64, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			return nil, err
		}

		to = append(to, int8(i64))
	}

	return to, nil
}

// SliceStringToUint8 []string 转 []uint8
func SliceStringToUint8(from []string) ([]uint8, error) {
	if len(from) == 0 {
		return []uint8{}, nil
	}

	to := make([]uint8, 0, len(from))

	for _, s := range from {
		i64, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			return nil, err
		}

		to = append(to, uint8(i64))
	}

	return to, nil
}

// SliceStringToInt16 []string 转 []int6
func SliceStringToInt16(from []string) ([]int16, error) {
	if len(from) == 0 {
		return []int16{}, nil
	}

	to := make([]int16, 0, len(from))

	for _, s := range from {
		i64, err := strconv.ParseInt(s, 10, 16)
		if err != nil {
			return nil, err
		}

		to = append(to, int16(i64))
	}

	return to, nil
}

// SliceStringToUint16 []string 转 []uint16
func SliceStringToUint16(from []string) ([]uint16, error) {
	if len(from) == 0 {
		return []uint16{}, nil
	}

	to := make([]uint16, 0, len(from))

	for _, s := range from {
		i64, err := strconv.ParseUint(s, 10, 16)
		if err != nil {
			return nil, err
		}

		to = append(to, uint16(i64))
	}

	return to, nil
}

// SliceStringToInt32 []string 转 []int32
func SliceStringToInt32(from []string) ([]int32, error) {
	if len(from) == 0 {
		return []int32{}, nil
	}

	to := make([]int32, 0, len(from))

	for _, s := range from {
		i64, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return nil, err
		}

		to = append(to, int32(i64))
	}

	return to, nil
}

// SliceStringToUint32 []string 转 []uint32
func SliceStringToUint32(from []string) ([]uint32, error) {
	if len(from) == 0 {
		return []uint32{}, nil
	}

	to := make([]uint32, 0, len(from))

	for _, s := range from {
		i64, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return nil, err
		}

		to = append(to, uint32(i64))
	}

	return to, nil
}

// SliceStringToInt []string 转 []int
func SliceStringToInt(from []string) ([]int, error) {
	if len(from) == 0 {
		return []int{}, nil
	}

	to := make([]int, 0, len(from))

	for _, s := range from {
		i64, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return nil, err
		}

		to = append(to, int(i64))
	}

	return to, nil
}

// SliceStringToUint []string 转 []uint
func SliceStringToUint(from []string) ([]uint, error) {
	if len(from) == 0 {
		return []uint{}, nil
	}

	to := make([]uint, 0, len(from))

	for _, s := range from {
		i64, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return nil, err
		}

		to = append(to, uint(i64))
	}

	return to, nil
}

// SliceStringToInt64 []string 转 []int64
func SliceStringToInt64(from []string) ([]int64, error) {
	if len(from) == 0 {
		return []int64{}, nil
	}

	to := make([]int64, 0, len(from))

	for _, s := range from {
		i64, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}

		to = append(to, i64)
	}

	return to, nil
}

// SliceStringToUint64 []string 转 []uint64
func SliceStringToUint64(from []string) ([]uint64, error) {
	if len(from) == 0 {
		return []uint64{}, nil
	}

	to := make([]uint64, 0, len(from))

	for _, s := range from {
		i64, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, err
		}

		to = append(to, uint64(i64))
	}

	return to, nil
}

// SliceInt8ToString []int8 转 []string
func SliceInt8ToString(from []int8) []string {
	if len(from) == 0 {
		return nil
	}

	to := make([]string, 0, len(from))

	for _, i8 := range from {
		s := strconv.FormatInt(int64(i8), 10)
		to = append(to, s)
	}

	return to
}

// SliceUint8ToString []uint8 转 []string
func SliceUint8ToString(from []uint8) []string {
	if len(from) == 0 {
		return nil
	}

	to := make([]string, 0, len(from))

	for _, i8 := range from {
		s := strconv.FormatUint(uint64(i8), 10)
		to = append(to, s)
	}

	return to
}

// SliceInt16ToString []int16 转 []string
func SliceInt16ToString(from []int16) []string {
	if len(from) == 0 {
		return nil
	}

	to := make([]string, 0, len(from))

	for _, i16 := range from {
		s := strconv.FormatInt(int64(i16), 10)
		to = append(to, s)
	}

	return to
}

// SliceUint16ToString []uint16 转 []string
func SliceUint16ToString(from []uint16) []string {
	if len(from) == 0 {
		return nil
	}

	to := make([]string, 0, len(from))

	for _, i16 := range from {
		s := strconv.FormatUint(uint64(i16), 10)
		to = append(to, s)
	}

	return to
}

// SliceInt32ToString []int32 转 []string
func SliceInt32ToString(from []int32) []string {
	if len(from) == 0 {
		return nil
	}

	to := make([]string, 0, len(from))

	for _, i32 := range from {
		s := strconv.FormatInt(int64(i32), 10)
		to = append(to, s)
	}

	return to
}

// SliceUint32ToString []uint32 转 []string
func SliceUint32ToString(from []uint32) []string {
	if len(from) == 0 {
		return nil
	}

	to := make([]string, 0, len(from))

	for _, i32 := range from {
		s := strconv.FormatUint(uint64(i32), 10)
		to = append(to, s)
	}

	return to
}

// SliceIntToString []int 转 []string
func SliceIntToString(from []int) []string {
	if len(from) == 0 {
		return nil
	}

	to := make([]string, 0, len(from))

	for _, i := range from {
		s := strconv.FormatInt(int64(i), 10)
		to = append(to, s)
	}

	return to
}

// SliceUintToString []uint 转 []string
func SliceUintToString(from []uint) []string {
	if len(from) == 0 {
		return nil
	}

	to := make([]string, 0, len(from))

	for _, i := range from {
		s := strconv.FormatUint(uint64(i), 10)
		to = append(to, s)
	}

	return to
}

// SliceInt64ToString []int64 转 []string
func SliceInt64ToString(from []int64) []string {
	if len(from) == 0 {
		return nil
	}

	to := make([]string, 0, len(from))

	for _, i64 := range from {
		s := strconv.FormatInt(i64, 10)
		to = append(to, s)
	}

	return to
}

// SliceUint64ToString []uint64 转 []string
func SliceUint64ToString(from []uint64) []string {
	if len(from) == 0 {
		return nil
	}

	to := make([]string, 0, len(from))

	for _, i64 := range from {
		s := strconv.FormatUint(i64, 10)
		to = append(to, s)
	}

	return to
}
