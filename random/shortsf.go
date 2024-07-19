package random

import "errors"

// NewBit46Snowflake 缩短长度版本的 Snowflake，uint64 在一些脚本语言中默认会呈科学计数法显示甚至精度丢失
// 而平常并不需要这么长的长度
// 如 javascript 超过 2^53 9007199254740992 有可能会精度丢失。6145390195186705111 显示为 6145390195186705000
// 不能超过 16 位
//
// workID 取值 [0,1]
//
// 配置:
//
// 0(1 bit) - 毫秒时间戳(41 bit) - 节点 id(1 bit) - 序列号(3 bit)
//
// 毫秒时间戳(41 bit)：存储毫秒时间戳，取值范围 [0,1<<41)，目前存储的是当前毫秒时间戳与 originTime 的差值，可以在 69 年内保障唯一，可以设置 SetOriginTime 修改这 69 年的起始时间
//
// 节点 id(1 bit): 可以分布在 2 个节点上
//
// 序列号(3 bit)：每毫秒可以生成 15 个 UUID
func NewBit46Snowflake(workerID int) (*Snowflake, error) {
	if workerID < 0 || workerID > 1 {
		return nil, errors.New("invalid work id")
	}

	return NewSnowflakeBy(workerID, defaultSnowflakeOriginTime, 1, 3, nil, nil)
}
