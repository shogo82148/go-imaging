//go:generate go run internal/gen/main.go

// Package srgb handles [sRGB] colors.
//
// [sRGB]: https://en.wikipedia.org/wiki/SRGB
package srgb

import (
	"image"
	"image/color"

	"github.com/shogo82148/go-imaging/internal/parallels"
)

// Decode decodes an sRGB color encoded image to a linear color image.
func Decode(img image.Image) *image.NRGBA64 {
	bounds := img.Bounds()
	ret := image.NewNRGBA64(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			c = color.NRGBA64Model.Convert(c)
			rgba := c.(color.NRGBA64)
			ret.SetRGBA64(x, y, color.RGBA64{
				R: encodedToLinearTable16[rgba.R],
				G: encodedToLinearTable16[rgba.G],
				B: encodedToLinearTable16[rgba.B],
				A: rgba.A,
			})
		}
	})

	return ret
}

// Encode encodes a linear color image to an sRGB color encoded image.
func Encode(img *image.NRGBA64) *image.NRGBA64 {
	bounds := img.Bounds()
	ret := image.NewNRGBA64(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba := img.NRGBA64At(x, y)
			ret.SetRGBA64(x, y, color.RGBA64{
				R: linearToEncodedTable16[rgba.R],
				G: linearToEncodedTable16[rgba.G],
				B: linearToEncodedTable16[rgba.B],
				A: rgba.A,
			})
		}
	})

	return ret
}
