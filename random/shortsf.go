package random

import "errors"

// 缩短长度版本的 Snowflake，uint64 在一些脚本语言中会呈科学计数法显示
// 而平常并不需要这么长的长度
//
// shortsf: 46
// 0 - 毫秒时间戳(41 bit) - 序列号(4 bit) - 机器 id(1 bit)

// NewBit46Snowflake 生成46 bit 的 uint64
// workerId: 1-8
func NewBit46Snowflake(workerID int) (*Snowflake, error) {
	if workerID <= 0 || workerID > 8 {
		return nil, errors.New("invalid work id")
	}

	return NewSnowflakeBy(workerID, defaultSnowflakeOriginTime, 1, 4, nil, nil)
}
