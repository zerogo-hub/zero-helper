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

type CalcSignHandler func(client *HTTPClient, secret string, values map[string]interface{}) (string, error)

// WithSign 自动签名中间件，自动添加 timestamp, nonce, sign 值
// secret: 签名使用，签名方式见 github.com/zerogo-hub/zero-api-middleware/sign
//
// zerohttpclient.WithSign(client, "secret")
func WithSign(client *HTTPClient, secret string, calcSignFN CalcSignHandler) {
	handler := func(client *HTTPClient, method, url string) error {
		timestamp := zerotime.Now()
		nonce := zerorandom.LowerWithNumber(32)

		if method == http.MethodGet || method == http.MethodDelete {
			// 如果有缓存，需要事先设置缓存 key，否则自动生成的缓存 key 包含 timestamp/nonce 会使得缓存一直失效
			if method == http.MethodGet && client.cacheAble() {
				cacheKey := calcCacheKey(client, url)
				client.WithCacheKey(cacheKey)
			}

			// 参数添加到 params 中
			client.WithParams(map[string]interface{}{
				"timestamp": timestamp,
				"nonce":     nonce,
			})

			var sign string
			var err error

			if calcSignFN != nil {
				sign, err = calcSignFN(client, secret, client.Params())
			} else {
				sign, err = calcSign(client, secret, client.Params())
			}

			if err != nil {
				return err
			}

			client.WithParams(map[string]interface{}{
				"sign": sign,
			})
		} else if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
			// 参数添加到 body 中
			client.WithBody(map[string]interface{}{
				"timestamp": strconv.FormatUint(uint64(timestamp), 10),
				"nonce":     nonce,
			})

			var sign string
			var err error

			if calcSignFN != nil {
				sign, err = calcSignFN(client, secret, client.Body())
			} else {
				sign, err = calcSign(client, secret, client.Body())
			}

			if err != nil {
				return err
			}

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

// calcCacheKey 在 GET 请求下生成缓存 key
func calcCacheKey(client *HTTPClient, url string) string {
	b := buffer()
	defer releaseBuffer(b)
	b.Reset()

	b.WriteString(url)
	b.WriteByte('?')

	values, _ := client.ToURLValues(client.Params())

	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		if key == "" {
			continue
		}

		b.WriteString(key)
		b.WriteByte('=')
		vvs := values[key]
		for idx, vv := range vvs {
			// 相同的参数使用 , 连接
			// 例如: a=1&a=2&a=3，会构造成字符串: a=1,2,3
			b.WriteString(vv)
			if idx != len(vvs)-1 {
				b.WriteByte(',')
			}
		}
		b.WriteByte('&')
	}

	return b.String()
}

func calcSign(client *HTTPClient, secret string, values map[string]interface{}) (string, error) {
	// 所有参数按照字母顺序从小到大排列
	// 所有参数形成如 key1=value1key2=value2 的形式
	parsedValues, err := client.ToURLValues(values)
	if err != nil {
		return "", err
	}

	size := len(parsedValues)
	keys := make([]string, 0, size)
	for key := range parsedValues {
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
		vvs := parsedValues[key]
		for idx, vv := range vvs {
			// 相同的参数使用 , 连接
			// 例如: a=1&a=2&a=3，会构造成字符串: a=1,2,3
			b.WriteString(vv)
			if idx != len(vvs)-1 {
				b.WriteByte(',')
			}
		}
	}

	signStr := b.String()
	return calcWithHmacSha256(secret, signStr), nil
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
