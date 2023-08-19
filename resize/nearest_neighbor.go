package resize

import (
	"image"
	"math/bits"
)

// NearestNeighbor resizes the image using nearest neighbor interpolation.
func NearestNeighbor(dst, src *image.NRGBA64) {
	dstBounds := dst.Bounds()
	dstDx := dstBounds.Dx()
	dstDy := dstBounds.Dy()
	srcBounds := src.Bounds()
	srcDx := srcBounds.Dx()
	srcDy := srcBounds.Dy()

	for y := 0; y < dstDy; y++ {
		for x := 0; x < dstDx; x++ {
			dst.SetNRGBA64(x+dstBounds.Min.X, y+dstBounds.Min.Y, src.NRGBA64At(
				mulDiv(x, srcDx-1, dstDx-1),
				mulDiv(y, srcDy-1, dstDy-1),
			))
		}
	}
}

func mulDiv(a, b, c int) int {
	hi, lo := bits.Mul64(uint64(a), uint64(b))

	// round to nearest
	lo, carry := bits.Add64(lo, uint64(c)/2, 0)
	hi += carry

	quo, _ := bits.Div64(hi, lo, uint64(c))
	return int(quo)
}
