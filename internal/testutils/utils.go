package testutils

import (
	"image"
	"image/draw"
	"strings"
	"testing"

	"github.com/shogo82148/go-imaging/pnm"
)

// LoadPNM loads a PNM image from string and
// convert it to *image.NRGBA64.
func LoadPNM(data string) *image.NRGBA64 {
	r := strings.NewReader(data)
	img, err := pnm.Decode(r)
	if err != nil {
		panic(err)
	}
	return ToNRGBA64(img)
}

// ToNRGBA64 converts an image.Image to *image.NRGBA64.
func ToNRGBA64(img image.Image) *image.NRGBA64 {
	dst := image.NewNRGBA64(img.Bounds())
	draw.Draw(dst, dst.Bounds(), img, image.Point{}, draw.Src)
	return dst
}

var plainEncoder = &pnm.Encoder{Plain: true}

func Compare(t *testing.T, got, want *image.NRGBA64) {
	t.Helper()

	if got.Bounds() != want.Bounds() {
		t.Errorf("bounds mismatch: got %v, want %v", got.Bounds(), want.Bounds())
		return
	}

	// compare pixel by pixel
	match := true
	for y := got.Bounds().Min.Y; y < got.Bounds().Max.Y; y++ {
		for x := got.Bounds().Min.X; x < got.Bounds().Max.X; x++ {
			match = match && (got.RGBA64At(x, y) == want.RGBA64At(x, y))
		}
	}
	if !match {
		var buf1 strings.Builder
		plainEncoder.Encode(&buf1, got)
		var buf2 strings.Builder
		plainEncoder.Encode(&buf2, want)
		t.Errorf("pixel mismatch: got:\n%s\nwant:\n%s", buf1.String(), buf2.String())
	}
}
