package resize

import (
	"math/bits"

	"github.com/shogo82148/float16"
	"github.com/shogo82148/go-imaging/fp16"
	"github.com/shogo82148/go-imaging/fp16/fp16color"
	"github.com/shogo82148/go-imaging/internal/parallels"
)

// BiLinear resizes the image using bilinear interpolation.
func BiLinear(dst, src *fp16.NRGBAh) {
	dstBounds := dst.Bounds()
	dstDx := dstBounds.Dx()
	dstDy := dstBounds.Dy()
	srcBounds := src.Bounds()
	srcDx := srcBounds.Dx()
	srcDy := srcBounds.Dy()

	parallels.Parallel(0, dstDy, func(y int) {
		for x := 0; x < dstDx; x++ {
			var c fp16color.NRGBAh
			srcX, remX := mulDiv(x, srcDx-1, dstDx-1)
			srcY, remY := mulDiv(y, srcDy-1, dstDy-1)
			dx := float64(remX) / float64(dstDx-1)
			dy := float64(remY) / float64(dstDy-1)
			c0 := src.NRGBAhAt(srcBounds.Min.X+srcX+0, srcBounds.Min.Y+srcY+0)
			c1 := src.NRGBAhAt(srcBounds.Min.X+srcX+1, srcBounds.Min.Y+srcY+0)
			c2 := src.NRGBAhAt(srcBounds.Min.X+srcX+0, srcBounds.Min.Y+srcY+1)
			c3 := src.NRGBAhAt(srcBounds.Min.X+srcX+1, srcBounds.Min.Y+srcY+1)
			c.R = bilinear(c0.R, c1.R, c2.R, c3.R, dx, dy)
			c.G = bilinear(c0.G, c1.G, c2.G, c3.G, dx, dy)
			c.B = bilinear(c0.B, c1.B, c2.B, c3.B, dx, dy)
			c.A = bilinear(c0.A, c1.A, c2.A, c3.A, dx, dy)
			dst.SetNRGBAh(x+dstBounds.Min.X, y+dstBounds.Min.Y, c)
		}
	})
}

// mulDiv returns a * b / c.
func mulDiv(a, b, c int) (int, int) {
	hi, lo := bits.Mul64(uint64(a), uint64(b))
	quo, rem := bits.Div64(hi, lo, uint64(c))
	return int(quo), int(rem)
}

func bilinear(c0, c1, c2, c3 float16.Float16, dx, dy float64) float16.Float16 {
	return float16.FromFloat64(c0.Float64()*(1-dx)*(1-dy) +
		c1.Float64()*dx*(1-dy) +
		c2.Float64()*(1-dx)*dy +
		c3.Float64()*dx*dy,
	)
}
