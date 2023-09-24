package main

import (
	"compress/gzip"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/shogo82148/go-imaging/fp16"
	"github.com/shogo82148/go-imaging/srgb"
	"github.com/shogo82148/go-imaging/srgb/internal/golden"
)

func main() {
	dir := os.Args[1]

	generators := []struct {
		name  string
		input func() image.Image
	}{
		{"rgba", golden.NewRGBA},
		{"rgba64", golden.NewRGBA64},
		{"nrgba", golden.NewNRGBA},
		{"nrgba64", golden.NewNRGBA64},
		{"gray", golden.NewGray},
		{"gray16", golden.NewGray16},
		{"alpha", golden.NewAlpha},
		{"alpha16", golden.NewAlpha16},
		{"paletted", golden.NewPaletted},
	}

	for _, g := range generators {
		img := g.input()
		output := srgb.DecodeTone(img)
		saveGolden(filepath.Join(dir, g.name+".golden.gz"), output)
	}

	generateNonLinearizeGolden(dir)
}

func saveGolden(path string, img *fp16.NRGBAh) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := gzip.NewWriter(f)
	if _, err := w.Write(img.Pix); err != nil {
		log.Fatal(err)
	}
	if err := w.Close(); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func generateNonLinearizeGolden(dir string) {
	input := golden.NewNRGBAh()
	output := srgb.EncodeTone(input)

	// save the output as the raw data.
	f, err := os.OpenFile(filepath.Join(dir, "non-linearize.golden.gz"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := gzip.NewWriter(f)
	if _, err := w.Write(output.Pix); err != nil {
		log.Fatal(err)
	}
	if err := w.Close(); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	// save the output as the PNG image for the preview.
	f, err = os.OpenFile(filepath.Join(dir, "non-linearize.png"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, output); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
