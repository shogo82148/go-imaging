package resize

import (
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
