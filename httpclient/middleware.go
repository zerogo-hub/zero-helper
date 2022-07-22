package httpclient

import (
	"bytes"
	"net/http"
	"sort"
	"strconv"
	"sync"

	zerocrypto "github.com/zerogo-hub/zero-helper/crypto"
	zerorandom "github.com/zerogo-hub/zero-helper/random"
	zerotime "github.com/zerogo-hub/zero-helper/time"
)

// WithSign 自动签名中间件，自动添加 timestamp, nonce, sign 值
// secret: 签名使用，签名方式见 github.com/zerogo-hub/zero-api-middleware/sign
func WithSign(client *HTTPClient, secret string) {
	handler := func(client *HTTPClient, method, url string) error {
		timestamp := zerotime.Now()
		nonce := zerorandom.LowerWithNumber(32)

		if method == http.MethodGet || method == http.MethodDelete {
			// 参数添加到 params 中
			client.WithParams(map[string]string{
				"timestamp": strconv.FormatUint(uint64(timestamp), 10),
				"nonce":     nonce,
			})
			sign := calcParamsSign(secret, client.Params())
			client.WithParams(map[string]string{
				"sign": sign,
			})
		} else if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
			// 参数添加到 body 中
			client.WithBody(map[string]interface{}{
				"timestamp": strconv.FormatUint(uint64(timestamp), 10),
				"nonce":     nonce,
			})
			sign := calcBodySign(secret, client.Body())
			client.WithBody(map[string]interface{}{
				"sign": sign,
			})
		} else {
			return nil
		}

		return nil
	}

	client.WithDefaultBefores(handler)
}

func calcParamsSign(secret string, values map[string]string) string {
	// 所有参数按照字母顺序从小到大排列
	// 所有参数形成如 key1=value1key2=value2 的形式
	size := len(values)
	keys := make([]string, 0, size)
	for key := range values {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	b := buffer()
	defer releaseBuffer(b)
	b.Reset()

	for _, key := range keys {
		if key == "" {
			continue
		}

		b.WriteString(key)
		b.WriteByte('=')
		vvs := values[key]
		b.WriteString(vvs)
	}

	signStr := b.String()
	return calcWithHmacSha256(secret, signStr)
}

func calcBodySign(secret string, values map[string]interface{}) string {
	// 所有参数按照字母顺序从小到大排列
	// 所有参数形成如 key1=value1key2=value2 的形式
	size := len(values)
	keys := make([]string, 0, size)
	for key := range values {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	b := buffer()
	defer releaseBuffer(b)
	b.Reset()

	for _, key := range keys {
		if key == "" {
			continue
		}

		b.WriteString(key)
		b.WriteByte('=')
		vvs := values[key]

		switch vvs.(type) {
		case string:
			{
				b.WriteString(key)
				b.WriteByte('=')
				vvs := values[key]
				b.WriteString(vvs.(string))
			}
		case []string:
			{
				size := len(vvs.([]string))
				for idx, vv := range vvs.([]string) {
					// 相同的参数使用 , 连接
					// 例如: a=1&a=2&a=3，会构造成字符串: a=1,2,3
					b.WriteString(vv)
					if idx != size-1 {
						b.WriteByte(',')
					}
				}
			}
		}
	}

	signStr := b.String()
	return calcWithHmacSha256(secret, signStr)
}

// calcWithHmacSha256 使用 HmacSha256 进行签名
func calcWithHmacSha256(secretKey, signStr string) string {
	return zerocrypto.HmacSha256(signStr, secretKey)
}

var bufferPool *sync.Pool

func buffer() *bytes.Buffer {
	buff := bufferPool.Get().(*bytes.Buffer)
	buff.Reset()
	return buff
}

func releaseBuffer(buff *bytes.Buffer) {
	bufferPool.Put(buff)
}

func init() {
	bufferPool = &sync.Pool{}
	bufferPool.New = func() interface{} {
		return &bytes.Buffer{}
	}
}
