package jwt

import (
	"errors"
	"fmt"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
)

// JWT ..
type JWT interface {
	// Token 生成 token
	Token(data ...map[string]interface{}) (string, error)

	// TokenWithKey 生成 token
	TokenWithKey(key []byte, data ...map[string]interface{}) (string, error)

	// Verify 验证 token，并获取自定义内容
	Verify(s string) (map[string]interface{}, error)

	// VerifyWithKey 验证 token，并获取自定义内容
	VerifyWithKey(key []byte, s string) (map[string]interface{}, error)
}

type jwt struct {
	opt Option
}

// NewJWT 创建一个 jwt
func NewJWT(opts ...Option) JWT {
	jwt := &jwt{
		opt: defaultOption(),
	}

	if len(opts) > 0 {
		jwt.opt = opts[0]
	}

	return jwt
}

// Token 生成 token
func (jwt *jwt) Token(data ...map[string]interface{}) (string, error) {
	return jwt.createToken(jwt.opt.Secret, data...)
}

// TokenWithKey 生成 token
func (jwt *jwt) TokenWithKey(key []byte, data ...map[string]interface{}) (string, error) {
	return jwt.createToken(key, data...)
}

// Verify 验证 token，并获取自定义内容
func (jwt *jwt) Verify(s string) (map[string]interface{}, error) {
	return jwt.verify(jwt.opt.Secret, s)
}

// VerifyWithKey 验证 token，并获取自定义内容
func (jwt *jwt) VerifyWithKey(key []byte, s string) (map[string]interface{}, error) {
	return jwt.verify(key, s)
}

func (jwt *jwt) createToken(key []byte, data ...map[string]interface{}) (string, error) {
	claims := jwtgo.MapClaims{
		// 签发时间
		"iat": time.Now().Unix(),
		// 过期时间
		"exp": time.Now().Add(jwt.opt.Exp).Unix(),
		// 签发者
		"iss": jwt.opt.ISS,
	}

	if len(data) > 0 {
		for k, v := range data[0] {
			claims[k] = v
		}
	}

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)

	return token.SignedString(key)
}

func (jwt *jwt) verify(key []byte, s string) (map[string]interface{}, error) {
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
