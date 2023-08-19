package pnm

import (
	"bytes"
	"image"
	"testing"

	"github.com/shogo82148/go-imaging/bitmap"
)

func TestEncode(t *testing.T) {
	t.Run("encode plain PBM", func(t *testing.T) {
		img := &bitmap.Image{
			Pix: []byte{
				0b000010_00,
				0b000010_00,
				0b000010_00,
				0b000010_00,
				0b000010_00,
				0b000010_00,
				0b100010_00,
				0b011100_00,
				0b000000_00,
				0b000000_00,
			},
			Rect:   image.Rect(0, 0, 6, 10),
			Stride: 1,
		}
		buf := &bytes.Buffer{}
		if err := (&Encoder{Plain: true}).Encode(buf, img); err != nil {
			t.Fatal(err)
		}
		if got, want := buf.String(), "P1\n6 10\n000010\n000010\n000010\n000010\n000010\n000010\n100010\n011100\n000000\n000000\n"; got != want {
			t.Errorf("unexpected output: got %v, want %v", got, want)
		}
	})
}
