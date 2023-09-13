package resize

import (
	"math/bits"

	"github.com/shogo82148/go-imaging/fp16"
	"github.com/shogo82148/go-imaging/internal/parallels"
)

// NearestNeighbor resizes the image using nearest neighbor interpolation.
func NearestNeighbor(dst, src *fp16.NRGBAh) {
	dstBounds := dst.Bounds()
	dstDx := dstBounds.Dx()
	dstDy := dstBounds.Dy()
	srcBounds := src.Bounds()
	srcDx := srcBounds.Dx()
	srcDy := srcBounds.Dy()

	parallels.Parallel(dstBounds.Min.Y, dstBounds.Max.Y, func(y int) {
		for x := 0; x < dstDx; x++ {
			dst.SetNRGBAh(
				x+dstBounds.Min.X,
				y+dstBounds.Min.Y,
				src.NRGBAhAt(
					mulDivRound(x, srcDx-1, dstDx-1),
					mulDivRound(y, srcDy-1, dstDy-1),
				),
			)
		}
	})
}

func mulDivRound(a, b, c int) int {
	hi, lo := bits.Mul64(uint64(a), uint64(b))

	// round to nearest
	lo, carry := bits.Add64(lo, uint64(c)/2, 0)
	hi += carry

	quo, _ := bits.Div64(hi, lo, uint64(c))
	return int(quo)
}
