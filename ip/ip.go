package ip

import (
	"encoding/hex"
	"math/big"
	"net"
	"strings"
)

// ToUint64 将 IPV4 或者 IPV6 转为 uint64
// s: IPV4 或者 IPV6
func ToUint64(s string) *big.Int {
	ipV6Int := big.NewInt(0)

	ip := net.ParseIP(s)
	// 无效ip
	if ip == nil {
		return ipV6Int
	}

	ipV6Int.SetBytes(ip.To16())
	return ipV6Int
}

// ToString 将 uint64 转为 IPV4 或者 IPV6
// n: 通过 ToUint64 得到的值
// https://blog.csdn.net/swingLau007/article/details/116170235
func ToString(ipInt *big.Int) string {
	b255 := big.NewInt(0).SetBytes([]byte{255})
	var buf = make([]byte, 2)
	p := make([]string, 8)
	j := 0
	var i uint
	tmpint := big.NewInt(0)
	for i = 0; i < 16; i += 2 {
		tmpint.Rsh(ipInt, 120-i*8).And(tmpint, b255)
		bytes := tmpint.Bytes()
		if len(bytes) > 0 {
			buf[0] = bytes[0]
		} else {
			buf[0] = 0
		}
		tmpint.Rsh(ipInt, 120-(i+1)*8).And(tmpint, b255)
		bytes = tmpint.Bytes()
		if len(bytes) > 0 {
			buf[1] = bytes[0]
		} else {
			buf[1] = 0
		}
		p[j] = hex.EncodeToString(buf)
		j++
	}

	return strings.Join(p, ":")
}

// ToIPString 根据 IP 类型，输出 x.x.x.x 和 x:x:x:x
func ToIPString(ipInt *big.Int) string {
	s := ToString(ipInt)
	ip := net.ParseIP(s)

	for _, c := range s {
		switch c {
		case ':':
			return ip.To16().String()
		case '.':
			return ip.To4().String()
		}
	}

	return s
}

// GetLocalAddr 获取本地地址，只限 ipv4
func GetLocalAddr() ([]string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return nil, err
	}

	var out []string

	for _, addr := range addrs {
		if n, ok := addr.(*net.IPNet); ok && !n.IP.IsLoopback() {
			if n.IP.To4() != nil /** || n.IP.To16() != nil */ {
				out = append(out, n.IP.String())
			}
		}
	}

	return out, nil
}
