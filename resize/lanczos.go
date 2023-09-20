package resize

import (
	"image"
	"math"

	"github.com/shogo82148/float16"
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

	// resize horizontally
	parallels.Parallel(0, srcDy, func(y int) {
		for x := 0; x < dstDx; x++ {
			var c fp16color.NRGBAh
			srcX, dx := scale(x, srcDx, dstDx)
			c0 := nrgbhAt(src, srcBounds.Min.X+srcX-1, srcBounds.Min.Y+y)
			c1 := nrgbhAt(src, srcBounds.Min.X+srcX+0, srcBounds.Min.Y+y)
			c2 := nrgbhAt(src, srcBounds.Min.X+srcX+1, srcBounds.Min.Y+y)
			c3 := nrgbhAt(src, srcBounds.Min.X+srcX+2, srcBounds.Min.Y+y)

			c.R = lanczos2(c0.R, c1.R, c2.R, c3.R, dx)
			c.G = lanczos2(c0.G, c1.G, c2.G, c3.G, dx)
			c.B = lanczos2(c0.B, c1.B, c2.B, c3.B, dx)
			c.A = lanczos2(c0.A, c1.A, c2.A, c3.A, dx)
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

			c.R = lanczos2(c0.R, c1.R, c2.R, c3.R, dy)
			c.G = lanczos2(c0.G, c1.G, c2.G, c3.G, dy)
			c.B = lanczos2(c0.B, c1.B, c2.B, c3.B, dy)
			c.A = lanczos2(c0.A, c1.A, c2.A, c3.A, dy)
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

			c.R = lanczos3(c0.R, c1.R, c2.R, c3.R, c4.R, c5.R, dx)
			c.G = lanczos3(c0.G, c1.G, c2.G, c3.G, c4.G, c5.G, dx)
			c.B = lanczos3(c0.B, c1.B, c2.B, c3.B, c4.B, c5.B, dx)
			c.A = lanczos3(c0.A, c1.A, c2.A, c3.A, c4.A, c5.A, dx)
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

			c.R = lanczos3(c0.R, c1.R, c2.R, c3.R, c4.R, c5.R, dy)
			c.G = lanczos3(c0.G, c1.G, c2.G, c3.G, c4.G, c5.G, dy)
			c.B = lanczos3(c0.B, c1.B, c2.B, c3.B, c4.B, c5.B, dy)
			c.A = lanczos3(c0.A, c1.A, c2.A, c3.A, c4.A, c5.A, dy)
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

			c.R = lanczos4(c0.R, c1.R, c2.R, c3.R, c4.R, c5.R, c6.R, c7.R, dx)
			c.G = lanczos4(c0.G, c1.G, c2.G, c3.G, c4.G, c5.G, c6.G, c7.G, dx)
			c.B = lanczos4(c0.B, c1.B, c2.B, c3.B, c4.B, c5.B, c6.B, c7.B, dx)
			c.A = lanczos4(c0.A, c1.A, c2.A, c3.A, c4.A, c5.A, c6.A, c7.A, dx)
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

			c.R = lanczos4(c0.R, c1.R, c2.R, c3.R, c4.R, c5.R, c6.R, c7.R, dy)
			c.G = lanczos4(c0.G, c1.G, c2.G, c3.G, c4.G, c5.G, c6.G, c7.G, dy)
			c.B = lanczos4(c0.B, c1.B, c2.B, c3.B, c4.B, c5.B, c6.B, c7.B, dy)
			c.A = lanczos4(c0.A, c1.A, c2.A, c3.A, c4.A, c5.A, c6.A, c7.A, dy)
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

func lanczos2(c0, c1, c2, c3 float16.Float16, d float64) float16.Float16 {
	var c float64
	c = math.FMA(lanczos(d+1, 2), c0.Float64(), c)
	c = math.FMA(lanczos(d+0, 2), c1.Float64(), c)
	c = math.FMA(lanczos(d-1, 2), c2.Float64(), c)
	c = math.FMA(lanczos(d-2, 2), c3.Float64(), c)
	return float16.FromFloat64(c)
}

func lanczos3(c0, c1, c2, c3, c4, c5 float16.Float16, d float64) float16.Float16 {
	var c float64
	c = math.FMA(lanczos(d+2, 3), c0.Float64(), c)
	c = math.FMA(lanczos(d+1, 3), c1.Float64(), c)
	c = math.FMA(lanczos(d+0, 3), c2.Float64(), c)
	c = math.FMA(lanczos(d-1, 3), c3.Float64(), c)
	c = math.FMA(lanczos(d-2, 3), c4.Float64(), c)
	c = math.FMA(lanczos(d-3, 3), c5.Float64(), c)
	return float16.FromFloat64(c)
}

func lanczos4(c0, c1, c2, c3, c4, c5, c6, c7 float16.Float16, d float64) float16.Float16 {
	var c float64
	c = math.FMA(lanczos(d+3, 4), c0.Float64(), c)
	c = math.FMA(lanczos(d+2, 4), c1.Float64(), c)
	c = math.FMA(lanczos(d+1, 4), c2.Float64(), c)
	c = math.FMA(lanczos(d+0, 4), c3.Float64(), c)
	c = math.FMA(lanczos(d-1, 4), c4.Float64(), c)
	c = math.FMA(lanczos(d-2, 4), c5.Float64(), c)
	c = math.FMA(lanczos(d-3, 4), c6.Float64(), c)
	c = math.FMA(lanczos(d-4, 4), c7.Float64(), c)
	return float16.FromFloat64(c)
}
