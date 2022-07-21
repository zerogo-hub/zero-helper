package httpclient

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
)

// Context 一次调用过程的上下文
type Context struct {
	req  *http.Request
	resp *http.Response

	err error

	// 暂存响应内容，可以反复调用 ToBytes()
	respBody []byte

	fromCache bool
}

func newContextWithError(err error) *Context {
	return &Context{err: err}
}

func newContextWithCache(req *http.Request, b []byte) *Context {
	return &Context{
		req:       req,
		respBody:  b,
		fromCache: true,
	}
}

func newContext(req *http.Request, resp *http.Response) *Context {
	return &Context{
		req:  req,
		resp: resp,
	}
}

// Request 获取原始请求
func (ctx *Context) Request() *http.Request {
	return ctx.req
}

// Response 获取原始响应
func (ctx *Context) Response() *http.Response {
	return ctx.resp
}

// ToBytes .
func (ctx *Context) ToByes() ([]byte, error) {
	if ctx.err != nil {
		return nil, ctx.err
	}

	if ctx.respBody != nil {
		return ctx.respBody, nil
	}

	if ctx.resp == nil {
		return nil, nil
	}

	defer ctx.resp.Body.Close()
	body, err := ioutil.ReadAll(ctx.resp.Body)
	if err != nil {
		ctx.err = err
		return nil, err
	}

	ctx.respBody = body
	return ctx.respBody, nil
}

// ToString .
func (ctx *Context) ToString() (string, error) {
	bytes, err := ctx.ToByes()
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// ToJSON .
func (ctx *Context) ToJSON(out interface{}) error {
	bytes, err := ctx.ToByes()
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, out)
}

// IsOK 结果是否正常
func (ctx *Context) IsOK() bool {
	if ctx.fromCache {
		return true
	}

	return ctx.err == nil && ctx.resp != nil && ctx.resp.StatusCode == 200
}

// IsTimeout 是否超时
func (ctx *Context) IsTimeout() bool {
	if ctx.err == nil {
		return false
	}

	nerr, ok := ctx.err.(net.Error)
	if !ok {
		return false
	}

	return nerr.Timeout()
}
