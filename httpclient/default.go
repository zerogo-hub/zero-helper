package httpclient

import (
	"net/http"
	"time"

	zerologger "github.com/zerogo-hub/zero-helper/logger"
)

// DefaultClient 默认全局对象
var DefaultClient *HTTPClient

// WithProxy 设置代理地址
func WithProxy(proxy string) *HTTPClient {
	return DefaultClient.WithProxy(proxy)
}

// WithDefaultHeaders 设置默认消息头
func WithDefaultHeaders(headers map[string]string) *HTTPClient {
	return DefaultClient.WithDefaultHeaders(headers)
}

// WithHeaders 设置本次调用的消息头
func WithHeaders(headers map[string]string) *HTTPClient {
	return DefaultClient.WithHeaders(headers)
}

// WithHeader 设置本次调用的消息头
func WithHeader(k, v string) *HTTPClient {
	return DefaultClient.WithHeader(k, v)
}

// WithCookie 设置本次调用的 cookie
func WithCookie(cookies ...*http.Cookie) *HTTPClient {
	return DefaultClient.WithCookie(cookies...)
}

// WithParams 设置 params
func WithParams(params map[string]string) *HTTPClient {
	return DefaultClient.WithParams(params)
}

// WithBody 设置 body
// body 格式: map[string]string 或者 map[string][]string
func WithBody(body map[string]interface{}) *HTTPClient {
	return DefaultClient.WithBody(body)
}

// WithLock 手动加锁，会在调用结束时自动解锁
func WithLock() *HTTPClient {
	return DefaultClient.WithLock()
}

// WithDebug 设置调试开关
func WithDebug(isDebug bool) *HTTPClient {
	return DefaultClient.WithDebug(isDebug)
}

// WithLogger 设置日志
func WithLogger(logger zerologger.Logger) *HTTPClient {
	return DefaultClient.WithLogger(logger)
}

// WithContextTypeJSON 设置 Context-Type 格式
func WithContextTypeJSON() *HTTPClient {
	return DefaultClient.WithContextTypeJSON()
}

// WithContextTypeURLEncoded 设置 Context-Type 格式
func WithContextTypeURLEncoded() *HTTPClient {
	return DefaultClient.WithContextTypeURLEncoded()
}

// WithDialTimeout 设置拨号的超时时间
func WithDialTimeout(timeout time.Duration) *HTTPClient {
	return DefaultClient.WithDialTimeout(timeout)
}

// WithTimeout 设置每次调用的超时时间
func WithTimeout(timeout time.Duration) *HTTPClient {
	return DefaultClient.WithTimeout(timeout)
}

// WithDefaultCache 设置全局缓存开关
func WithDefaultCache(isCache bool) *HTTPClient {
	return DefaultClient.WithDefaultCache(isCache)
}

// WithCache 设置缓存开关
func WithCache(isCache bool) *HTTPClient {
	return DefaultClient.WithCache(isCache)
}

// WithCacheTTL 设置缓存时长
func WithCacheTTL(ttl time.Duration) *HTTPClient {
	return DefaultClient.WithCacheTTL(ttl)
}

// Get .
func Get(url string) *Context {
	return DefaultClient.Get(url)
}

// Post .
func Post(url string) *Context {
	return DefaultClient.Post(url)
}

// Put .
func Put(url string) *Context {
	return DefaultClient.Put(url)
}

// Patch .
func Patch(url string) *Context {
	return DefaultClient.Patch(url)
}

// Delete .
func Delete(url string) *Context {
	return DefaultClient.Delete(url)
}

// Options .
func Options(url string) *Context {
	return DefaultClient.Options(url)
}

// Connect .
func Connect(url string) *Context {
	return DefaultClient.Connect(url)
}

// Trace .
func Trace(url string) *Context {
	return DefaultClient.Trace(url)
}

func init() {
	DefaultClient = NewClient()
}
