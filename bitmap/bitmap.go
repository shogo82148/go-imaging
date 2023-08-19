// Package bitmap support binary image format.
package bitmap

import (
	"image"
	"image/color"
	"image/draw"
)

// Color is binary color.
type Color bool

var _ color.Color = White

const (
	// White is white color.
	White Color = false

	// Black is black color.
	Black Color = true
)

// RGBA implements [image/color.Color].
func (c Color) RGBA() (r, g, b, a uint32) {
	if c {
		return 0, 0, 0, 0xffff
	}
	return 0xffff, 0xffff, 0xffff, 0xffff
}

var ColorModel color.Model = color.ModelFunc(binaryModel)

func binaryModel(c color.Color) color.Color {
	bin, ok := c.(Color)
	if ok {
		return bin
	}

	r, g, b, _ := c.RGBA()

	// These coefficients (the fractions 0.299, 0.587 and 0.114) are the same
	// as those given by the JFIF specification and used by func RGBToYCbCr in
	// ycbcr.go.
	//
	// Note that 19595 + 38470 + 7471 equals 65536.
	y := (19595*r + 38470*g + 7471*b + 1<<15) >> 16

	return Color(y < 0x8000)
}

var _ image.Image = (*Image)(nil)
var _ draw.Image = (*Image)(nil)

// Image is a binary image.
type Image struct {
	// Pix holds the image's pixels, in binary format.
	// Each bit represents a pixel: 1 is black, 0 is white.
	// The order of the pixels is left to right.
	// The pixel at (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*2].
	Pix []uint8

	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int

	// Rect is the image's bounds.
	Rect image.Rectangle
}

// New returns a new binary image with the bounds r.
func New(r image.Rectangle) *Image {
	stride := (r.Dx() + 7) / 8
	return &Image{
		Pix:    make([]uint8, r.Dy()*stride),
		Stride: stride,
		Rect:   r,
	}
}

// ColorModel implements [image.Image].
func (img *Image) ColorModel() color.Model {
	return ColorModel
}

// Bounds implements [image.Image].
func (img *Image) Bounds() image.Rectangle {
	return img.Rect
}

// At implements [image.Image].
func (img *Image) At(x, y int) color.Color {
	return img.BinaryAt(x, y)
}

// BinaryAt returns the color of the pixel at (x, y).
func (img *Image) BinaryAt(x, y int) Color {
	if !(image.Point{x, y}).In(img.Rect) {
		return White
	}
	offset := (y-img.Rect.Min.Y)*img.Stride + (x-img.Rect.Min.X)/8
	shift := 7 - (x-img.Rect.Min.X)%8
	return Color((img.Pix[offset]>>shift)&0x01 != 0)
}

// Set implements [draw.Image].
func (img *Image) Set(x, y int, c color.Color) {
	img.SetBinary(x, y, binaryModel(c).(Color))
}

// SetBinary sets the color of the pixel at (x, y) to c.
func (img *Image) SetBinary(x, y int, c Color) {
	if !(image.Point{x, y}).In(img.Rect) {
		return
	}
	offset := (y-img.Rect.Min.Y)*img.Stride + (x-img.Rect.Min.X)/8
	shift := (x - img.Rect.Min.X) % 8
	mask := byte(0x80 >> shift)
	if c {
		img.Pix[offset] |= mask
	} else {
		img.Pix[offset] &^= mask
	}
}
