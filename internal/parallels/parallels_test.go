package parallels

import (
	"fmt"
	"image"
	"sync"
	"testing"
)

func parallel2(start, stop, n int, fn func(<-chan int)) {
	// https://github.com/disintegration/imaging/blob/d40f48ce0f098c53ab1fcd6e0e402da682262da5/utils.go#L20-L50
	count := stop - start
	c := make(chan int, count)
	for i := start; i < stop; i++ {
		c <- i
	}
	close(c)

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn(c)
		}()
	}
	wg.Wait()
}

// RunWorkers executes the specified worker function using n goroutines, passing
// each a workerNum from 0-n and a workerCount of n. This function returns after
// all workers have run to completion.
//
// https://github.com/mandykoh/go-parallel/blob/467982c50985631149f6cb6d7d5a34541703b2f4/parallel.go#L7-L22
func RunWorkers(n int, worker func(workerNum, workerCount int)) {
	allDone := sync.WaitGroup{}
	allDone.Add(n)

	for workerNum := 0; workerNum < n; workerNum++ {
		go func(workerNum int) {
			defer allDone.Done()
			worker(workerNum, n)
		}(workerNum)
	}

	allDone.Wait()
}

func parallel4(from, to, n int, f func(i int)) {
	if n <= 1 {
		for i := from; i < to; i++ {
			f(i)
		}
		return
	}
	var wg sync.WaitGroup
	wg.Add(n)
	step := (to - from + n - 1) / n
	for i := from; i < to; i += step {
		go func(from, to int) {
			defer wg.Done()
			for i := from; i < to; i++ {
				f(i)
			}
		}(i, min(i+step, to))
	}
	wg.Wait()
}

func BenchmarkParallel(b *testing.B) {
	src := image.NewNRGBA64(image.Rect(0, 0, 512, 512))
	dst := image.NewNRGBA64(image.Rect(0, 0, 512, 512))

	for n := 1; n <= 16; n++ {
		n := n
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Parallel(dst.Bounds().Min.Y, dst.Bounds().Max.Y, n, func(from, to int) {
					for y := from; y < to; y++ {
						for x := 0; x < 512; x++ {
							dst.SetNRGBA64(x, y, src.NRGBA64At(x, y))
						}
					}
				})
			}
		})
	}
}

func BenchmarkParallel2(b *testing.B) {
	src := image.NewNRGBA64(image.Rect(0, 0, 512, 512))
	dst := image.NewNRGBA64(image.Rect(0, 0, 512, 512))

	for n := 1; n <= 16; n++ {
		n := n
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				parallel2(dst.Bounds().Min.Y, dst.Bounds().Max.Y, n, func(ch <-chan int) {
					for y := range ch {
						for x := 0; x < 512; x++ {
							dst.SetNRGBA64(x, y, src.NRGBA64At(x, y))
						}
					}
				})
			}
		})
	}
}

func BenchmarkParallel3(b *testing.B) {
	src := image.NewNRGBA64(image.Rect(0, 0, 512, 512))
	dst := image.NewNRGBA64(image.Rect(0, 0, 512, 512))

	for n := 1; n <= 16; n++ {
		n := n
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				RunWorkers(n, func(num, count int) {
					for y := num; y < 512; y += count {
						for x := 0; x < 512; x++ {
							dst.SetNRGBA64(x, y, src.NRGBA64At(x, y))
						}
					}
				})
			}
		})
	}
}

func BenchmarkParallel4(b *testing.B) {
	src := image.NewNRGBA64(image.Rect(0, 0, 512, 512))
	dst := image.NewNRGBA64(image.Rect(0, 0, 512, 512))

	for n := 1; n <= 16; n++ {
		n := n
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				parallel4(dst.Bounds().Min.Y, dst.Bounds().Max.Y, n, func(y int) {
					for x := 0; x < 512; x++ {
						dst.SetNRGBA64(x, y, src.NRGBA64At(x, y))
					}
				})
			}
		})
	}
}
