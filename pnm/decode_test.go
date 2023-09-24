package pnm

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/shogo82148/go-imaging/bitmap"
	"github.com/shogo82148/go-imaging/graymap"
	"github.com/shogo82148/go-imaging/pixmap"
)

func TestDecode(t *testing.T) {
	t.Run("ascii PBM Example from Wikipedia", func(t *testing.T) {
		// example from Wikipedia https://en.wikipedia.org/wiki/Netpbm#PBM_example
		f, err := os.Open("testdata/wikipedia_example_j.pbm")
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		img, err := Decode(f)
		if err != nil {
			t.Error(err)
		}
		want := []byte{
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
		if !bytes.Equal(img.(*bitmap.Image).Pix, want) {
			t.Errorf("expected %v, got %v", want, img.(*bitmap.Image).Pix)
		}
		if img.ColorModel() != bitmap.ColorModel {
			t.Errorf("expected bitmap.ColorModel, got %v", img.ColorModel())
		}
		if img.Bounds().Dx() != 6 {
			t.Errorf("expected width 6, got %d", img.Bounds().Dx())
		}
		if img.Bounds().Dy() != 10 {
			t.Errorf("expected height 10, got %d", img.Bounds().Dy())
		}
	})

	t.Run("another ascii PBM Example from Wikipedia", func(t *testing.T) {
		// example from Wikipedia https://en.wikipedia.org/wiki/Netpbm#PBM_example
		f, err := os.Open("testdata/wikipedia_example_j2.pbm")
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		img, err := Decode(f)
		if err != nil {
			t.Error(err)
		}
		want := []byte{
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
		if !bytes.Equal(img.(*bitmap.Image).Pix, want) {
			t.Errorf("expected %v, got %v", want, img.(*bitmap.Image).Pix)
		}
		if img.ColorModel() != bitmap.ColorModel {
			t.Errorf("expected bitmap.ColorModel, got %v", img.ColorModel())
		}
		if img.Bounds().Dx() != 6 {
			t.Errorf("expected width 6, got %d", img.Bounds().Dx())
		}
		if img.Bounds().Dy() != 10 {
			t.Errorf("expected height 10, got %d", img.Bounds().Dy())
		}
	})

	t.Run("ascii PBM", func(t *testing.T) {
		// example from https://netpbm.sourceforge.net/doc/pbm.html
		f, err := os.Open("testdata/feep-ascii.pbm")
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		img, err := Decode(f)
		if err != nil {
			t.Error(err)
		}
		want := []byte{
			0b00000000, 0b00000000, 0b00000000,
			0b01111001, 0b11100111, 0b10011110,
			0b01000001, 0b00000100, 0b00010010,
			0b01110001, 0b11000111, 0b00011110,
			0b01000001, 0b00000100, 0b00010000,
			0b01000001, 0b11100111, 0b10010000,
			0b00000000, 0b00000000, 0b00000000,
		}
		if !bytes.Equal(img.(*bitmap.Image).Pix, want) {
			t.Errorf("expected %v, got %v", want, img.(*bitmap.Image).Pix)
		}
		if img.ColorModel() != bitmap.ColorModel {
			t.Errorf("expected bitmap.ColorModel, got %v", img.ColorModel())
		}
		if img.Bounds().Dx() != 24 {
			t.Errorf("expected width 6, got %d", img.Bounds().Dx())
		}
		if img.Bounds().Dy() != 7 {
			t.Errorf("expected height 10, got %d", img.Bounds().Dy())
		}
	})

	t.Run("ascii PGM", func(t *testing.T) {
		// example from https://netpbm.sourceforge.net/doc/pgm.html
		f, err := os.Open("testdata/feep-ascii.pgm")
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		img, err := Decode(f)
		if err != nil {
			t.Error(err)
		}
		want := []byte{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 3, 3, 3, 3, 0, 0, 7, 7, 7, 7, 0, 0, 11, 11, 11, 11, 0, 0, 15, 15, 15, 15, 0,
			0, 3, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 11, 0, 0, 0, 0, 0, 15, 0, 0, 15, 0,
			0, 3, 3, 3, 0, 0, 0, 7, 7, 7, 0, 0, 0, 11, 11, 11, 0, 0, 0, 15, 15, 15, 15, 0,
			0, 3, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 11, 0, 0, 0, 0, 0, 15, 0, 0, 0, 0,
			0, 3, 0, 0, 0, 0, 0, 7, 7, 7, 7, 0, 0, 11, 11, 11, 11, 0, 0, 15, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		}
		if !bytes.Equal(img.(*graymap.Image).Pix, want) {
			t.Errorf("expected %v, got %v", want, img.(*graymap.Image).Pix)
		}
		if img.ColorModel() != graymap.Model(15) {
			t.Errorf("expected graymap.Model(15), got %v", img.ColorModel())
		}
		if img.Bounds().Dx() != 24 {
			t.Errorf("expected width 24, got %d", img.Bounds().Dx())
		}
		if img.Bounds().Dy() != 7 {
			t.Errorf("expected height 7, got %d", img.Bounds().Dy())
		}
	})

	t.Run("ascii PPM 1", func(t *testing.T) {
		// example from https://en.wikipedia.org/wiki/Netpbm#PPM_example
		f, err := os.Open("testdata/wikipedia_example_ppm1.ppm")
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		img, err := Decode(f)
		if err != nil {
			t.Error(err)
		}
		want := []byte{
			255, 0, 0,
			0, 255, 0,
			0, 0, 255,
			255, 255, 0,
			255, 255, 255,
			0, 0, 0,
		}
		if !bytes.Equal(img.(*pixmap.Image).Pix, want) {
			t.Errorf("expected %v, got %v", want, img.(*pixmap.Image).Pix)
		}
		if img.ColorModel() != pixmap.Model(255) {
			t.Errorf("expected pixmap.Model(255), got %v", img.ColorModel())
		}
		if img.Bounds().Dx() != 3 {
			t.Errorf("expected width 3, got %d", img.Bounds().Dx())
		}
		if img.Bounds().Dy() != 2 {
			t.Errorf("expected height 2, got %d", img.Bounds().Dy())
		}
	})

	t.Run("ascii PPM 2", func(t *testing.T) {
		// example from https://en.wikipedia.org/wiki/Netpbm#PPM_example
		f, err := os.Open("testdata/wikipedia_example_ppm2.ppm")
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		img, err := Decode(f)
		if err != nil {
			t.Error(err)
		}
		want := []byte{
			1, 0, 0, 0, 1, 0, 0, 0, 1,
			1, 1, 0, 1, 1, 1, 0, 0, 0,
		}
		if !bytes.Equal(img.(*pixmap.Image).Pix, want) {
			t.Errorf("expected %v, got %v", want, img.(*pixmap.Image).Pix)
		}
		if img.ColorModel() != pixmap.Model(1) {
			t.Errorf("expected pixmap.Model(1), got %v", img.ColorModel())
		}
		if img.Bounds().Dx() != 3 {
			t.Errorf("expected width 3, got %d", img.Bounds().Dx())
		}
		if img.Bounds().Dy() != 2 {
			t.Errorf("expected height 2, got %d", img.Bounds().Dy())
		}
	})

	t.Run("ascii PPM 3", func(t *testing.T) {
		// example from https://en.wikipedia.org/wiki/Netpbm#PPM_example
		f, err := os.Open("testdata/wikipedia_example_ppm3.ppm")
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		img, err := Decode(f)
		if err != nil {
			t.Error(err)
		}
		want := []byte{
			1, 0, 0, 0, 1, 0, 0, 0, 1,
			1, 1, 0, 1, 1, 1, 0, 0, 0,
		}
		if !bytes.Equal(img.(*pixmap.Image).Pix, want) {
			t.Errorf("expected %v, got %v", want, img.(*pixmap.Image).Pix)
		}
		if img.ColorModel() != pixmap.Model(1) {
			t.Errorf("expected pixmap.Model(1), got %v", img.ColorModel())
		}
		if img.Bounds().Dx() != 3 {
			t.Errorf("expected width 3, got %d", img.Bounds().Dx())
		}
		if img.Bounds().Dy() != 2 {
			t.Errorf("expected height 2, got %d", img.Bounds().Dy())
		}
	})

	t.Run("binary PBM maze", func(t *testing.T) {
		// netpbm test data https://sourceforge.net/p/netpbm/code/HEAD/tree/trunk/test/maze.pbm
		f, err := os.Open("testdata/maze.pbm")
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		img, err := Decode(f)
		if err != nil {
			t.Error(err)
		}
		if img.ColorModel() != bitmap.ColorModel {
			t.Errorf("expected bitmap.ColorModel, got %v", img.ColorModel())
		}
		if img.Bounds().Dx() != 57 {
			t.Errorf("expected width 6, got %d", img.Bounds().Dx())
		}
		if img.Bounds().Dy() != 59 {
			t.Errorf("expected height 10, got %d", img.Bounds().Dy())
		}
	})

	t.Run("binary PBM testgrid", func(t *testing.T) {
		// netpbm test data https://sourceforge.net/p/netpbm/code/HEAD/tree/trunk/test/maze.pbm
		f, err := os.Open("testdata/testgrid.pbm")
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		img, err := Decode(f)
		if err != nil {
			t.Error(err)
		}
		if img.ColorModel() != bitmap.ColorModel {
			t.Errorf("expected bitmap.ColorModel, got %v", img.ColorModel())
		}
		if img.Bounds().Dx() != 14 {
			t.Errorf("expected width 6, got %d", img.Bounds().Dx())
		}
		if img.Bounds().Dy() != 16 {
			t.Errorf("expected height 10, got %d", img.Bounds().Dy())
		}
	})

	t.Run("binary PGM", func(t *testing.T) {
		// converted feep-ascii.pgm by using GIMP.
		f, err := os.Open("testdata/feep-binary.pgm")
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		img, err := Decode(f)
		if err != nil {
			t.Error(err)
		}
		if img.ColorModel() != graymap.Model(255) {
			t.Errorf("expected graymap.Model(255), got %v", img.ColorModel())
		}
		if img.Bounds().Dx() != 24 {
			t.Errorf("expected width 24, got %d", img.Bounds().Dx())
		}
		if img.Bounds().Dy() != 7 {
			t.Errorf("expected height 7, got %d", img.Bounds().Dy())
		}
	})

	t.Run("binary PPM testimg", func(t *testing.T) {
		f, err := os.Open("testdata/testimg.ppm")
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		img, err := Decode(f)
		if err != nil {
			t.Error(err)
		}
		if img.ColorModel() != pixmap.Model(255) {
			t.Errorf("expected pixmap.Model(255), got %v", img.ColorModel())
		}
		if img.Bounds().Dx() != 227 {
			t.Errorf("expected width 227, got %d", img.Bounds().Dx())
		}
		if img.Bounds().Dy() != 149 {
			t.Errorf("expected height 149, got %d", img.Bounds().Dy())
		}
	})
}

func TestDecodeConfig(t *testing.T) {
	t.Run("PBM Example from Wikipedia", func(t *testing.T) {
		// example from Wikipedia https://en.wikipedia.org/wiki/Netpbm#PBM_example
		f, err := os.Open("testdata/wikipedia_example_j.pbm")
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		cfg, err := DecodeConfig(f)
		if err != nil {
			t.Error(err)
		}
		if cfg.ColorModel != bitmap.ColorModel {
			t.Errorf("expected bitmap.ColorModel, got %v", cfg.ColorModel)
		}
		if cfg.Width != 6 {
			t.Errorf("expected width 6, got %d", cfg.Width)
		}
		if cfg.Height != 10 {
			t.Errorf("expected height 10, got %d", cfg.Height)
		}
	})

	t.Run("maze", func(t *testing.T) {
		// netpbm test data https://sourceforge.net/p/netpbm/code/HEAD/tree/trunk/test/maze.pbm
		f, err := os.Open("testdata/maze.pbm")
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		cfg, err := DecodeConfig(f)
		if err != nil {
			t.Error(err)
		}
		if cfg.ColorModel != bitmap.ColorModel {
			t.Errorf("expected bitmap.ColorModel, got %v", cfg.ColorModel)
		}
		if cfg.Width != 57 {
			t.Errorf("expected width 6, got %d", cfg.Width)
		}
		if cfg.Height != 59 {
			t.Errorf("expected height 10, got %d", cfg.Height)
		}
	})
}

func TestInvalidFormat(t *testing.T) {
	t.Run("unknown magic number", func(t *testing.T) {
		input := "P7\n"
		r := strings.NewReader(input)
		_, err := Decode(r)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("empty magic number", func(t *testing.T) {
		input := ""
		r := strings.NewReader(input)
		_, err := Decode(r)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("incomplete magic number", func(t *testing.T) {
		input := "P"
		r := strings.NewReader(input)
		_, err := Decode(r)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("missing width", func(t *testing.T) {
		input := "P1\n"
		r := strings.NewReader(input)
		_, err := Decode(r)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("parsing width error", func(t *testing.T) {
		input := "P1\nNotANumber\n"
		r := strings.NewReader(input)
		_, err := Decode(r)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("0 width", func(t *testing.T) {
		input := "P1\n0\n"
		r := strings.NewReader(input)
		_, err := Decode(r)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("missing height", func(t *testing.T) {
		input := "P1\n1\n"
		r := strings.NewReader(input)
		_, err := Decode(r)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("parsing height error", func(t *testing.T) {
		input := "P1\n1\nNotANumber\n"
		r := strings.NewReader(input)
		_, err := Decode(r)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("0 height", func(t *testing.T) {
		input := "P1\n1\n0\n"
		r := strings.NewReader(input)
		_, err := Decode(r)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}
