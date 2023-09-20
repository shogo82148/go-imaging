package resize

import (
	"math"

	"github.com/shogo82148/float16"
	"github.com/shogo82148/go-imaging/fp16"
	"github.com/shogo82148/go-imaging/fp16/fp16color"
)

// scale returns the source point and the distance from the source point.
func scale(x, srcDx, dstDx int) (srcX int, dx float64) {
	quo, rem := mulDiv(2*x+1, srcDx, dstDx, 2*dstDx)
	srcX = quo - 1
	dx = float64(rem) / float64(2*dstDx)
	return
}

func mulDiv(a, b, c, d int) (quo, rem int) {
	e := uint64(a)*uint64(b) + uint64(c)
	quo = int(e / uint64(d))
	rem = int(e % uint64(d))
	return
}

// nrgbhAt is similar to src.NRGBAhAt(x, y), but it returns the color at the
// nearest point if the point is out of bounds.
func nrgbhAt(img *fp16.NRGBAh, x, y int) fp16color.NRGBAh {
	bounds := img.Bounds()
	x = max(bounds.Min.X, min(bounds.Max.X-1, x))
	y = max(bounds.Min.Y, min(bounds.Max.Y-1, y))
	return img.NRGBAhAt(x, y)
}

// product4 calculates inner product of [a0, a1, a2, a3] and [c0, c1, c2, c3].
func product4(a0, a1, a2, a3 float64, c0, c1, c2, c3 float16.Float16) float16.Float16 {
	var c float64
	c = math.FMA(a0, c0.Float64(), c)
	c = math.FMA(a1, c1.Float64(), c)
	c = math.FMA(a2, c2.Float64(), c)
	c = math.FMA(a3, c3.Float64(), c)
	return float16.FromFloat64(c)
}

// product6 calculates inner product of [a0, a1, a2, a3, a4, a5] and [c0, c1, c2, c3, c4, c5].
func product6(a0, a1, a2, a3, a4, a5 float64, c0, c1, c2, c3, c4, c5 float16.Float16) float16.Float16 {
	var c float64
	c = math.FMA(a0, c0.Float64(), c)
	c = math.FMA(a1, c1.Float64(), c)
	c = math.FMA(a2, c2.Float64(), c)
	c = math.FMA(a3, c3.Float64(), c)
	c = math.FMA(a4, c4.Float64(), c)
	c = math.FMA(a5, c5.Float64(), c)
	return float16.FromFloat64(c)
}

// product8 calculates inner product of [a0, a1, a2, a3, a4, a5, a6, a7] and [c0, c1, c2, c3, c4, c5, c6, c7].
func product8(a0, a1, a2, a3, a4, a5, a6, a7 float64, c0, c1, c2, c3, c4, c5, c6, c7 float16.Float16) float16.Float16 {
	var c float64
	c = math.FMA(a0, c0.Float64(), c)
	c = math.FMA(a1, c1.Float64(), c)
	c = math.FMA(a2, c2.Float64(), c)
	c = math.FMA(a3, c3.Float64(), c)
	c = math.FMA(a4, c4.Float64(), c)
	c = math.FMA(a5, c5.Float64(), c)
	c = math.FMA(a6, c6.Float64(), c)
	c = math.FMA(a7, c7.Float64(), c)
	return float16.FromFloat64(c)
}
