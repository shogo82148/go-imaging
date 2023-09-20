package resize

import (
	"math"

	"github.com/shogo82148/go-imaging/fp16"
	"github.com/shogo82148/go-imaging/fp16/fp16color"
)

// scale returns the source point and the distance from the source point.
func scale(x, srcDx, dstDx int) (srcX int, dx float64) {
	fx := float64(x) + 0.5
	fx = (fx * float64(srcDx)) / float64(dstDx)
	fx -= 0.5
	dx = fx - math.Floor(fx)
	srcX = int(math.Floor(fx))
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
