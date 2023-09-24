package golden

import (
	"bytes"
	"compress/gzip"
	"image"
	"image/color"
	"image/color/palette"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/shogo82148/float16"
	"github.com/shogo82148/go-imaging/fp16"
	"github.com/shogo82148/go-imaging/fp16/fp16color"
)

func NewRGBA() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 256, 256))
	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			img.SetRGBA(x, y, color.RGBA{
				R: uint8(x),
				G: uint8(x),
				B: uint8(x),
				A: uint8(y),
			})
		}
	}
	return img
}

func NewRGBA64() image.Image {
	img := image.NewRGBA64(image.Rect(0, 0, 512, 512))
	for y := 0; y < 512; y++ {
		for x := 0; x < 512; x++ {
			v := uint16(x * 65535 / 511)
			a := uint16(y * 65535 / 511)
			img.SetRGBA64(x, y, color.RGBA64{
				R: v,
				G: v,
				B: v,
				A: a,
			})
		}
	}
	return img
}

func NewNRGBA() image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, 256, 256))
	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			img.SetNRGBA(x, y, color.NRGBA{
				R: uint8(x),
				G: uint8(x),
				B: uint8(x),
				A: uint8(y),
			})
		}
	}
	return img
}

func NewNRGBA64() image.Image {
	img := image.NewNRGBA64(image.Rect(0, 0, 512, 512))
	for y := 0; y < 512; y++ {
		for x := 0; x < 512; x++ {
			v := uint16(x * 65535 / 511)
			a := uint16(y * 65535 / 511)
			img.SetNRGBA64(x, y, color.NRGBA64{
				R: v,
				G: v,
				B: v,
				A: a,
			})
		}
	}
	return img
}

func NewGray() image.Image {
	img := image.NewGray(image.Rect(0, 0, 16, 16))
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			v := y*16 + x
			img.SetGray(x, y, color.Gray{uint8(v)})
		}
	}
	return img
}

func NewGray16() image.Image {
	img := image.NewGray16(image.Rect(0, 0, 256, 256))
	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			v := y*256 + x
			img.SetGray16(x, y, color.Gray16{uint16(v)})
		}
	}
	return img
}

func NewAlpha() image.Image {
	img := image.NewAlpha(image.Rect(0, 0, 16, 16))
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			v := y*16 + x
			img.SetAlpha(x, y, color.Alpha{uint8(v)})
		}
	}
	return img
}

func NewAlpha16() image.Image {
	img := image.NewAlpha16(image.Rect(0, 0, 256, 256))
	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			v := y*256 + x
			img.SetAlpha16(x, y, color.Alpha16{uint16(v)})
		}
	}
	return img
}

func NewPaletted() image.Image {
	img := image.NewPaletted(image.Rect(0, 0, 16, 16), palette.Plan9)
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			v := y*16 + x
			img.SetColorIndex(x, y, uint8(v))
		}
	}
	return img
}

func NewNRGBAh() *fp16.NRGBAh {
	img := fp16.NewNRGBAh(image.Rect(0, 0, 256, 256))
	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			v := float16.FromFloat64(float64(x) / 255)
			a := float16.FromFloat64(float64(y) / 255)
			img.SetNRGBAh(x, y, fp16color.NRGBAh{
				R: v,
				G: v,
				B: v,
				A: a,
			})
		}
	}
	return img
}

func Assert(t *testing.T, name string, got *fp16.NRGBAh) {
	t.Helper()
	f, err := os.Open(filepath.Join("testdata", name+".golden.gz"))
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()

	r, err := gzip.NewReader(f)
	if err != nil {
		t.Error(err)
		return
	}
	defer r.Close()

	golden, err := io.ReadAll(r)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(golden, got.Pix) {
		t.Errorf("%s: mismatch", name)
		return
	}
}
