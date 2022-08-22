package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	bigcache "github.com/allegro/bigcache/v3"

	zerologger "github.com/zerogo-hub/zero-helper/logger"
)

// BeforeHandler 执行前中间件
type BeforeHandler func(client *HTTPClient, method, url string) error

// AfterHandler 执行后中间件
type AfterHandler func(ctx *Context)

// HTTPClient 封装客户端，可以使用该客户端重复调用
type HTTPClient struct {
	// 代理地址
	proxy string

	// 默认消息头，一开始就设置好，每次调用都会传入
	defaultHeaders map[string]string

	// 消息头，需手动添加 WithHeader
	headers map[string]string

	// cookie，需手动添加 WithCookie
	cookies []*http.Cookie

	// params 添加到地址 url 上
	params map[string]string

	// body 添加到请求体 body 上
	body map[string]interface{}

	// 可复用
	transport http.RoundTripper

	// body 格式，默认为 application/x-www-form-urlencoded
	// WithContextTypeJSON
	// WithContextTypeURLEncoded
	contentType string

	// 手动加锁
	lock *sync.Mutex

	isLocked bool

	// 是否打印日志
	isDebug bool

	// 打印日志
	logger zerologger.Logger

	// 连接超时时间，默认 2 秒
	dialTimeout time.Duration

	// 每次调用的超时时间，默认 5 秒
	timeout time.Duration

	// 是否启用全局缓存，只对 GET 请求有效，只需要设置一次
	isDefaultCache bool

	// 是否启用缓存，只对 GET 请求有效，需要每次调用时设置
	isCache bool

	// 缓存过期时间，默认 60 秒
	cacheTTL time.Duration

	// 内存型缓存
	cache *bigcache.BigCache

	// 调用 cacheKey 生成，也可以实现设置
	cacheKey_ string

	// 在执行前调用，比如可以额外添加参数，默认，先于 befores 执行
	defaultBefores []BeforeHandler

	// 在执行前调用，比如可以额外添加参数，每次调用前手动设置
	befores []BeforeHandler

	// 在执行后调用，默认，先于 afters 执行
	defaultAfters []AfterHandler

	// 每次调用前手动设置
	afters []AfterHandler
}

// NewClient .
func NewClient() *HTTPClient {
	return &HTTPClient{
		contentType: "application/x-www-form-urlencoded",
		dialTimeout: time.Second * time.Duration(2),
		timeout:     time.Second * time.Duration(5),
	}
}

// Params 获取参数
func (client *HTTPClient) Params() map[string]string {
	return client.params
}

// Body 获取 body
func (client *HTTPClient) Body() map[string]interface{} {
	return client.body
}

// WithProxy 设置代理地址
func (client *HTTPClient) WithProxy(proxy string) *HTTPClient {
	client.proxy = proxy

	return client
}

// WithDefaultHeaders 设置默认消息头
func (client *HTTPClient) WithDefaultHeaders(headers map[string]string) *HTTPClient {
	if client.defaultHeaders == nil {
		client.defaultHeaders = headers
	} else {
		for k, v := range headers {
			client.defaultHeaders[k] = v
		}
	}

	return client
}

// WithHeaders 设置本次调用的消息头
func (client *HTTPClient) WithHeaders(headers map[string]string) *HTTPClient {
	if client.headers == nil {
		client.headers = make(map[string]string)
	}
	for k, v := range headers {
		client.headers[k] = v
	}

	return client
}

// WithHeader 设置本次调用的消息头
func (client *HTTPClient) WithHeader(k, v string) *HTTPClient {
	if client.headers == nil {
		client.headers = make(map[string]string)
	}
	client.headers[k] = v

	return client
}

// WithCookie 设置本次调用的 cookie
func (client *HTTPClient) WithCookie(cookies ...*http.Cookie) *HTTPClient {
	client.cookies = append(client.cookies, cookies...)

	return client
}

// WithParams 设置 params, GET, Delete 时使用
func (client *HTTPClient) WithParams(params map[string]string) *HTTPClient {
	if client.params == nil {
		client.params = make(map[string]string)
	}

	for k, v := range params {
		client.params[k] = v
	}

	return client
}

// WithBody 设置 body，Post, Put, Patch 时使用
// body 格式: map[string]string 或者 map[string][]string
func (client *HTTPClient) WithBody(body map[string]interface{}) *HTTPClient {
	if client.body == nil {
		client.body = make(map[string]interface{})
	}

	for k, v := range body {
		client.body[k] = v
	}

	return client
}

// WithLock 手动加锁，会在调用结束时自动解锁
func (client *HTTPClient) WithLock() *HTTPClient {
	client.lock.Lock()
	client.isLocked = true

	return client
}

