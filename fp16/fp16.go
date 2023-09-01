package fp16

import (
	"image"
	"image/color"
	"math/bits"

	"github.com/shogo82148/go-imaging/fp16/fp16color"
)

var _ image.Image = &NRGBAh{}

// NRGBA64 is an in-memory image whose At method returns fp16color.NRGBAh values.
type NRGBAh struct {
	// Pix holds the image's pixels, in R, G, B, A order and big-endian format. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*8].
	Pix []uint8

	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int

	// Rect is the image's bounds.
	Rect image.Rectangle
}

// NewNRGBh returns a new NRGBA64 image with the given bounds.
func NewNRGBAh(r image.Rectangle) *NRGBAh {
	return &NRGBAh{
		Pix:    make([]uint8, pixelBufferLength(8, r, "NRGBAh")),
		Stride: 8 * r.Dx(),
		Rect:   r,
	}
}

func (p *NRGBAh) ColorModel() color.Model {
	return fp16color.NRGBAhModel
}

func (p *NRGBAh) Bounds() image.Rectangle {
	return p.Rect
}

func (p *NRGBAh) At(x, y int) color.Color {
	return p.NRGBAhAt(x, y)
}

func (p *NRGBAh) NRGBAhAt(x, y int) fp16color.NRGBAh {
	return fp16color.NRGBAh{}
}

// pixelBufferLength returns the length of the []uint8 typed Pix slice field
// for the NewXxx functions. Conceptually, this is just (bpp * width * height),
// but this function panics if at least one of those is negative or if the
// computation would overflow the int type.
//
// This panics instead of returning an error because of backwards
// compatibility. The NewXxx functions do not return an error.
func pixelBufferLength(bytesPerPixel int, r image.Rectangle, imageTypeName string) int {
	// https://github.com/golang/go/blob/38b623f42da899ba7fd6b3fd791a7a72ebd5fad0/src/image/image.go#L86-L99

	totalLength := mul3NonNeg(bytesPerPixel, r.Dx(), r.Dy())
	if totalLength < 0 {
		panic("image: New" + imageTypeName + " Rectangle has huge or negative dimensions")
	}
	return totalLength
}

// mul3NonNeg returns (x * y * z), unless at least one argument is negative or
// if the computation overflows the int type, in which case it returns -1.
func mul3NonNeg(x int, y int, z int) int {
	// https://github.com/golang/go/blob/38b623f42da899ba7fd6b3fd791a7a72ebd5fad0/src/image/geom.go#L285-L304

	if (x < 0) || (y < 0) || (z < 0) {
		return -1
	}
	hi, lo := bits.Mul64(uint64(x), uint64(y))
	if hi != 0 {
		return -1
	}
	hi, lo = bits.Mul64(lo, uint64(z))
	if hi != 0 {
		return -1
	}
	a := int(lo)
	if (a < 0) || (uint64(a) != lo) {
		return -1
	}
	return a
}
