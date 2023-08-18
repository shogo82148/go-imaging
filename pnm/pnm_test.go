package pnm

import (
	"os"
	"testing"

	"github.com/shogo82148/go-imaging/bitmap"
)

func TestDecode(t *testing.T) {
	t.Run("PBM Example from Wikipedia", func(t *testing.T) {
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
		_ = img // TODO: check the image
	})

	t.Run("maze", func(t *testing.T) {
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
