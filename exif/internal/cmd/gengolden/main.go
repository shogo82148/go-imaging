package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shogo82148/go-imaging/exif"
	"github.com/shogo82148/go-imaging/jpeg"
	"github.com/shogo82148/go-imaging/pnm"
	"github.com/shogo82148/go-imaging/srgb"
)

func main() {
	dir := os.Args[1]
	for i := 1; i <= 8; i++ {
		if err := generateGolden(dir, i); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func generateGolden(dir string, i int) error {
	name := fmt.Sprintf("a-%d.jpg", i)
	f, err := os.Open(filepath.Join(dir, name))
	if err != nil {
		return err
	}
	defer f.Close()
	img, err := jpeg.DecodeWithMeta(f)
	if err != nil {
		return err
	}

	// adjust the orientation
	o := img.Exif.Orientation
	src := srgb.DecodeTone(img.Image)
	dst := exif.AutoOrientation(o, src)
	got := srgb.EncodeTone(dst)

	f, err = os.Create(filepath.Join(dir, fmt.Sprintf("a-%d.golden.ppm", i)))
	if err != nil {
		return err
	}
	defer f.Close()

	if err := pnm.Encode(f, got); err != nil {
		return err
	}

	return f.Close()
}
