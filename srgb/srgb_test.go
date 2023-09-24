package srgb

import (
	"bytes"
	"compress/gzip"
	"image"
	"io"
	"os"
	"testing"

	"github.com/shogo82148/go-imaging/srgb/internal/golden"
)

func TestLinearize(t *testing.T) {
	t.Run("RGBA", func(t *testing.T) {
		img := golden.NewRGBA()
		output := Linearize(img)
		golden.Assert(t, "rgba", output)
	})

	t.Run("RGBA64", func(t *testing.T) {
		img := golden.NewRGBA64()
		output := Linearize(img)
		golden.Assert(t, "rgba64", output)
	})

	t.Run("NRGBA", func(t *testing.T) {
		img := golden.NewNRGBA()
		output := Linearize(img)
		golden.Assert(t, "nrgba", output)
	})

	t.Run("NRGBA64", func(t *testing.T) {
		img := golden.NewNRGBA64()
		output := Linearize(img)
		golden.Assert(t, "nrgba64", output)
	})

	t.Run("Gray", func(t *testing.T) {
		img := golden.NewGray()
		output := Linearize(img)
		golden.Assert(t, "gray", output)
	})

	t.Run("Gray16", func(t *testing.T) {
		img := golden.NewGray16()
		output := Linearize(img)
		golden.Assert(t, "gray16", output)
	})

	t.Run("Alpha", func(t *testing.T) {
		img := golden.NewAlpha()
		output := Linearize(img)
		golden.Assert(t, "alpha", output)
	})

	t.Run("Alpha16", func(t *testing.T) {
		img := golden.NewAlpha16()
		output := Linearize(img)
		golden.Assert(t, "alpha16", output)
	})

	t.Run("Paletted", func(t *testing.T) {
		img := golden.NewPaletted()
		output := Linearize(img)
		golden.Assert(t, "paletted", output)
	})
}

func TestNonLinearize(t *testing.T) {
	input := golden.NewNRGBAh()
	output := NonLinearize(input)

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

func BenchmarkLinearize_RGBA(b *testing.B) {
	input := image.NewRGBA(image.Rect(0, 0, 512, 512))
	for i := 0; i < b.N; i++ {
		Linearize(input)
	}
}

func BenchmarkLinearize_RGBA64(b *testing.B) {
	input := image.NewRGBA64(image.Rect(0, 0, 512, 512))
	for i := 0; i < b.N; i++ {
		Linearize(input)
	}
}

func BenchmarkLinearize_NRGBA(b *testing.B) {
	input := image.NewNRGBA(image.Rect(0, 0, 512, 512))
	for i := 0; i < b.N; i++ {
		Linearize(input)
	}
}

func BenchmarkLinearize_NRGBA64(b *testing.B) {
	input := image.NewNRGBA64(image.Rect(0, 0, 512, 512))
	for i := 0; i < b.N; i++ {
		Linearize(input)
	}
}
