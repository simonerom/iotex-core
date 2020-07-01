before 100
```
go test -bench=. -benchmem -benchtime 100x
goos: darwin
goarch: amd64
pkg: github.com/iotexproject/iotex-core/db/trie/mptrie
BenchmarkTrie-8   	     100	  84320630 ns/op	   96246 B/op	     664 allocs/op
--- BENCH: BenchmarkTrie-8
    merklepatriciatrie_benchmark_test.go:123: iter 1
    merklepatriciatrie_benchmark_test.go:140: iter+3 4
    merklepatriciatrie_benchmark_test.go:168: Warning: test 4 entries
    merklepatriciatrie_benchmark_test.go:123: iter 100
    merklepatriciatrie_benchmark_test.go:140: iter+3 103
    merklepatriciatrie_benchmark_test.go:168: Warning: test 103 entries
PASS
ok  	github.com/iotexproject/iotex-core/db/trie/mptrie	55.871s
```


now 100
```
go test -bench=. -benchmem -benchtime 100x
goos: darwin
goarch: amd64
pkg: github.com/iotexproject/iotex-core/db/trie/mptrie
BenchmarkTrie-8   	     100	  83607639 ns/op	  101763 B/op	     808 allocs/op
--- BENCH: BenchmarkTrie-8
    merklepatriciatrie_benchmark_test.go:44: iter 1
    merklepatriciatrie_benchmark_test.go:62: iter+3 4
    merklepatriciatrie_benchmark_test.go:91: Warning: test 4 entries
    merklepatriciatrie_benchmark_test.go:44: iter 100
    merklepatriciatrie_benchmark_test.go:62: iter+3 103
    merklepatriciatrie_benchmark_test.go:91: Warning: test 103 entries
PASS
ok  	github.com/iotexproject/iotex-core/db/trie/mptrie	32.626s
```

before 1000
```
 go test -bench=. -benchmem -benchtime 1000x
goos: darwin
goarch: amd64
pkg: github.com/iotexproject/iotex-core/db/trie/mptrie
BenchmarkTrie-8   	    1000	 111029180 ns/op	  205267 B/op	    1198 allocs/op
--- BENCH: BenchmarkTrie-8
    merklepatriciatrie_benchmark_test.go:122: iter 1
    merklepatriciatrie_benchmark_test.go:139: iter+3 4
    merklepatriciatrie_benchmark_test.go:167: Warning: test 4 entries
    merklepatriciatrie_benchmark_test.go:122: iter 1000
    merklepatriciatrie_benchmark_test.go:139: iter+3 1003
    merklepatriciatrie_benchmark_test.go:167: Warning: test 1003 entries
PASS
ok  	github.com/iotexproject/iotex-core/db/trie/mptrie	133.808s
```

now 1000
```
go test -bench=. -benchmem -benchtime 1000x
goos: darwin
goarch: amd64
pkg: github.com/iotexproject/iotex-core/db/trie/mptrie
BenchmarkTrie-8   	    1000	 110362176 ns/op	  286114 B/op	    2751 allocs/op
--- BENCH: BenchmarkTrie-8
    merklepatriciatrie_benchmark_test.go:44: iter 1
    merklepatriciatrie_benchmark_test.go:62: iter+3 4
    merklepatriciatrie_benchmark_test.go:91: Warning: test 4 entries
    merklepatriciatrie_benchmark_test.go:44: iter 1000
    merklepatriciatrie_benchmark_test.go:62: iter+3 1003
    merklepatriciatrie_benchmark_test.go:91: Warning: test 1003 entries
PASS
ok  	github.com/iotexproject/iotex-core/db/trie/mptrie	115.199s
```

before 4000
```
go test -bench=. -benchmem -benchtime 4000x
goos: darwin
goarch: amd64
pkg: github.com/iotexproject/iotex-core/db/trie/mptrie
BenchmarkTrie-8   	    4000	 124749396 ns/op	  261113 B/op	    1497 allocs/op
--- BENCH: BenchmarkTrie-8
    merklepatriciatrie_benchmark_test.go:122: iter 1
    merklepatriciatrie_benchmark_test.go:139: iter+3 4
    merklepatriciatrie_benchmark_test.go:167: Warning: test 4 entries
    merklepatriciatrie_benchmark_test.go:122: iter 4000
    merklepatriciatrie_benchmark_test.go:139: iter+3 4003
    merklepatriciatrie_benchmark_test.go:167: Warning: test 4003 entries
PASS
ok  	github.com/iotexproject/iotex-core/db/trie/mptrie	521.103s
```

now 4000
```
go test -bench=. -benchmem -benchtime 4000x
goos: darwin
goarch: amd64
pkg: github.com/iotexproject/iotex-core/db/trie/mptrie
BenchmarkTrie-8   	    4000	 123461392 ns/op	  493906 B/op	    6360 allocs/op
--- BENCH: BenchmarkTrie-8
    merklepatriciatrie_benchmark_test.go:44: iter 1
    merklepatriciatrie_benchmark_test.go:62: iter+3 4
    merklepatriciatrie_benchmark_test.go:91: Warning: test 4 entries
    merklepatriciatrie_benchmark_test.go:44: iter 4000
    merklepatriciatrie_benchmark_test.go:62: iter+3 4003
    merklepatriciatrie_benchmark_test.go:91: Warning: test 4003 entries
PASS
ok  	github.com/iotexproject/iotex-core/db/trie/mptrie	509.923s
```
