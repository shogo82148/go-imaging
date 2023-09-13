package fp16

import (
	"bytes"
	"image"
	"testing"

	"github.com/shogo82148/go-imaging/fp16/fp16color"
)

func TestNRGBAhAt(t *testing.T) {
	img := NewNRGBAh(image.Rect(0, 0, 2, 2))
	img.SetNRGBAhAt(0, 0, fp16color.NewNRGBAh(0, 0, 0, 0))
	img.SetNRGBAhAt(1, 0, fp16color.NewNRGBAh(1, 1, 1, 1))
	img.SetNRGBAhAt(0, 1, fp16color.NewNRGBAh(0.5, 0.5, 0.5, 1))
	img.SetNRGBAhAt(1, 1, fp16color.NewNRGBAh(0.75, 0.25, 0.3, 1))

	got := img.Pix
	want := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x3c, 0x00, 0x3c, 0x00, 0x3c, 0x00, 0x3c, 0x00,
		0x38, 0x00, 0x38, 0x00, 0x38, 0x00, 0x3c, 0x00,
		0x3a, 0x00, 0x34, 0x00, 0x34, 0xcd, 0x3c, 0x00,
	}
	if !bytes.Equal(got, want) {
		t.Errorf("got %x, want %x", got, want)
	}
}
