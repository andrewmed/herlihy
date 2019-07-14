package herlihy

import (
	"runtime"
	"strconv"
	"sync"
	"testing"
)

func runLinear(tb testing.TB, list list, repetition int, val int, fail bool) {
	for i := 0; i < repetition; i++ {
		if !list.add(val) && fail {
			tb.Error(i, val)
		}
		if list.add(val) && fail {
			tb.Error(i, val)
		}
		if !list.contains(val) && fail {
			tb.Error(i, val)
		}
		if !list.remove(val) && fail {
			tb.Error(i, val)
		}
		if list.remove(val) && fail {
			tb.Error(i, val)
		}
		if list.contains(val) && fail {
			tb.Error(i, val)
		}
	}
}

func runConcurrent(tb testing.TB, list list, concurrency int, fail bool) {
	start := sync.WaitGroup{}
	start.Add(concurrency)
	stop := sync.WaitGroup{}
	stop.Add(concurrency)
	repeat := 1000 / concurrency
	for val := 0; val < concurrency; val++ {
		go func(val int) {
			start.Done()
			start.Wait()
			runLinear(tb, list, repeat, val, fail)
			stop.Done()
		}(val)
	}
	stop.Wait()
}

func TestCoarseLock(t *testing.T) {
	coarse := newCoarse()
	runConcurrent(t, &coarse, 100, true)
}

func TestFineLock(t *testing.T) {
	fine := newFineLocked()
	runConcurrent(t, &fine, 100, true)
}

func TestOptimisticLock(t *testing.T) {
	optimistic := newOptimistic()
	runConcurrent(t, &optimistic, 100, true)
}

func TestLazyLock(t *testing.T) {
	lazy := newLazy()
	runConcurrent(t, &lazy, 100, true)
}

func TestChan(t *testing.T) {
	lazy := newChan()
	runConcurrent(t, &lazy, 100, true)
}

func TestAtomic(t *testing.T) {
	lazy := newAtomic()
	runConcurrent(t, &lazy, 100, true)
}

func runBundle(b *testing.B, list list, fail bool) {
	var tests = []int{1, runtime.NumCPU(), 10, 100, 1000}
	for _, routines := range tests {
		b.Run(strconv.Itoa(routines), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				runConcurrent(b, list, routines, fail)
			}
		})
	}
}

func BenchmarkRaces(b *testing.B) {
	list := newRaces()
	runBundle(b, &list, false)
}

func BenchmarkCoarseLock(b *testing.B) {
	list := newCoarse()
	runBundle(b, &list, true)
}

func BenchmarkFineLock(b *testing.B) {
	list := newFineLocked()
	runBundle(b, &list, true)
}

func BenchmarkOptimisticLock(b *testing.B) {
	list := newOptimistic()
	runBundle(b, &list, true)
}

func BenchmarkLazyLock(b *testing.B) {
	list := newLazy()
	runBundle(b, &list, true)
}

func BenchmarkChan(b *testing.B) {
	list := newChan()
	runBundle(b, &list, true)
}

func BenchmarkAtomic(b *testing.B) {
	list := newAtomic()
	runBundle(b, &list, true)
}
