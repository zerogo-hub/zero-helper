package jwt

import (
	"errors"
	"fmt"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
)

// JWT ..
type JWT struct {
	opt Option
}

var (
	defaultJWT *JWT
)

// Token 生成 token
func Token(data ...map[string]interface{}) (string, error) {
	return defaultJWT.Token(data...)
}

// TokenWithKey 生成 token
func TokenWithKey(key []byte, data ...map[string]interface{}) (string, error) {
	return defaultJWT.TokenWithKey(key, data...)
}

// Verify 验证 token，并获取自定义内容
func Verify(s string) (map[string]interface{}, error) {
	return defaultJWT.Verify(s)
}

// VerifyWithKey 验证 token，并获取自定义内容
func VerifyWithKey(key []byte, s string) (map[string]interface{}, error) {
	return defaultJWT.VerifyWithKey(key, s)
}

// NewJWT 创建一个 jwt
func NewJWT(opts ...Option) *JWT {
	jwt := &JWT{
		opt: defaultOption(),
	}

	if len(opts) > 0 {
		jwt.opt = opts[0]
	}

	return jwt
}

// Token 生成 token
func (jwt *JWT) Token(data ...map[string]interface{}) (string, error) {
	return jwt.createToken(jwt.opt.Secret, data...)
}

// TokenWithKey 生成 token
func (jwt *JWT) TokenWithKey(key []byte, data ...map[string]interface{}) (string, error) {
	return jwt.createToken(key, data...)
}

// Verify 验证 token，并获取自定义内容
func (jwt *JWT) Verify(s string) (map[string]interface{}, error) {
	return jwt.verify(jwt.opt.Secret, s)
}

// VerifyWithKey 验证 token，并获取自定义内容
func (jwt *JWT) VerifyWithKey(key []byte, s string) (map[string]interface{}, error) {
	return jwt.verify(key, s)
}

func (jwt *JWT) createToken(key []byte, data ...map[string]interface{}) (string, error) {
	claims := jwtgo.MapClaims{
		// 签发时间
		"iat": time.Now().Unix(),
		// 过期时间
		"exp": time.Now().Add(jwt.opt.Exp).Unix(),
	}

	if len(data) > 0 {
		for k, v := range data[0] {
			claims[k] = v
		}
	}

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)

	return token.SignedString(key)
}

func (jwt *JWT) verify(key []byte, s string) (map[string]interface{}, error) {
	parse, err := jwtgo.Parse(s, func(token *jwtgo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid jwt method")
		}

		return key, nil
	})

	if err != nil {
		return nil, errors.New("token validate failed")
	}

	if claims, ok := parse.Claims.(jwtgo.MapClaims); ok && parse.Valid {
		// 验证是否超时
		exp, ok := claims["exp"]
		if !ok {
			return nil, errors.New("no exp in token")
		}
		f64, ok := exp.(float64)
		if !ok {
			return nil, errors.New("invalid exp type")
		}

		i64exp := int64(f64)

		if i64exp < time.Now().Unix() {
			return nil, fmt.Errorf("token expired at %d", i64exp)
		}

		payload := make(map[string]interface{}, len(claims))
		for k, v := range claims {
			payload[k] = v
		}

		return payload, nil
	}

	return nil, errors.New("token not match MapClaims")
}

// SetDefaultOption 设置默认配置
func SetDefaultOption(opt Option) {
	defaultJWT.opt = opt
}

func init() {
	defaultJWT = NewJWT()
}
