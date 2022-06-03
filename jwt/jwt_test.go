package jwt_test

import (
	"testing"
	"time"

	zerojwt "github.com/zerogo-hub/zero-helper/jwt"
	zerotime "github.com/zerogo-hub/zero-helper/time"
)

func TestToken(t *testing.T) {

	j := zerojwt.NewJWT(zerojwt.Option{
		Secret: []byte("12345"),
		Exp:    zerotime.Minute(5),
	})

	data := map[string]interface{}{
		"id":  "1",
		"age": "18",
	}

	token, err := j.Token(data)
	if err != nil {
		t.Fatalf("create token failed: %s", err.Error())
	}
	t.Logf("token: %s", token)

	p, err := j.Verify(token)
	if err != nil {
		t.Fatalf("verify token failed: %s", err.Error())
	}

	if p["id"] != data["id"] {
		t.Error("invalid id")
	}

	if p["age"] != data["age"] {
		t.Error("invalid age")
	}

	t.Logf("exp: %d, iat: %d", p["exp"], p["iat"])
}

func TestTokenWithKey(t *testing.T) {
	key := []byte("GIiK325IynHKxEAZ")

	j := zerojwt.NewJWT(zerojwt.Option{
		Exp: zerotime.Minute(5),
	})

	data := map[string]interface{}{
		"id":  "1",
		"age": "18",
	}

	token, err := j.TokenWithKey(key, data)
	if err != nil {
		t.Fatalf("create token failed: %s", err.Error())
	}
	t.Logf("token: %s", token)

	p, err := j.VerifyWithKey(key, token)
	if err != nil {
		t.Fatalf("verify token failed: %s", err.Error())
	}

	if p["id"] != data["id"] {
		t.Error("invalid id")
	}

	if p["age"] != data["age"] {
		t.Error("invalid age")
	}

	t.Logf("exp: %d, iat: %d", p["exp"], p["iat"])
}

func TestInvalidToken(t *testing.T) {

	j := zerojwt.NewJWT(zerojwt.Option{
		Secret: []byte("12345"),
		Exp:    zerotime.Millisecond(1),
	})

	data := map[string]interface{}{
		"id":  "1",
		"age": "18",
	}

	token, err := j.Token(data)
	if err != nil {
		t.Fatalf("create token failed: %s", err.Error())
	}

	time.Sleep(zerotime.Millisecond(15))

	if _, err = j.Verify(token); err != nil {
		t.Fatalf("verify token timeout failed: %s", err.Error())
	}

	// 无效 token
	if _, err = j.Verify(token + "s"); err == nil {
		t.Fatalf("verify invalid token failed: %s", err.Error())
	}
}
