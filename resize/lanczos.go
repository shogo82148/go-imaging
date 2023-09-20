package resize

import (
	"image"
	"math"

	"github.com/shogo82148/go-imaging/fp16"
	"github.com/shogo82148/go-imaging/fp16/fp16color"
	"github.com/shogo82148/go-imaging/internal/parallels"
)

// Lanczos2 resizes the image using Lanczos interpolation with lobe 2.
func Lanczos2(dst, src *fp16.NRGBAh) {
	dstBounds := dst.Bounds()
	dstDx := dstBounds.Dx()
	dstDy := dstBounds.Dy()
	srcBounds := src.Bounds()
	srcDx := srcBounds.Dx()
	srcDy := srcBounds.Dy()

	tmp := fp16.NewNRGBAh(image.Rect(0, 0, dstDx, srcDy))

	// pre-calculate lanczos weights
	cache := make(map[float64]float64, 4*(dstDx+dstDy))
	for x := 0; x < dstDx; x++ {
		_, dx := scale(x, srcDx, dstDx)
		cache[dx+1] = lanczos(dx+1, 2)
		cache[dx+0] = lanczos(dx+0, 2)
		cache[dx-1] = lanczos(dx-1, 2)
		cache[dx-2] = lanczos(dx-2, 2)
	}
	for y := 0; y < dstDy; y++ {
		_, dy := scale(y, srcDy, dstDy)
		cache[dy+1] = lanczos(dy+1, 2)
		cache[dy+0] = lanczos(dy+0, 2)
		cache[dy-1] = lanczos(dy-1, 2)
		cache[dy-2] = lanczos(dy-2, 2)
	}

	// resize horizontally
	parallels.Parallel(0, srcDy, func(y int) {
		for x := 0; x < dstDx; x++ {
			var c fp16color.NRGBAh
			srcX, dx := scale(x, srcDx, dstDx)
			c0 := nrgbhAt(src, srcBounds.Min.X+srcX-1, srcBounds.Min.Y+y)
			c1 := nrgbhAt(src, srcBounds.Min.X+srcX+0, srcBounds.Min.Y+y)
			c2 := nrgbhAt(src, srcBounds.Min.X+srcX+1, srcBounds.Min.Y+y)
			c3 := nrgbhAt(src, srcBounds.Min.X+srcX+2, srcBounds.Min.Y+y)

			a0 := cache[dx+1]
			a1 := cache[dx+0]
			a2 := cache[dx-1]
			a3 := cache[dx-2]

			c.R = product4(a0, a1, a2, a3, c0.R, c1.R, c2.R, c3.R)
			c.G = product4(a0, a1, a2, a3, c0.G, c1.G, c2.G, c3.G)
			c.B = product4(a0, a1, a2, a3, c0.B, c1.B, c2.B, c3.B)
			c.A = product4(a0, a1, a2, a3, c0.A, c1.A, c2.A, c3.A)
			tmp.SetNRGBAh(x, y, c)
		}
	})

	// resize vertically
	parallels.Parallel(0, dstDy, func(y int) {
		for x := 0; x < dstDx; x++ {
			var c fp16color.NRGBAh
			srcY, dy := scale(y, srcDy, dstDy)
			c0 := nrgbhAt(tmp, x, srcY-1)
			c1 := nrgbhAt(tmp, x, srcY+0)
			c2 := nrgbhAt(tmp, x, srcY+1)
			c3 := nrgbhAt(tmp, x, srcY+2)

			a0 := cache[dy+1]
			a1 := cache[dy+0]
			a2 := cache[dy-1]
			a3 := cache[dy-2]

			c.R = product4(a0, a1, a2, a3, c0.R, c1.R, c2.R, c3.R)
			c.G = product4(a0, a1, a2, a3, c0.G, c1.G, c2.G, c3.G)
			c.B = product4(a0, a1, a2, a3, c0.B, c1.B, c2.B, c3.B)
			c.A = product4(a0, a1, a2, a3, c0.A, c1.A, c2.A, c3.A)
			dst.SetNRGBAh(x+dstBounds.Min.X, y+dstBounds.Min.Y, c)
		}
	})
}

