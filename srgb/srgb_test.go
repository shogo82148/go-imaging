package srgb

import (
	"bytes"
	"compress/gzip"
	"image"
	"image/color"
	"image/color/palette"
	"image/png"
	"io"
	"os"
	"testing"

	"github.com/shogo82148/go-imaging/fp16"
	"github.com/shogo82148/go-imaging/srgb/internal/golden"
)

func TestDecodeTone(t *testing.T) {
	t.Run("RGBA", func(t *testing.T) {
		img := golden.NewRGBA()
		output := DecodeTone(img)
		golden.Assert(t, "rgba", output)
	})

	t.Run("RGBA64", func(t *testing.T) {
		img := golden.NewRGBA64()
		output := DecodeTone(img)
		golden.Assert(t, "rgba64", output)
	})

	t.Run("NRGBA", func(t *testing.T) {
		img := golden.NewNRGBA()
		output := DecodeTone(img)
		golden.Assert(t, "nrgba", output)
	})

	t.Run("NRGBA64", func(t *testing.T) {
		img := golden.NewNRGBA64()
		output := DecodeTone(img)
		golden.Assert(t, "nrgba64", output)
	})

	t.Run("Gray", func(t *testing.T) {
		img := golden.NewGray()
		output := DecodeTone(img)
		golden.Assert(t, "gray", output)
	})

	t.Run("Gray16", func(t *testing.T) {
		img := golden.NewGray16()
		output := DecodeTone(img)
		golden.Assert(t, "gray16", output)
	})

	t.Run("Alpha", func(t *testing.T) {
		img := golden.NewAlpha()
		output := DecodeTone(img)
		golden.Assert(t, "alpha", output)
	})

	t.Run("Alpha16", func(t *testing.T) {
		img := golden.NewAlpha16()
		output := DecodeTone(img)
		golden.Assert(t, "alpha16", output)
	})

	t.Run("Paletted", func(t *testing.T) {
		img := golden.NewPaletted()
		output := DecodeTone(img)
		golden.Assert(t, "paletted", output)
	})
}

func TestEncodeTone(t *testing.T) {
	input := golden.NewNRGBAh()
	output := EncodeTone(input)

	f, err := os.Open("testdata/non-linearize.golden.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	r, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	golden, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(golden, output.Pix) {
		t.Errorf("mismatch")
		return
	}
}

func readPNG(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return png.Decode(f)
}

func compareImage(t *testing.T, got, want image.Image) {
	if got.Bounds() != want.Bounds() {
		t.Errorf("bounds mismatch: got %v, want %v", got.Bounds(), want.Bounds())
		return
	}

	for y := got.Bounds().Min.Y; y < got.Bounds().Max.Y; y++ {
		for x := got.Bounds().Min.X; x < got.Bounds().Max.X; x++ {
			c0 := color.NRGBA64Model.Convert(got.At(x, y)).(color.NRGBA64)
			c1 := color.NRGBA64Model.Convert(want.At(x, y)).(color.NRGBA64)
			if c0 != c1 {
				t.Errorf("color mismatch at (%d, %d): got %v, want %v", x, y, c0, c1)
				return
			}
		}
	}
}

func TestProfile_Linearize(t *testing.T) {
	input, err := readPNG("../testdata/senkakuwan.png")
	if err != nil {
		t.Fatal(err)
	}

	want, err := readPNG("./testdata/senkakuwan.golden.png")
	if err != nil {
		t.Fatal(err)
	}

	got := DecodeTone(input)

	compareImage(t, got, want)
}

func BenchmarkDecodeTone_RGBA(b *testing.B) {
	input := image.NewRGBA(image.Rect(0, 0, 512, 512))
	for i := 0; i < b.N; i++ {
		DecodeTone(input)
	}
}

func BenchmarkDecodeTone_RGBA64(b *testing.B) {
	input := image.NewRGBA64(image.Rect(0, 0, 512, 512))
	for i := 0; i < b.N; i++ {
		DecodeTone(input)
	}
}

func BenchmarkDecodeTone_NRGBA(b *testing.B) {
	input := image.NewNRGBA(image.Rect(0, 0, 512, 512))
	for i := 0; i < b.N; i++ {
		DecodeTone(input)
	}
}

func BenchmarkDecodeTone_NRGBA64(b *testing.B) {
	input := image.NewNRGBA64(image.Rect(0, 0, 512, 512))
	for i := 0; i < b.N; i++ {
		DecodeTone(input)
	}
}

func BenchmarkDecodeTone_Gray(b *testing.B) {
	input := image.NewGray(image.Rect(0, 0, 512, 512))
	for i := 0; i < b.N; i++ {
		DecodeTone(input)
	}
}

func BenchmarkDecodeTone_Gray16(b *testing.B) {
	input := image.NewGray(image.Rect(0, 0, 512, 512))
	for i := 0; i < b.N; i++ {
		DecodeTone(input)
	}
}

func BenchmarkDecodeTone_Alpha(b *testing.B) {
	input := image.NewAlpha(image.Rect(0, 0, 512, 512))
	for i := 0; i < b.N; i++ {
		DecodeTone(input)
	}
}

func BenchmarkDecodeTone_Alpha16(b *testing.B) {
	input := image.NewAlpha16(image.Rect(0, 0, 512, 512))
	for i := 0; i < b.N; i++ {
		DecodeTone(input)
	}
}

func BenchmarkDecodeTone_Paletted(b *testing.B) {
	input := image.NewPaletted(image.Rect(0, 0, 512, 512), palette.Plan9)
	for i := 0; i < b.N; i++ {
		DecodeTone(input)
	}
}

func BenchmarkEncodeTone(b *testing.B) {
	input := fp16.NewNRGBAh(image.Rect(0, 0, 512, 512))
	for i := 0; i < b.N; i++ {
		EncodeTone(input)
	}
}
