package cache

// GetBit 获取指定偏移量上的位
func (c *cache) GetBit(key string, offset int64) (int64, error) {
	return c.Int64(c.DO("GETBIT", key, offset))
}

// SetBit 设置或清除指定偏移量上的位(bit)
func (c *cache) SetBit(key string, offset, value int64) (int64, error) {
	return c.Int64(c.DO("SETBIT", key, offset, value))
}