// Lanczos3 resizes the image using Lanczos interpolation with lobe 3.
func Lanczos3(dst, src *fp16.NRGBAh) {
	dstBounds := dst.Bounds()
	dstDx := dstBounds.Dx()
	dstDy := dstBounds.Dy()
	srcBounds := src.Bounds()
	srcDx := srcBounds.Dx()
	srcDy := srcBounds.Dy()

	tmp := fp16.NewNRGBAh(image.Rect(0, 0, dstDx, srcDy))

	// pre-calculate lanczos weights
	cache := make(map[float64]float64, 6*(dstDx+dstDy))
	for x := 0; x < dstDx; x++ {
		_, dx := scale(x, srcDx, dstDx)
		cache[dx+2] = lanczos(dx+2, 3)
		cache[dx+1] = lanczos(dx+1, 3)
		cache[dx+0] = lanczos(dx+0, 3)
		cache[dx-1] = lanczos(dx-1, 3)
		cache[dx-2] = lanczos(dx-2, 3)
		cache[dx-3] = lanczos(dx-3, 3)
	}
	for y := 0; y < dstDy; y++ {
		_, dy := scale(y, srcDy, dstDy)
		cache[dy+2] = lanczos(dy+2, 3)
		cache[dy+1] = lanczos(dy+1, 3)
		cache[dy+0] = lanczos(dy+0, 3)
		cache[dy-1] = lanczos(dy-1, 3)
		cache[dy-2] = lanczos(dy-2, 3)
		cache[dy-3] = lanczos(dy-3, 3)
	}

	// resize horizontally
	parallels.Parallel(0, srcDy, func(y int) {
		for x := 0; x < dstDx; x++ {
			var c fp16color.NRGBAh
			srcX, dx := scale(x, srcDx, dstDx)
			c0 := nrgbhAt(src, srcBounds.Min.X+srcX-2, srcBounds.Min.Y+y)
			c1 := nrgbhAt(src, srcBounds.Min.X+srcX-1, srcBounds.Min.Y+y)
			c2 := nrgbhAt(src, srcBounds.Min.X+srcX+0, srcBounds.Min.Y+y)
			c3 := nrgbhAt(src, srcBounds.Min.X+srcX+1, srcBounds.Min.Y+y)
			c4 := nrgbhAt(src, srcBounds.Min.X+srcX+2, srcBounds.Min.Y+y)
			c5 := nrgbhAt(src, srcBounds.Min.X+srcX+3, srcBounds.Min.Y+y)

			a0 := cache[dx+2]
			a1 := cache[dx+1]
			a2 := cache[dx+0]
			a3 := cache[dx-1]
			a4 := cache[dx-2]
			a5 := cache[dx-3]

			c.R = product6(a0, a1, a2, a3, a4, a5, c0.R, c1.R, c2.R, c3.R, c4.R, c5.R)
			c.G = product6(a0, a1, a2, a3, a4, a5, c0.G, c1.G, c2.G, c3.G, c4.G, c5.G)
			c.B = product6(a0, a1, a2, a3, a4, a5, c0.B, c1.B, c2.B, c3.B, c4.B, c5.B)
			c.A = product6(a0, a1, a2, a3, a4, a5, c0.A, c1.A, c2.A, c3.A, c4.A, c5.A)
			tmp.SetNRGBAh(x, y, c)
		}
	})

	// resize vertically
	parallels.Parallel(0, dstDy, func(y int) {
		for x := 0; x < dstDx; x++ {
			var c fp16color.NRGBAh
			srcY, dy := scale(y, srcDy, dstDy)
			c0 := nrgbhAt(tmp, x, srcY-2)
			c1 := nrgbhAt(tmp, x, srcY-1)
			c2 := nrgbhAt(tmp, x, srcY+0)
			c3 := nrgbhAt(tmp, x, srcY+1)
			c4 := nrgbhAt(tmp, x, srcY+2)
			c5 := nrgbhAt(tmp, x, srcY+3)

			a0 := cache[dy+2]
			a1 := cache[dy+1]
			a2 := cache[dy+0]
			a3 := cache[dy-1]
			a4 := cache[dy-2]
			a5 := cache[dy-3]

			c.R = product6(a0, a1, a2, a3, a4, a5, c0.R, c1.R, c2.R, c3.R, c4.R, c5.R)
			c.G = product6(a0, a1, a2, a3, a4, a5, c0.G, c1.G, c2.G, c3.G, c4.G, c5.G)
			c.B = product6(a0, a1, a2, a3, a4, a5, c0.B, c1.B, c2.B, c3.B, c4.B, c5.B)
			c.A = product6(a0, a1, a2, a3, a4, a5, c0.A, c1.A, c2.A, c3.A, c4.A, c5.A)
			dst.SetNRGBAh(x+dstBounds.Min.X, y+dstBounds.Min.Y, c)
		}
	})
}

