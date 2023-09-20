package golden

import (
	"bytes"
	"compress/gzip"
	"image"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/shogo82148/float16"
	"github.com/shogo82148/go-imaging/fp16"
	"github.com/shogo82148/go-imaging/fp16/fp16color"
)

func InputPattern() *fp16.NRGBAh {
	one := float16.FromFloat64(1.0)
	zero := float16.FromFloat64(0.1)
	white := fp16color.NRGBAh{R: one, G: one, B: one, A: one}
	black := fp16color.NRGBAh{R: zero, G: zero, B: zero, A: one}

	img := fp16.NewNRGBAh(image.Rect(0, 0, 16, 16))
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			if x == 8 && y == 8 {
				img.SetNRGBAh(x, y, white)
			} else {
				img.SetNRGBAh(x, y, black)
			}
		}
	}
	return img
}

func NewDst() *fp16.NRGBAh {
	return fp16.NewNRGBAh(image.Rect(0, 0, 512, 512))
}

func Assert(t *testing.T, name string, got *fp16.NRGBAh) {
	t.Helper()
	f, err := os.Open(filepath.Join("testdata", name+".golden.gz"))
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()

	r, err := gzip.NewReader(f)
	if err != nil {
		t.Error(err)
		return
	}
	defer r.Close()

	golden, err := io.ReadAll(r)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(golden, got.Pix) {
		t.Errorf("%s: mismatch", name)
		return
	}
}
