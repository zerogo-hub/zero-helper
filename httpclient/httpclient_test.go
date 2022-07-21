package httpclient_test

import (
	"testing"
	"time"

	zerohttpclient "github.com/zerogo-hub/zero-helper/httpclient"
)

func TestGet(t *testing.T) {
	client := zerohttpclient.NewClient()

	url := "https://www.keylala.cn"
	ctx := client.WithCache(true).WithCacheTTL(time.Minute * time.Duration(10)).Get(url)
	if !ctx.IsOK() {
		if ctx.IsTimeout() {
			t.Fatal("Get timeout")
		}

		t.Fatal("Get failed")
	}
	ctx = client.WithCache(true).WithCacheTTL(time.Minute * time.Duration(10)).Get(url)

	result, err := ctx.ToString()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}
