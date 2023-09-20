package resize

import (
	"math"

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

	parallels.Parallel(0, dstDy, func(y int) {
		for x := 0; x < dstDx; x++ {
			srcX, dx := scale(x, srcDx, dstDx)
			srcY, dy := scale(y, srcDy, dstDy)
			srcX += int(math.Round(dx))
			srcY += int(math.Round(dy))
			dst.SetNRGBAh(
				x+dstBounds.Min.X,
				y+dstBounds.Min.Y,
				nrgbhAt(src, srcBounds.Min.X+srcX, srcBounds.Min.Y+srcY),
			)
		}
	})
}
