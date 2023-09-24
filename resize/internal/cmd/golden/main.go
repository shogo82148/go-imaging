// generate golden files

package main

import (
	"bytes"
	"compress/gzip"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/shogo82148/go-imaging/fp16"
	"github.com/shogo82148/go-imaging/resize"
	"github.com/shogo82148/go-imaging/resize/internal/golden"
	"github.com/shogo82148/go-imaging/srgb"
)

func main() {
	dir := os.Args[1]

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

	src := golden.InputPattern()
	for _, f := range funcs {
		buf := new(bytes.Buffer)

		dst := golden.NewDst()
		f.f(dst, src)
		w := gzip.NewWriter(buf)
		w.Write(dst.Pix)
		w.Close()
		err := os.WriteFile(filepath.Join(dir, f.name+".golden.gz"), buf.Bytes(), 0644)
		if err != nil {
			log.Println(err)
		}

		buf.Reset()
		out := srgb.NonLinearize(dst)
		if err := png.Encode(buf, out); err != nil {
			log.Println(err)
		}
		if err := os.WriteFile(filepath.Join(dir, f.name+".png"), buf.Bytes(), 0644); err != nil {
			log.Println(err)
		}
	}
}
