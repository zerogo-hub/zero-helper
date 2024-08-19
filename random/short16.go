package random

// New16Snowflake 支持 16 个节点，每个节点每毫秒生成 64
// node: [0,15]
//
// 生成的 ID:
// 2024年: 56765078912
// 2044年: 646299648000384
// 2054年: 969405235200384
func New16Snowflake(workerID int) (*Snowflake, error) {
	return NewSnowflakeBy(workerID, defaultSnowflakeOriginTime, 4, 6, nil, nil)
}
