package pnm

import (
	"bytes"
	"image"
	"testing"

	"github.com/shogo82148/go-imaging/bitmap"
	"github.com/shogo82148/go-imaging/graymap"
)

func TestEncode(t *testing.T) {
	plainEncoder := &Encoder{Plain: true}

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
		if err := plainEncoder.Encode(buf, img); err != nil {
			t.Fatal(err)
		}
		want := `P1
6 10
000010
000010
000010
000010
000010
000010
100010
011100
000000
000000
`
		got := buf.String()
		if got != want {
			t.Errorf("unexpected output: got %v, want %v", got, want)
		}
	})

	t.Run("encode plain PGM", func(t *testing.T) {
		img := &graymap.Image{
			Pix: []uint8{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 3, 3, 3, 3, 0, 0, 7, 7, 7, 7, 0, 0, 11, 11, 11, 11, 0, 0, 15, 15, 15, 15, 0,
				0, 3, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 11, 0, 0, 0, 0, 0, 15, 0, 0, 15, 0,
				0, 3, 3, 3, 0, 0, 0, 7, 7, 7, 0, 0, 0, 11, 11, 11, 0, 0, 0, 15, 15, 15, 15, 0,
				0, 3, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 11, 0, 0, 0, 0, 0, 15, 0, 0, 0, 0,
				0, 3, 0, 0, 0, 0, 0, 7, 7, 7, 7, 0, 0, 11, 11, 11, 11, 0, 0, 15, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			},
			Stride: 24,
			Rect:   image.Rect(0, 0, 24, 7),
			Max:    15,
		}
		buf := &bytes.Buffer{}
		if err := plainEncoder.Encode(buf, img); err != nil {
			t.Fatal(err)
		}
		want := `P2
24 7
15
0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
0 3 3 3 3 0 0 7 7 7 7 0 0 11 11 11 11 0 0 15 15 15 15 0
0 3 0 0 0 0 0 7 0 0 0 0 0 11 0 0 0 0 0 15 0 0 15 0
0 3 3 3 0 0 0 7 7 7 0 0 0 11 11 11 0 0 0 15 15 15 15 0
0 3 0 0 0 0 0 7 0 0 0 0 0 11 0 0 0 0 0 15 0 0 0 0
0 3 0 0 0 0 0 7 7 7 7 0 0 11 11 11 11 0 0 15 0 0 0 0
0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
`
		got := buf.String()
		if got != want {
			t.Errorf("unexpected output: got %v, want %v", got, want)
		}
	})
}
