package pnm

import (
	"bytes"
	"image"
	"testing"

	"github.com/shogo82148/go-imaging/bitmap"
	"github.com/shogo82148/go-imaging/graymap"
	"github.com/shogo82148/go-imaging/pixmap"
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

	t.Run("encode plain PPM", func(t *testing.T) {
		img := &pixmap.Image{
			Pix: []uint8{
				255, 0, 0,
				0, 255, 0,
				0, 0, 255,
				255, 255, 0,
				255, 255, 255,
				0, 0, 0,
			},
			Stride: 9,
			Rect:   image.Rect(0, 0, 3, 2),
			Max:    255,
		}
		buf := &bytes.Buffer{}
		if err := plainEncoder.Encode(buf, img); err != nil {
			t.Fatal(err)
		}
		want := `P3
3 2
255
255 0 0 0 255 0 0 0 255
255 255 0 255 255 255 0 0 0
`
		got := buf.String()
		if got != want {
			t.Errorf("unexpected output: got %v, want %v", got, want)
		}
	})

	t.Run("encode raw PBM", func(t *testing.T) {
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
		if err := defaultEncoder.Encode(buf, img); err != nil {
			t.Fatal(err)
		}
		want := []byte{
			'P', '4', '\n',
			'6', ' ', '1', '0', '\n',
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
		}
		got := buf.Bytes()
		if !bytes.Equal(got, want) {
			t.Errorf("unexpected output: got %v, want %v", got, want)
		}
	})

	t.Run("encode raw PGM", func(t *testing.T) {
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
		if err := defaultEncoder.Encode(buf, img); err != nil {
			t.Fatal(err)
		}
		want := []byte{
			'P', '5', '\n',
			'2', '4', ' ', '7', '\n',
			'1', '5', '\n',
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 3, 3, 3, 3, 0, 0, 7, 7, 7, 7, 0, 0, 11, 11, 11, 11, 0, 0, 15, 15, 15, 15, 0,
			0, 3, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 11, 0, 0, 0, 0, 0, 15, 0, 0, 15, 0,
			0, 3, 3, 3, 0, 0, 0, 7, 7, 7, 0, 0, 0, 11, 11, 11, 0, 0, 0, 15, 15, 15, 15, 0,
			0, 3, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 11, 0, 0, 0, 0, 0, 15, 0, 0, 0, 0,
			0, 3, 0, 0, 0, 0, 0, 7, 7, 7, 7, 0, 0, 11, 11, 11, 11, 0, 0, 15, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		}
		got := buf.Bytes()
		if !bytes.Equal(got, want) {
			t.Errorf("unexpected output: got %v, want %v", got, want)
		}
	})
}
