# benchmark 结果

```text
> go test -bench .

goos: windows
goarch: amd64
pkg: github.com/zerogo-hub/zero-helper/codec
cpu: 12th Gen Intel(R) Core(TM) i7-12700F

BenchmarkMarshal_LibJSON-20              9258480               124.9 ns/op
BenchmarkMarshal_ZeroJSON-20             9724638               126.7 ns/op
BenchmarkMarshal_MsgPack-20              7323261               164.5 ns/op
BenchmarkMarshal_Protobuf-20            17642466                69.60 ns/op

BenchmarkUnmarshal_LibJSON-20            2282698               512.2 ns/op
BenchmarkUnmarshal_ZeroJSON-20          10912116               108.9 ns/op
BenchmarkUnmarshal_Msgpack-20            6285486               190.0 ns/op
BenchmarkUnmarshal_Protobuf-20          24339882                49.18 ns/op
PASS
ok      github.com/zerogo-hub/zero-helper/codec 11.332s
```
