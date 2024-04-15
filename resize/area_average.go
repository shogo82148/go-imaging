package resize

import (
	"math"

	"github.com/shogo82148/float16"
	"github.com/shogo82148/go-imaging/fp16"
	"github.com/shogo82148/go-imaging/fp16/fp16color"
	"github.com/shogo82148/go-imaging/internal/parallels"
)

// AreaAverage resizes the image using area average interpolation.
func AreaAverage(dst, src *fp16.NRGBAh) {
	dstBounds := dst.Bounds()
	dstDx := dstBounds.Dx()
	dstDy := dstBounds.Dy()
	srcBounds := src.Bounds()
	srcDx := srcBounds.Dx()
	srcDy := srcBounds.Dy()

	parallels.Parallel(0, dstDy, func(y int) {
		for x := 0; x < dstDx; x++ {
			var rr, gg, bb, aa, ww float64
			srcX0, dx0 := scaleArea(x, srcDx, dstDx)
			srcY0, dy0 := scaleArea(y, srcDy, dstDy)
			srcX1, dx1 := scaleArea(x+1, srcDx, dstDx)
			srcY1, dy1 := scaleArea(y+1, srcDy, dstDy)

			// top right corner
			c := nrgbhAt(src, srcBounds.Min.X+srcX0, srcBounds.Min.Y+srcY0)
			w := (1 - dx0) * (1 - dy0)
			rr = math.FMA(c.R.Float64(), w, rr)
			gg = math.FMA(c.G.Float64(), w, gg)
			bb = math.FMA(c.B.Float64(), w, bb)
			aa = math.FMA(c.A.Float64(), w, aa)
			ww += w

			// top left corner
			c = nrgbhAt(src, srcBounds.Min.X+srcX1, srcBounds.Min.Y+srcY0)
			w = dx1 * (1 - dy0)
			rr = math.FMA(c.R.Float64(), w, rr)
			gg = math.FMA(c.G.Float64(), w, gg)
			bb = math.FMA(c.B.Float64(), w, bb)
			aa = math.FMA(c.A.Float64(), w, aa)
			ww += w

			// bottom right corner
			c = nrgbhAt(src, srcBounds.Min.X+srcX0, srcBounds.Min.Y+srcY1)
			w = (1 - dx0) * dy1
			rr = math.FMA(c.R.Float64(), w, rr)
			gg = math.FMA(c.G.Float64(), w, gg)
			bb = math.FMA(c.B.Float64(), w, bb)
			aa = math.FMA(c.A.Float64(), w, aa)
			ww += w

			// bottom left corner
			c = nrgbhAt(src, srcBounds.Min.X+srcX1, srcBounds.Min.Y+srcY1)
			w = dx1 * dy1
			rr = math.FMA(c.R.Float64(), w, rr)
			gg = math.FMA(c.G.Float64(), w, gg)
			bb = math.FMA(c.B.Float64(), w, bb)
			aa = math.FMA(c.A.Float64(), w, aa)
			ww += w

			for srcX := srcX0 + 1; srcX < srcX1; srcX++ {
				// top edge
				c := nrgbhAt(src, srcBounds.Min.X+srcX, srcBounds.Min.Y+srcY0)
				w := 1 - dy0
				rr = math.FMA(c.R.Float64(), w, rr)
				gg = math.FMA(c.G.Float64(), w, gg)
				bb = math.FMA(c.B.Float64(), w, bb)
				aa = math.FMA(c.A.Float64(), w, aa)
				ww += w

				// bottom edge
				c = nrgbhAt(src, srcBounds.Min.X+srcX, srcBounds.Min.Y+srcY1)
				w = dy1
				rr = math.FMA(c.R.Float64(), w, rr)
				gg = math.FMA(c.G.Float64(), w, gg)
				bb = math.FMA(c.B.Float64(), w, bb)
				aa = math.FMA(c.A.Float64(), w, aa)
				ww += w
			}

			for srcY := srcY0 + 1; srcY < srcY1; srcY++ {
				// left edge
				c := nrgbhAt(src, srcBounds.Min.X+srcX0, srcBounds.Min.Y+srcY)
				w := 1 - dx0
				rr = math.FMA(c.R.Float64(), w, rr)
				gg = math.FMA(c.G.Float64(), w, gg)
				bb = math.FMA(c.B.Float64(), w, bb)
				aa = math.FMA(c.A.Float64(), w, aa)
				ww += w

				// right edge
				c = nrgbhAt(src, srcBounds.Min.X+srcX1, srcBounds.Min.Y+srcY)
				w = dx1
				rr = math.FMA(c.R.Float64(), w, rr)
				gg = math.FMA(c.G.Float64(), w, gg)
				bb = math.FMA(c.B.Float64(), w, bb)
				aa = math.FMA(c.A.Float64(), w, aa)
				ww += w
			}

			for srcY := srcY0 + 1; srcY < srcY1; srcY++ {
				for srcX := srcX0 + 1; srcX < srcX1; srcX++ {
					c := nrgbhAt(src, srcBounds.Min.X+srcX, srcBounds.Min.Y+srcY)
					rr += c.R.Float64()
					gg += c.G.Float64()
					bb += c.B.Float64()
					aa += c.A.Float64()
					ww += 1.0
				}
			}

			var cc fp16color.NRGBAh
			cc.R = float16.FromFloat64(rr / ww)
			cc.G = float16.FromFloat64(gg / ww)
			cc.B = float16.FromFloat64(bb / ww)
			cc.A = float16.FromFloat64(aa / ww)
			dst.SetNRGBAh(x+dstBounds.Min.X, y+dstBounds.Min.Y, cc)
		}
	})
}

// scaleArea returns the source point and the distance from the source point.
func scaleArea(x, srcDx, dstDx int) (srcX int, dx float64) {
	quo, rem := mulDiv(x, srcDx, 0, dstDx)
	srcX = quo
	dx = float64(rem) / float64(dstDx)
	return
}
