package cache

// Eval 执行脚本
// keyCount: 参数个数
// args: 参数
func (c *cache) Eval(script string, keyCount int, keysAndArgs ...interface{}) (interface{}, error) {
	args := make([]interface{}, len(keysAndArgs)+2)
	args[0] = script
	args[1] = keyCount
	copy(args[2:], keysAndArgs)
	return c.DO("EVAL", args...)
}
