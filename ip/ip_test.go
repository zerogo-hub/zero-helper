package ip_test

import (
	"testing"

	zeroip "github.com/zerogo-hub/zero-helper/ip"
)

func TestIPV4(t *testing.T) {
	ipV4 := "192.168.1.1"

	n := zeroip.ToUint64(ipV4)
	s := zeroip.ToIPString(n)
	if s != ipV4 {
		t.Fatalf("s must be %s, now: %s", ipV4, s)
	}
}

func TestIPV6(t *testing.T) {
	ipV6 := "2801:0137:0000:0000:0000:ffff:ffff:ffff"

	n := zeroip.ToUint64(ipV6)
	s := zeroip.ToString(n)
	if s != ipV6 {
		t.Fatalf("s must be %s, now: %s", ipV6, s)
	}
}
