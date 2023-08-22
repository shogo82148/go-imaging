package resize

import (
	"image"
	"image/color"
	"math/bits"
)

func BiLinear(dst, src *image.NRGBA64) {
	dstBounds := dst.Bounds()
	dstDx := dstBounds.Dx()
	dstDy := dstBounds.Dy()
	srcBounds := src.Bounds()
	srcDx := srcBounds.Dx()
	srcDy := srcBounds.Dy()

	for y := 0; y < dstDy; y++ {
		for x := 0; x < dstDx; x++ {
			var c color.NRGBA64
			srcX, remX := mulDiv(x, srcDx-1, dstDx-1)
			srcY, remY := mulDiv(y, srcDy-1, dstDy-1)
			dx := float64(remX) / float64(dstDx-1)
			dy := float64(remY) / float64(dstDy-1)
			c0 := src.NRGBA64At(srcX, srcY)
			c1 := src.NRGBA64At(srcX+1, srcY)
			c2 := src.NRGBA64At(srcX, srcY+1)
			c3 := src.NRGBA64At(srcX+1, srcY+1)
			c.R = bilinear(c0.R, c1.R, c2.R, c3.R, dx, dy)
			c.G = bilinear(c0.G, c1.G, c2.G, c3.G, dx, dy)
			c.B = bilinear(c0.B, c1.B, c2.B, c3.B, dx, dy)
			c.A = bilinear(c0.A, c1.A, c2.A, c3.A, dx, dy)
			dst.SetNRGBA64(x+dstBounds.Min.X, y+dstBounds.Min.Y, c)
		}
	}
}

// mulDiv returns a * b / c.
func mulDiv(a, b, c int) (int, int) {
	hi, lo := bits.Mul64(uint64(a), uint64(b))
	quo, rem := bits.Div64(hi, lo, uint64(c))
	return int(quo), int(rem)
}

func bilinear(c0, c1, c2, c3 uint16, dx, dy float64) uint16 {
	return uint16(
		float64(c0)*(1-dx)*(1-dy) +
			float64(c1)*dx*(1-dy) +
			float64(c2)*(1-dx)*dy +
			float64(c3)*dx*dy +
			0.5,
	)
}
