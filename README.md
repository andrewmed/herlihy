# Various synchronized implementations for a linked list, in Go

Implementation of ch. 9 of Prof. Herlihy's ["The Art of Multiprocessor Programming"](https://www.amazon.com/Art-Multiprocessor-Programming-Revised-Reprint/dp/0123973376), plus a Go-specific channels version

Done as a self-preparation for [Hydraconf](https://hydraconf.com/)

## Run
`go test -bench .`

## Benchmarks
On a 4-core machine it gives:

```
BenchmarkRaces/1-4                         30000             52828 ns/op
BenchmarkRaces/4-4                         30000             56115 ns/op
BenchmarkRaces/10-4                        30000             59021 ns/op
BenchmarkRaces/100-4                       20000             94784 ns/op
BenchmarkRaces/1000-4                       3000            436254 ns/op
BenchmarkCoarseLock/1-4                     5000            309904 ns/op
BenchmarkCoarseLock/4-4                     3000            521100 ns/op
BenchmarkCoarseLock/10-4                    2000            552875 ns/op
BenchmarkCoarseLock/100-4                   2000           1546961 ns/op
BenchmarkCoarseLock/1000-4                  1000           2037702 ns/op
BenchmarkFineLock/1-4                      10000            200247 ns/op
BenchmarkFineLock/4-4                       5000            344624 ns/op
BenchmarkFineLock/10-4                      5000            371864 ns/op
BenchmarkFineLock/100-4                     2000            871999 ns/op
BenchmarkFineLock/1000-4                     500           3266271 ns/op
BenchmarkOptimisticLock/1-4                10000            141566 ns/op
BenchmarkOptimisticLock/4-4                10000            181976 ns/op
BenchmarkOptimisticLock/10-4               10000            192385 ns/op
BenchmarkOptimisticLock/100-4               5000            287794 ns/op
BenchmarkOptimisticLock/1000-4              2000            811499 ns/op
BenchmarkLazyLock/1-4                      10000            114197 ns/op
BenchmarkLazyLock/4-4                      10000            150818 ns/op
BenchmarkLazyLock/10-4                     10000            157579 ns/op
BenchmarkLazyLock/100-4                     5000            251170 ns/op
BenchmarkLazyLock/1000-4                    2000            754783 ns/op
BenchmarkChan/1-4                            500           3138205 ns/op
BenchmarkChan/4-4                            500           2868732 ns/op
BenchmarkChan/10-4                           500           2918365 ns/op
BenchmarkChan/100-4                          500           3162133 ns/op
BenchmarkChan/1000-4                         300           4506547 ns/op
BenchmarkAtomic/1-4                        10000            111681 ns/op
BenchmarkAtomic/4-4                        10000            215645 ns/op
BenchmarkAtomic/10-4                       10000            230367 ns/op
BenchmarkAtomic/100-4                       5000            221398 ns/op
BenchmarkAtomic/1000-4                      3000            523708 ns/op

```

We can see that atomic scales better for a higher number of routines
Go-specific channels version scales particularly bad because message passing implies only one worker


p.s. Prof. Herlihy did not review this code, so suboptimalities are possible.


