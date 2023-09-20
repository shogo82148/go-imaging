package resize

import (
	"github.com/shogo82148/float16"
	"github.com/shogo82148/go-imaging/fp16"
	"github.com/shogo82148/go-imaging/fp16/fp16color"
	"github.com/shogo82148/go-imaging/internal/parallels"
)

// Hermite resizes the image using Hermite (Bicubic B:0, C:0) interpolation.
func Hermite(dst, src *fp16.NRGBAh) {
	dstBounds := dst.Bounds()
	dstDx := dstBounds.Dx()
	dstDy := dstBounds.Dy()
	srcBounds := src.Bounds()
	srcDx := srcBounds.Dx()
	srcDy := srcBounds.Dy()

	parallels.Parallel(0, dstDy, func(y int) {
		for x := 0; x < dstDx; x++ {
			var c fp16color.NRGBAh
			srcX, dx0 := scale(x, srcDx, dstDx)
			srcY, dy0 := scale(y, srcDy, dstDy)
			c0 := nrgbhAt(src, srcBounds.Min.X+srcX, srcBounds.Min.Y+srcY)
			c1 := nrgbhAt(src, srcBounds.Min.X+srcX+1, srcBounds.Min.Y+srcY)
			c2 := nrgbhAt(src, srcBounds.Min.X+srcX, srcBounds.Min.Y+srcY+1)
			c3 := nrgbhAt(src, srcBounds.Min.X+srcX+1, srcBounds.Min.Y+srcY+1)

			// https://qiita.com/yoya/items/f167b2598fec98679422
			// https://legacy.imagemagick.org/Usage/filter/#mitchell
			// https://www.cs.utexas.edu/~fussell/courses/cs384g-fall2013/lectures/mitchell/Mitchell.pdf
			dx1 := 1 - dx0
			dy1 := 1 - dy0
			kx0 := dx0*dx0*(2*dx0-3) + 1
			kx1 := dx1*dx1*(2*dx1-3) + 1
			ky0 := dy0*dy0*(2*dy0-3) + 1
			ky1 := dy1*dy1*(2*dy1-3) + 1

			c.R = hermite(c0.R, c1.R, c2.R, c3.R, kx0, kx1, ky0, ky1)
			c.G = hermite(c0.G, c1.G, c2.G, c3.G, kx0, kx1, ky0, ky1)
			c.B = hermite(c0.B, c1.B, c2.B, c3.B, kx0, kx1, ky0, ky1)
			c.A = hermite(c0.A, c1.A, c2.A, c3.A, kx0, kx1, ky0, ky1)
			dst.SetNRGBAh(x+dstBounds.Min.X, y+dstBounds.Min.Y, c)
		}
	})
}

func hermite(c0, c1, c2, c3 float16.Float16, kx0, kx1, ky0, ky1 float64) float16.Float16 {
	c01 := c0.Float64()*kx0 + c1.Float64()*kx1
	c23 := c2.Float64()*kx0 + c3.Float64()*kx1
	return float16.FromFloat64(c01*ky0 + c23*ky1)
}
