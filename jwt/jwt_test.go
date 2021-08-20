package jwt_test

import (
	"testing"

	zerojwt "github.com/zerogo-hub/zero-helper/jwt"
	zerotime "github.com/zerogo-hub/zero-helper/time"
)

func TestToken(t *testing.T) {

	j := zerojwt.NewJWT(zerojwt.Option{
		Secret: []byte("12345"),
		Exp:    zerotime.Minute(5),
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
