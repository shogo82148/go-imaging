package srgb

import (
	"bytes"
	"compress/gzip"
	"image"
	"image/color/palette"
	"image/jpeg"
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

func TestProfile_Linearize(t *testing.T) {
	img, err := os.ReadFile("testdata/senkakuwan.jpg")
	if err != nil {
		t.Fatal(err)
	}
	m, err := jpeg.Decode(bytes.NewReader(img))
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.OpenFile("senkakuwan-linearized.jpg", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		t.Fatal(err)
	}
	out := DecodeTone(m)
	jpeg.Encode(f, out, nil)
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
