package main

import (
	"bytes"
	_ "embed"
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
	resize2x(filepath.Join(dir, "autumn_bilinear_2x.png"), autumn, resize.BiLinear)
	resize2x(filepath.Join(dir, "plush_bilinear_2x.png"), plush, resize.BiLinear)
	resize2x(filepath.Join(dir, "autumn_nearest_neighbor_2x.png"), autumn, resize.NearestNeighbor)
	resize2x(filepath.Join(dir, "plush_nearest_neighbor_2x.png"), plush, resize.NearestNeighbor)
	resize2x(filepath.Join(dir, "autumn_hermite_2x.png"), autumn, resize.Hermite)
	resize2x(filepath.Join(dir, "plush_hermite_2x.png"), plush, resize.Hermite)
	resize2x(filepath.Join(dir, "autumn_general_2x.png"), autumn, resize.General)
	resize2x(filepath.Join(dir, "plush_general_2x.png"), plush, resize.General)
	resize2x(filepath.Join(dir, "autumn_catmull_rom_2x.png"), autumn, resize.CatmullRom)
	resize2x(filepath.Join(dir, "plush_catmull_rom_2x.png"), plush, resize.CatmullRom)
	resize2x(filepath.Join(dir, "autumn_mitchell_2x.png"), autumn, resize.Mitchell)
	resize2x(filepath.Join(dir, "plush_mitchell_rom_2x.png"), plush, resize.Mitchell)

	resize0_5x(filepath.Join(dir, "autumn_bilinear_0.5x.png"), autumn, resize.BiLinear)
	resize0_5x(filepath.Join(dir, "plush_bilinear_0.5x.png"), plush, resize.BiLinear)
	resize0_5x(filepath.Join(dir, "autumn_nearest_neighbor_0.5x.png"), autumn, resize.NearestNeighbor)
	resize0_5x(filepath.Join(dir, "plush_nearest_neighbor_0.5x.png"), plush, resize.NearestNeighbor)
	resize0_5x(filepath.Join(dir, "autumn_hermite_0.5x.png"), autumn, resize.Hermite)
	resize0_5x(filepath.Join(dir, "plush_hermite_0.5x.png"), plush, resize.Hermite)
	resize0_5x(filepath.Join(dir, "autumn_catmull_rom_0.5x.png"), autumn, resize.CatmullRom)
	resize0_5x(filepath.Join(dir, "plush_catmull_rom_0.5x.png"), plush, resize.CatmullRom)
	resize0_5x(filepath.Join(dir, "autumn_mitchell_0.5x.png"), autumn, resize.Mitchell)
	resize0_5x(filepath.Join(dir, "plush_mitchell_0.5x.png"), plush, resize.Mitchell)
}

func resize2x(filename string, img []byte, f func(dst *fp16.NRGBAh, src *fp16.NRGBAh)) error {
	orig, _, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		return err
	}
	src := srgb.Decode(orig)
	dst := fp16.NewNRGBAh(image.Rect(0, 0, orig.Bounds().Dx()*2, orig.Bounds().Dy()*2))
	f(dst, src)

	buf := new(bytes.Buffer)
	out := srgb.Encode(dst)
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
	src := srgb.Decode(orig)
	dst := fp16.NewNRGBAh(image.Rect(0, 0, orig.Bounds().Dx()/2, orig.Bounds().Dy()/2))
	f(dst, src)

	buf := new(bytes.Buffer)
	out := srgb.Encode(dst)
	if err := png.Encode(buf, out); err != nil {
		return err
	}
	return os.WriteFile(filename, buf.Bytes(), 0644)
}
