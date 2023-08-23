package parallels

import "sync"

func Parallel(from, to, n int, f func(from, to int)) {
	if n <= 1 {
		f(from, to)
		return
	}
	var wg sync.WaitGroup
	wg.Add(n)
	step := (to - from + n - 1) / n
	for i := from; i < to; i += step {
		go func(from, to int) {
			defer wg.Done()
			f(from, to)
		}(i, min(i+step, to))
	}
	wg.Wait()
}
