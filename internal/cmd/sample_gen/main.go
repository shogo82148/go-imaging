package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/shogo82148/go-imaging/fp16"
	"github.com/shogo82148/go-imaging/resize"
	"github.com/shogo82148/go-imaging/srgb"
)

// https://github.com/tenntenn/gopher-stickers/blob/6fa428fa62c6a7d5ae2bb536636e59952a18b7ee/png/autumn.png
// The Go gopher was designed by Renee French (http://reneefrench.blogspot.com/).
// The gopher stickers was made by Takuya Ueda (https://twitter.com/tenntenn).
//
//go:embed autumn.png
var autumn []byte

// https://go.dev/blog/gopher
//
//go:embed plush.jpg
var plush []byte

func main() {
	dir := os.Args[1]
	grid := gridPattern()
	os.WriteFile(filepath.Join(dir, "grid.png"), grid, 0644)

	funcs := []struct {
		name string
		f    func(dst *fp16.NRGBAh, src *fp16.NRGBAh)
	}{
		{"bilinear", resize.BiLinear},
		{"nearest_neighbor", resize.NearestNeighbor},
		{"hermite", resize.Hermite},
		{"general", resize.General},
		{"catmull_rom", resize.CatmullRom},
		{"mitchell", resize.Mitchell},
		{"lanczos2", resize.Lanczos2},
		{"lanczos3", resize.Lanczos3},
		{"lanczos4", resize.Lanczos4},
	}

	images := []struct {
		name string
		img  []byte
	}{
		{"autumn", autumn},
		{"plush", plush},
		{"grid", grid},
	}

	for _, f := range funcs {
		for _, img := range images {
			var filename string
			filename = filepath.Join(dir, fmt.Sprintf("%s_%s_2x.png", img.name, f.name))
			resize2x(filename, img.img, f.f)
			filename = filepath.Join(dir, fmt.Sprintf("%s_%s_0.5x.png", img.name, f.name))
			resize0_5x(filename, img.img, f.f)
		}
	}
}

func resize2x(filename string, img []byte, f func(dst *fp16.NRGBAh, src *fp16.NRGBAh)) error {
	orig, _, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		return err
	}
	src := srgb.DecodeTone(orig)
	dst := fp16.NewNRGBAh(image.Rect(0, 0, orig.Bounds().Dx()*2, orig.Bounds().Dy()*2))
	f(dst, src)

	buf := new(bytes.Buffer)
	out := srgb.EncodeTone64(dst)
	if err := png.Encode(buf, out); err != nil {
		return err
	}
	return os.WriteFile(filename, buf.Bytes(), 0644)
}

func resize0_5x(filename string, img []byte, f func(dst *fp16.NRGBAh, src *fp16.NRGBAh)) error {
	orig, _, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		return err
	}
	src := srgb.DecodeTone(orig)
	dst := fp16.NewNRGBAh(image.Rect(0, 0, orig.Bounds().Dx()/2, orig.Bounds().Dy()/2))
	f(dst, src)

	buf := new(bytes.Buffer)
	out := srgb.EncodeTone64(dst)
	if err := png.Encode(buf, out); err != nil {
		return err
	}
	return os.WriteFile(filename, buf.Bytes(), 0644)
}

func gridPattern() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 256, 256))
	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			if (x+y)%2 == 0 {
				img.Set(x, y, image.White)
			} else {
				img.Set(x, y, image.Black)
			}
		}
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