// WithDebug 设置调试开关
func (client *HTTPClient) WithDebug(isDebug bool) *HTTPClient {
	client.isDebug = isDebug

	return client
}

// WithLogger 设置日志
func (client *HTTPClient) WithLogger(logger zerologger.Logger) *HTTPClient {
	client.logger = logger

	return client
}

// WithContextTypeJSON 设置 Context-Type 格式
func (client *HTTPClient) WithContextTypeJSON() *HTTPClient {
	client.contentType = "application/json;charset=utf-8"

	return client
}

// WithContextTypeURLEncoded 设置 Context-Type 格式
func (client *HTTPClient) WithContextTypeURLEncoded() *HTTPClient {
	client.contentType = "application/x-www-form-urlencoded"

	return client
}

// WithDialTimeout 设置拨号的超时时间
func (client *HTTPClient) WithDialTimeout(timeout time.Duration) *HTTPClient {
	client.dialTimeout = timeout

	return client
}

// WithTimeout 设置每次调用的超时时间
func (client *HTTPClient) WithTimeout(timeout time.Duration) *HTTPClient {
	client.timeout = timeout

	return client
}

// WithDefaultCache 设置全局缓存开关
func (client *HTTPClient) WithDefaultCache(isCache bool) *HTTPClient {
	client.isDefaultCache = isCache

	return client
}

// WithCache 设置缓存开关
func (client *HTTPClient) WithCache(isCache bool) *HTTPClient {
	client.isCache = isCache

	return client
}

// WithCacheTTL 设置缓存时长
func (client *HTTPClient) WithCacheTTL(ttl time.Duration) *HTTPClient {
	client.cacheTTL = ttl

	return client
}

func (client *HTTPClient) WithCacheKey(cacheKey string) *HTTPClient {
	client.cacheKey_ = cacheKey

	return client
}

// WithDefaultBefores 设置默认执行前函数
func (client *HTTPClient) WithDefaultBefores(handlers ...BeforeHandler) *HTTPClient {
	client.defaultBefores = append(client.defaultBefores, handlers...)

	return client
}

// WithBefores 设置执行前函数
func (client *HTTPClient) WithBefores(handlers ...BeforeHandler) *HTTPClient {
	client.befores = append(client.befores, handlers...)

	return client
}

// WithDefaultAfters 设置默认执行后函数
func (client *HTTPClient) WithDefaultAfters(handler ...AfterHandler) *HTTPClient {
	client.defaultAfters = append(client.defaultAfters, handler...)

	return client
}

// WithAfters 设置执行后函数
func (client *HTTPClient) WithAfters(handlers ...AfterHandler) *HTTPClient {
	client.afters = append(client.afters, handlers...)

	return client
}

// Get .
func (client *HTTPClient) Get(url string) *Context {
	return client.do(http.MethodGet, url)
}

// Post .
func (client *HTTPClient) Post(url string) *Context {
	return client.do(http.MethodPost, url)
}

// Put .
func (client *HTTPClient) Put(url string) *Context {
	return client.do(http.MethodPut, url)
}

// Patch .
func (client *HTTPClient) Patch(url string) *Context {
	return client.do(http.MethodPatch, url)
}

// Delete .
func (client *HTTPClient) Delete(url string) *Context {
	return client.do(http.MethodDelete, url)
}

// Options .
func (client *HTTPClient) Options(url string) *Context {
	return client.do(http.MethodOptions, url)
}

// Connect .
func (client *HTTPClient) Connect(url string) *Context {
	return client.do(http.MethodConnect, url)
}

// Trace .
func (client *HTTPClient) Trace(url string) *Context {
	return client.do(http.MethodTrace, url)
}

// reset 每次调用结束后重置
func (client *HTTPClient) reset() {
	client.headers = nil
	client.cookies = nil
	client.params = nil
	client.body = nil
	client.isCache = false
	client.cacheKey_ = ""
	if client.isLocked {
		client.isLocked = false
		client.lock.Unlock()
	}
}

func (client *HTTPClient) runBefores(method, url string) *Context {
	for _, handler := range client.defaultBefores {
		if err := handler(client, method, url); err != nil {
			return newContextWithError(err)
		}
	}

	for _, handler := range client.befores {
		if err := handler(client, method, url); err != nil {
			return newContextWithError(err)
		}
	}

	return nil
}

func (client *HTTPClient) runAfters(ctx *Context) {
	for _, handler := range client.defaultAfters {
		handler(ctx)
	}

	for _, handler := range client.afters {
		handler(ctx)
	}
}

