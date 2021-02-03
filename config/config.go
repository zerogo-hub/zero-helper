package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
	yaml "gopkg.in/yaml.v2"
)

// eg:
// import (
// 	"github.com/alex-my/ghelper/config"
// )
//
// c := config.NewConfig()
// err := c.FileJSON("./test.json")
// if err != nil {
// 	return
// }
//
// 使用配置的实例来获取配置
// framework, _ := c.Any("framework")
//
// 如果只有一个配置实例，推荐使用 config.C, config.D 等辅助函数
// version, _ := config.C("framework")
//
// 如果配置不存在，使用默认值
// addr := config.D("addr", "127.0.0.1:8080")

// Load 加载配置文件
type Load interface {
	// LoadJSON 从 bytes 数据中读取配置
	LoadJSON(bytes []byte) error

	// LoadTOML 从 bytes 数据中读取配置
	LoadTOML(bytes []byte) error

	// LoadYAML 从 bytes 数据中读取配置
	LoadYAML(bytes []byte) error

	// FileJSON 从 json 文件中读取配置
	FileJSON(path string) error

	// FileTOML 从 toml 文件中读取配置
	FileTOML(path string) error

	// FileYAML 从 yaml 文件中读取配置
	FileYAML(path string) error
}

// Get 获取配置数据
type Get interface {
	// Any 根据 key 获取对应数据
	Any(key string) (interface{}, error)

	// C 获取配置
	C(key string) (string, error)

	// CB 获取配置，结果转为 bool
	CB(key string) (bool, error)

	// CI 获取配置，结果转为 int
	CI(key string) (int, error)

	// CI32 获取配置，结果转为 CI32
	CI32(key string) (int32, error)

	// CI64 获取配置，结果转为 int64
	CI64(key string) (int64, error)

	// CF32 获取配置，结果转为 float32
	CF32(key string) (float32, error)

	// CF64 获取配置，结果转为 float64
	CF64(key string) (float64, error)

	// DB 获取配置，结果转为 bool
	DB(key string, def bool) bool

	// DI 获取配置，结果转为 int
	DI(key string, def int) int

	// DI32 获取配置，结果转为 int32
	DI32(key string, def int32) int32

	// DI64 获取配置，结果转为 int64
	DI64(key string, def int64) int64

	// DF32 获取配置，结果转为 float32
	DF32(key string, def float32) float32

	// DF64 获取配置，结果转为 float64
	DF64(key string, def float64) float64
}

// Config 配置文件
type Config interface {
	Load
	Get
}

type config struct {
	// init true: 初始化完毕；false 尚未初始化完毕
	init bool
	// data 存储数据
	data map[string]interface{}
}

// NewConfig 生成一个配置文件实例
func NewConfig() Config {
	data := make(map[string]interface{})
	c := &config{data: data}
	return c
}

// LoadJSON 从 bytes 数据中读取 JSON 配置
func (c *config) LoadJSON(bytes []byte) error {
	if bytes == nil {
		return errors.New("Bytes cannot be empty")
	}

	err := json.Unmarshal(bytes, &c.data)
	if err != nil {
		return err
	}

	c.init = true
	return nil
}

// LoadTOML 从 bytes 数据中读取 TOML 配置
func (c *config) LoadTOML(bytes []byte) error {
	if bytes == nil {
		return errors.New("Bytes cannot be empty")
	}

	// m := map[string]interface{}{}
	if _, err := toml.Decode(string(bytes), &c.data); err != nil {
		return err
	}

	c.init = true
	return nil
}

// LoadYAML 从 bytes 数据中读取 YAML 配置
func (c *config) LoadYAML(bytes []byte) error {
	if bytes == nil {
		return errors.New("Bytes cannot be empty")
	}

	err := yaml.Unmarshal(bytes, &c.data)
	if err != nil {
		return err
	}

	c.init = true
	return nil
}

func loadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// FileJSON 从 json 文件中读取配置
func (c *config) FileJSON(path string) error {
	bytes, err := loadFile(path)
	if err != nil {
		return err
	}
	return c.LoadJSON(bytes)
}

// FileTOML 从 toml 文件中读取配置
func (c *config) FileTOML(path string) error {
	bytes, err := loadFile(path)
	if err != nil {
		return err
	}
	return c.LoadTOML(bytes)
}

// FileYAML 从 yaml 文件中读取配置
func (c *config) FileYAML(path string) error {
	bytes, err := loadFile(path)
	if err != nil {
		return err
	}
	return c.LoadYAML(bytes)
}

// Any 根据 key 获取对应数据
func (c *config) Any(key string) (interface{}, error) {
	if value, exist := c.data[key]; exist {
		return value, nil
	}
	return nil, errors.New("Configuration does not exist")
}

// C 获取配置
func (c *config) C(key string) (string, error) {
	value, err := c.Any(key)
	if err != nil {
		return "", err
	}
	return value.(string), nil
}

// CB 获取配置，结果转为 bool
func (c *config) CB(key string) (bool, error) {
	value, err := c.C(key)
	if err != nil {
		return false, err
	}

	switch value {
	case "1", "t", "T", "true", "TRUE", "True":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False":
		return false, nil
	default:
		errMsg := "\"1\", \"t\", \"T\", \"true\", \"TRUE\", \"True\" turned true, \"0\", \"f\", \"F\", \"false\", \"FALSE\", \"False\" turned false"
		return false, errors.New(errMsg)
	}
}

// CI 获取配置，结果转为 int
func (c *config) CI(key string) (int, error) {
	value, err := c.C(key)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(value)
}

// CI32 获取配置，结果转为 CI32
func (c *config) CI32(key string) (int32, error) {
	value, err := c.C(key)
	if err != nil {
		return 0, err
	}

	i32, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(i32), nil
}

// CI64 获取配置，结果转为 int64
func (c *config) CI64(key string) (int64, error) {
	value, err := c.C(key)
	if err != nil {
		return 0, err
	}

	i64, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return i64, nil
}

// CF32 获取配置，结果转为 float32
func (c *config) CF32(key string) (float32, error) {
	value, err := c.C(key)
	if err != nil {
		return 0, err
	}

	f32, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return 0, err
	}
	return float32(f32), nil
}

// CF64 获取配置，结果转为 float64
func (c *config) CF64(key string) (float64, error) {
	value, err := c.C(key)
	if err != nil {
		return 0, err
	}

	f64, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	return f64, nil
}

// D 获取配置
func (c *config) D(key, def string) string {
	value, err := c.C(key)
	if err == nil {
		return value
	}
	return def
}

// DB 获取配置，结果转为 bool
func (c *config) DB(key string, def bool) bool {
	value, err := c.CB(key)
	if err == nil {
		return value
	}
	return def
}

// DI 获取配置，结果转为 int
func (c *config) DI(key string, def int) int {
	value, err := c.CI(key)
	if err == nil {
		return value
	}
	return def
}

// DI32 获取配置，结果转为 int32
func (c *config) DI32(key string, def int32) int32 {
	value, err := c.CI32(key)
	if err == nil {
		return value
	}
	return def
}

// DI64 获取配置，结果转为 int64
func (c *config) DI64(key string, def int64) int64 {
	value, err := c.CI64(key)
	if err == nil {
		return value
	}
	return def
}

// DF32 获取配置，结果转为 float32
func (c *config) DF32(key string, def float32) float32 {
	value, err := c.CF32(key)
	if err == nil {
		return value
	}
	return def
}

// DF64 获取配置，结果转为 float64
func (c *config) DF64(key string, def float64) float64 {
	value, err := c.CF64(key)
	if err == nil {
		return value
	}
	return def
}