// Lanczos4 resizes the image using Lanczos interpolation with lobe 4.
func Lanczos4(dst, src *fp16.NRGBAh) {
	dstBounds := dst.Bounds()
	dstDx := dstBounds.Dx()
	dstDy := dstBounds.Dy()
	srcBounds := src.Bounds()
	srcDx := srcBounds.Dx()
	srcDy := srcBounds.Dy()

	tmp := fp16.NewNRGBAh(image.Rect(0, 0, dstDx, srcDy))

	// pre-calculate lanczos weights
	cache := make(map[float64]float64, 8*(dstDx+dstDy))
	for x := 0; x < dstDx; x++ {
		_, dx := scale(x, srcDx, dstDx)
		cache[dx+3] = lanczos(dx+3, 4)
		cache[dx+2] = lanczos(dx+2, 4)
		cache[dx+1] = lanczos(dx+1, 4)
		cache[dx+0] = lanczos(dx+0, 4)
		cache[dx-1] = lanczos(dx-1, 4)
		cache[dx-2] = lanczos(dx-2, 4)
		cache[dx-3] = lanczos(dx-3, 4)
		cache[dx-4] = lanczos(dx-4, 4)
	}
	for y := 0; y < dstDy; y++ {
		_, dy := scale(y, srcDy, dstDy)
		cache[dy+3] = lanczos(dy+3, 4)
		cache[dy+2] = lanczos(dy+2, 4)
		cache[dy+1] = lanczos(dy+1, 4)
		cache[dy+0] = lanczos(dy+0, 4)
		cache[dy-1] = lanczos(dy-1, 4)
		cache[dy-2] = lanczos(dy-2, 4)
		cache[dy-3] = lanczos(dy-3, 4)
		cache[dy-4] = lanczos(dy-4, 4)
	}

	// resize horizontally
	parallels.Parallel(0, srcDy, func(y int) {
		for x := 0; x < dstDx; x++ {
			var c fp16color.NRGBAh
			srcX, dx := scale(x, srcDx, dstDx)
			c0 := nrgbhAt(src, srcBounds.Min.X+srcX-3, srcBounds.Min.Y+y)
			c1 := nrgbhAt(src, srcBounds.Min.X+srcX-2, srcBounds.Min.Y+y)
			c2 := nrgbhAt(src, srcBounds.Min.X+srcX-1, srcBounds.Min.Y+y)
			c3 := nrgbhAt(src, srcBounds.Min.X+srcX+0, srcBounds.Min.Y+y)
			c4 := nrgbhAt(src, srcBounds.Min.X+srcX+1, srcBounds.Min.Y+y)
			c5 := nrgbhAt(src, srcBounds.Min.X+srcX+2, srcBounds.Min.Y+y)
			c6 := nrgbhAt(src, srcBounds.Min.X+srcX+3, srcBounds.Min.Y+y)
			c7 := nrgbhAt(src, srcBounds.Min.X+srcX+4, srcBounds.Min.Y+y)

			a0 := cache[dx+3]
			a1 := cache[dx+2]
			a2 := cache[dx+1]
			a3 := cache[dx+0]
			a4 := cache[dx-1]
			a5 := cache[dx-2]
			a6 := cache[dx-3]
			a7 := cache[dx-4]

			c.R = product8(a0, a1, a2, a3, a4, a5, a6, a7, c0.R, c1.R, c2.R, c3.R, c4.R, c5.R, c6.R, c7.R)
			c.G = product8(a0, a1, a2, a3, a4, a5, a6, a7, c0.G, c1.G, c2.G, c3.G, c4.G, c5.G, c6.G, c7.G)
			c.B = product8(a0, a1, a2, a3, a4, a5, a6, a7, c0.B, c1.B, c2.B, c3.B, c4.B, c5.B, c6.B, c7.B)
			c.A = product8(a0, a1, a2, a3, a4, a5, a6, a7, c0.A, c1.A, c2.A, c3.A, c4.A, c5.A, c6.A, c7.A)
			tmp.SetNRGBAh(x, y, c)
		}
	})

	// resize vertically
	parallels.Parallel(0, dstDy, func(y int) {
		for x := 0; x < dstDx; x++ {
			var c fp16color.NRGBAh
			srcY, dy := scale(y, srcDy, dstDy)
			c0 := nrgbhAt(tmp, x, srcY-3)
			c1 := nrgbhAt(tmp, x, srcY-2)
			c2 := nrgbhAt(tmp, x, srcY-1)
			c3 := nrgbhAt(tmp, x, srcY+0)
			c4 := nrgbhAt(tmp, x, srcY+1)
			c5 := nrgbhAt(tmp, x, srcY+2)
			c6 := nrgbhAt(tmp, x, srcY+3)
			c7 := nrgbhAt(tmp, x, srcY+4)

			a0 := cache[dy+3]
			a1 := cache[dy+2]
			a2 := cache[dy+1]
			a3 := cache[dy+0]
			a4 := cache[dy-1]
			a5 := cache[dy-2]
			a6 := cache[dy-3]
			a7 := cache[dy-4]

			c.R = product8(a0, a1, a2, a3, a4, a5, a6, a7, c0.R, c1.R, c2.R, c3.R, c4.R, c5.R, c6.R, c7.R)
			c.G = product8(a0, a1, a2, a3, a4, a5, a6, a7, c0.G, c1.G, c2.G, c3.G, c4.G, c5.G, c6.G, c7.G)
			c.B = product8(a0, a1, a2, a3, a4, a5, a6, a7, c0.B, c1.B, c2.B, c3.B, c4.B, c5.B, c6.B, c7.B)
			c.A = product8(a0, a1, a2, a3, a4, a5, a6, a7, c0.A, c1.A, c2.A, c3.A, c4.A, c5.A, c6.A, c7.A)
			dst.SetNRGBAh(x+dstBounds.Min.X, y+dstBounds.Min.Y, c)
		}
	})
}

func sinc(x float64) float64 {
	if x == 0 {
		return 1
	}
	x *= math.Pi
	return math.Sin(x) / x
}

func lanczos(x, lobe float64) float64 {
	x = math.Abs(x)
	if x < lobe {
		return sinc(x) * sinc(x/lobe)
	}
	return 0
}
