package parallels

import (
	"runtime"
	"sync"
)

// Parallel runs f in parallel for i in [from, to).
func Parallel(from, to int, f func(i int)) {
	parallel(from, to, runtime.NumCPU(), f)
}

func parallel(from, to, n int, f func(i int)) {
	if n <= 1 {
		for i := from; i < to; i++ {
			f(i)
		}
		return
	}

	var wg sync.WaitGroup
	step := max((to-from+n-1)/n, 1)
	for i := from; i < to; i += step {
		wg.Add(1)
		go func(from, to int) {
			defer wg.Done()
			for i := from; i < to; i++ {
				f(i)
			}
		}(i, min(i+step, to))
	}
	wg.Wait()
}
