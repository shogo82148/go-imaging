package exif

import (
	"image"

	"github.com/shogo82148/go-imaging/fp16"
)

// AutoOrientation returns an image with the orientation applied.
func AutoOrientation(o Orientation, src *fp16.NRGBAh) *fp16.NRGBAh {
	switch o {
	case OrientationTopLeft:
		return src
	case OrientationTopRight:
		bounds := src.Bounds()
		dy := bounds.Dy()
		dx := bounds.Dx()
		dst := fp16.NewNRGBAh(image.Rect(0, 0, dx, dy))
		for y := 0; y < dy; y++ {
			offset := y * dst.Stride
			for x := 0; x < dx; x++ {
				copy(dst.Pix[offset+(dx-1-x)*8:], src.Pix[offset+x*8+0:offset+x*8+8])
			}
		}
		return dst
	case OrientationBottomRight:
		bounds := src.Bounds()
		dy := bounds.Dy()
		dx := bounds.Dx()
		dst := fp16.NewNRGBAh(image.Rect(0, 0, dx, dy))
		for y := 0; y < dy; y++ {
			srcOffset := y * src.Stride
			dstOffset := (dy - y - 1) * dst.Stride
			for x := 0; x < dx; x++ {
				copy(dst.Pix[dstOffset+(dx-1-x)*8:], src.Pix[srcOffset+x*8+0:srcOffset+x*8+8])
			}
		}
		return dst
	case OrientationBottomLeft:
		dst := fp16.NewNRGBAh(src.Bounds())
		dy := dst.Bounds().Dy()
		for y := 0; y < dy; y++ {
			copy(dst.Pix[(dy-y-1)*dst.Stride:(dy-y)*dst.Stride], src.Pix[y*src.Stride:(y+1)*src.Stride])
		}
		return dst
	case OrientationLeftTop:
		bounds := src.Bounds()
		dy := bounds.Dy()
		dx := bounds.Dx()
		dst := fp16.NewNRGBAh(image.Rect(0, 0, dy, dx))
		for y := 0; y < dy; y++ {
			srcOffset := y * src.Stride
			for x := 0; x < dx; x++ {
				dstOffset := x * dst.Stride
				copy(dst.Pix[dstOffset+y*8:], src.Pix[srcOffset+x*8+0:srcOffset+x*8+8])
			}
		}
		return dst
	case OrientationRightTop:
		bounds := src.Bounds()
		dy := bounds.Dy()
		dx := bounds.Dx()
		dst := fp16.NewNRGBAh(image.Rect(0, 0, dy, dx))
		for y := 0; y < dy; y++ {
			srcOffset := y * src.Stride
			for x := 0; x < dx; x++ {
				dstOffset := x * dst.Stride
				copy(dst.Pix[dstOffset+(dy-1-y)*8:], src.Pix[srcOffset+x*8+0:srcOffset+x*8+8])
			}
		}
		return dst
	case OrientationRightBottom:
		bounds := src.Bounds()
		dy := bounds.Dy()
		dx := bounds.Dx()
		dst := fp16.NewNRGBAh(image.Rect(0, 0, dy, dx))
		for y := 0; y < dy; y++ {
			srcOffset := y * src.Stride
			for x := 0; x < dx; x++ {
				dstOffset := (dx - 1 - x) * dst.Stride
				copy(dst.Pix[dstOffset+(dy-1-y)*8:], src.Pix[srcOffset+x*8+0:srcOffset+x*8+8])
			}
		}
		return dst
	case OrientationLeftBottom:
		bounds := src.Bounds()
		dy := bounds.Dy()
		dx := bounds.Dx()
		dst := fp16.NewNRGBAh(image.Rect(0, 0, dy, dx))
		for y := 0; y < dy; y++ {
			srcOffset := y * src.Stride
			for x := 0; x < dx; x++ {
				dstOffset := (dx - 1 - x) * dst.Stride
				copy(dst.Pix[dstOffset+y*8:], src.Pix[srcOffset+x*8+0:srcOffset+x*8+8])
			}
		}
		return dst
	default:
		return src
	}
}
