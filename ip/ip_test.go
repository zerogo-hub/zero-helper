package ip_test

import (
	"testing"

	zeroip "github.com/zerogo-hub/zero-helper/ip"
)

func TestIPV4(t *testing.T) {
	ipV4 := "192.168.1.1"

	n := zeroip.ToUint64(ipV4)
	s := zeroip.ToString(n)
	if s != ipV4 {
		t.Fatalf("s must be %s, now: %s", ipV4, s)
	}

	if zeroip.ToUint64("256.256.256.256").Int64() > 0 {
		t.Fatal("test invalid ip falied")
	}
}

func TestIPV6(t *testing.T) {
	ipV6 := "2801:0137:0000:0000:ffff:0000:0000:ffff"

	n := zeroip.ToUint64(ipV6)
	s := zeroip.ToString(n)
	if s != "2801:137::ffff:0:0:ffff" {
		t.Fatalf("s must be %s, now: %s", ipV6, s)
	}
}

func TestIP(t *testing.T) {
	_, err := zeroip.GetLocalAddr()
	if err != nil {
		t.Fatalf("get ips failed: %s", err.Error())
	}
}
