package jwt_test

import (
	"testing"

	"github.com/zerogo-hub/zero-helper/jwt"
	"github.com/zerogo-hub/zero-helper/time"
)

func TestToken(t *testing.T) {

	j := jwt.NewJWT(jwt.Option{
		Secret: []byte("12345"),
		Exp:    time.Minute(5),
	})

	token, err := j.Token()
	if err != nil {
		t.Fatalf("create token failed: %s", err.Error())
	}
	t.Logf("token: %s", token)

	p, err := j.Verify(token)
	if err != nil {
		t.Fatalf("verify token failed: %s", err.Error())
	}

	t.Logf("exp: %d, iat: %d", p["exp"], p["iat"])
}
