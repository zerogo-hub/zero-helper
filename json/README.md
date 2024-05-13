# benchmark 结果

```text
> go test -bench .

goos: windows
goarch: amd64
pkg: github.com/zerogo-hub/zero-helper/json
cpu: 12th Gen Intel(R) Core(TM) i7-12700F

BenchmarkMarshal_1_StdLib-20             7648665               161.2 ns/op
BenchmarkMarshal_1_ZeroJSON-20           8444180               146.9 ns/op
BenchmarkUnmarshal_1_StdLib-20          10990096               108.2 ns/op
BenchmarkUnmarshal_1_ZeroJSON-20         8223655               144.2 ns/op

BenchmarkMarshal_100_StdLib-20            176460              6307 ns/op
BenchmarkMarshal_100_ZeroJSON-20          521434              2209 ns/op
BenchmarkUnmarshal_100_StdLib-20          387099              3178 ns/op
BenchmarkUnmarshal_100_ZeroJSON-20       1000000              1125 ns/op

BenchmarkMarsha_500_StdLib-20              39156             30621 ns/op
BenchmarkMarshal_500_ZeroJSON-20          114260             10533 ns/op
BenchmarkUnmarshal_500_StdLib-20           79960             15646 ns/op
BenchmarkUnmarshal_500_ZeroJSON-20        187497              5632 ns/op

PASS
ok      github.com/zerogo-hub/zero-helper/json  15.859s
```
