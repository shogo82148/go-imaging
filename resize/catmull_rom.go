package resize

import (
	"image"

	"github.com/shogo82148/go-imaging/fp16"
	"github.com/shogo82148/go-imaging/fp16/fp16color"
	"github.com/shogo82148/go-imaging/internal/parallels"
)

// CatmullRom resizes the image using CatmullRom (Bicubic B:0, C:1/2) interpolation.
func CatmullRom(dst, src *fp16.NRGBAh) {
	dstBounds := dst.Bounds()
	dstDx := dstBounds.Dx()
	dstDy := dstBounds.Dy()
	srcBounds := src.Bounds()
	srcDx := srcBounds.Dx()
	srcDy := srcBounds.Dy()
	coeff := cubicBCcoefficient(0, 0.5)

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

			c.R = general(c0.R, c1.R, c2.R, c3.R, dx, coeff)
			c.G = general(c0.G, c1.G, c2.G, c3.G, dx, coeff)
			c.B = general(c0.B, c1.B, c2.B, c3.B, dx, coeff)
			c.A = general(c0.A, c1.A, c2.A, c3.A, dx, coeff)
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

			c.R = general(c0.R, c1.R, c2.R, c3.R, dy, coeff)
			c.G = general(c0.G, c1.G, c2.G, c3.G, dy, coeff)
			c.B = general(c0.B, c1.B, c2.B, c3.B, dy, coeff)
			c.A = general(c0.A, c1.A, c2.A, c3.A, dy, coeff)
			dst.SetNRGBAh(x+dstBounds.Min.X, y+dstBounds.Min.Y, c)
		}
	})
}