// do 执行
func (client *HTTPClient) do(method, url string) *Context {
	if ctx := client.runBefores(method, url); ctx != nil {
		return ctx
	}

	headers := client.headers
	cookies := client.cookies
	params := client.params
	body := client.body
	isCache := client.isCache
	cacheKey_ := client.cacheKey_

	client.reset()

	// transport
	if client.transport == nil {
		transport, err := client.prepareTransport()
		if err != nil {
			return newContextWithError(err)
		}
		client.transport = transport
	}

	req, err := client.prepareRequest(method, url, headers, params, body, cookies)
	if err != nil {
		return newContextWithError(err)
	}

	// 缓存
	cacheAble := method == "GET" && (client.isDefaultCache || isCache)
	if cacheAble {
		if ctx := client.getFromCache(req, cacheKey_); ctx != nil {
			client.runAfters(ctx)
			return ctx
		}
	}

	if client.isDebug {
		dump, err := httputil.DumpRequestOut(req, true)
		if err == nil {
			if client.logger != nil {
				client.logger.Info(dump)
			} else {
				fmt.Printf("%s\n", dump)
			}
		}
	}

	c := &http.Client{
		Transport: client.transport,
		Timeout:   client.timeout,
	}

	resp, err := c.Do(req)
	if err != nil {
		return newContextWithError(err)
	}

	ctx := newContext(req, resp)

	if cacheAble {
		client.setToCache(ctx, cacheKey_)
	}

	client.runAfters(ctx)
	return ctx
}

func (client *HTTPClient) prepareTransport() (*http.Transport, error) {
	transport := &http.Transport{}

	transport.Dial = func(network, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, addr, client.dialTimeout)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}

	if client.proxy != "" {
		proxy, err := url.Parse(client.proxy)
		if err != nil {
			return nil, err
		}
		transport.Proxy = http.ProxyURL(proxy)
	}

	return transport, nil
}

func (client *HTTPClient) prepareRequest(method, url string, headers, params map[string]string, body map[string]interface{}, cookies []*http.Cookie) (*http.Request, error) {
	var reader io.Reader
	if body != nil {
		var err error
		reader, err = client.parseBody(body)
		if err != nil {
			return nil, err
		}
		if headers == nil {
			headers = make(map[string]string)
		}
		headers["Content-Type"] = client.contentType
	} else {
		reader = nil
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}

	if params != nil {
		query := req.URL.Query()
		for k, v := range params {
			query.Add(k, v)
		}
		req.URL.RawQuery = query.Encode()
	}

	for k, v := range client.defaultHeaders {
		req.Header.Set(k, v)
	}

	for k, v := range client.headers {
		req.Header.Set(k, v)
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	return req, nil
}

func (client *HTTPClient) parseBody(body map[string]interface{}) (io.Reader, error) {
	if strings.HasPrefix(client.contentType, "application/json") {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		return bytes.NewReader(b), nil
	} else if strings.HasPrefix(client.contentType, "application/x-www-form-urlencoded") {
		values := url.Values{}
		for k, v := range body {
			switch v.(type) {
			case map[string]string:
				values.Set(k, v.(string))
			case map[string][]string:
				for _, vv := range v.([]string) {
					values.Set(k, vv)
				}
			default:
				return nil, errors.New("invalid body type")
			}
		}

		return strings.NewReader(values.Encode()), nil
	}

	return nil, nil
}

func (client *HTTPClient) cacheKey(req *http.Request) string {
	if client.cacheKey_ != "" {
		return client.cacheKey_
	}

	return req.URL.String()
}

func (client *HTTPClient) cacheAble() bool {
	return client.isDefaultCache || client.isCache
}

func (client *HTTPClient) getFromCache(req *http.Request, cacheKey_ string) *Context {
	if !client.ensureCache() {
		return nil
	}

	var key string
	if cacheKey_ != "" {
		key = cacheKey_
	} else {
		key = client.cacheKey(req)
	}

	b, err := client.cache.Get(key)
	if err != nil {
		return nil
	}

	return newContextWithCache(req, b)
}

func (client *HTTPClient) setToCache(ctx *Context, cacheKey_ string) {
	if !client.ensureCache() {
		return
	}

	b, err := ctx.ToByes()
	if err != nil {
		return
	}

	var key string
	if cacheKey_ != "" {
		key = cacheKey_
	} else {
		key = client.cacheKey(ctx.req)
	}

	_ = client.cache.Set(key, b)
}

// DelectCache 删除缓存
func (client *HTTPClient) DelectCache(key string) {
	if !client.ensureCache() {
		return
	}

	_ = client.cache.Delete(key)
}

func (client *HTTPClient) ensureCache() bool {
	if client.cache != nil {
		return true
	}

	ttl := client.cacheTTL
	if ttl <= 0 {
		ttl = time.Second * time.Duration(60)
	}

	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(ttl))
	if err != nil {
		return false
	}

	client.cache = cache
	return true
}
