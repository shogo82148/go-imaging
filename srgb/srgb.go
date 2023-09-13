//go:generate go run internal/gen/main.go

// Package srgb handles [sRGB] colors.
//
// [sRGB]: https://en.wikipedia.org/wiki/SRGB
package srgb

import (
	"image"
	"image/color"

	"github.com/shogo82148/float16"
	"github.com/shogo82148/go-imaging/fp16"
	"github.com/shogo82148/go-imaging/fp16/fp16color"
	"github.com/shogo82148/go-imaging/internal/parallels"
)

// Decode decodes an sRGB color encoded image to a linear color image.
func Decode(img image.Image) *fp16.NRGBAh {
	return decode(img)
}

func decode(img image.Image) *fp16.NRGBAh {
	bounds := img.Bounds()
	ret := fp16.NewNRGBAh(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, a := c.RGBA()
			fa := float64(a) / 0xffff
			ret.SetNRGBAhAt(x, y, fp16color.NRGBAh{
				R: encodedToLinearTable16[r],
				G: encodedToLinearTable16[g],
				B: encodedToLinearTable16[b],
				A: float16.FromFloat64(fa),
			})
		}
	})
	return ret
}

// Encode encodes a linear color image to an sRGB color encoded image.
func Encode(img *fp16.NRGBAh) *image.NRGBA64 {
	bounds := img.Bounds()
	ret := image.NewNRGBA64(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba := img.NRGBAhAt(x, y)
			a := uint16(rgba.A.Float64() * 0xffff)
			ret.SetRGBA64(x, y, color.RGBA64{
				R: linearToEncodedTable16[rgba.R],
				G: linearToEncodedTable16[rgba.G],
				B: linearToEncodedTable16[rgba.B],
				A: a,
			})
		}
	})

	return ret
}
