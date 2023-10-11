package exif_test

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/float16"
	. "github.com/shogo82148/go-imaging/exif"
	"github.com/shogo82148/go-imaging/fp16"
	"github.com/shogo82148/go-imaging/fp16/fp16color"
	"github.com/shogo82148/go-imaging/jpeg"
	"github.com/shogo82148/go-imaging/pnm"
	"github.com/shogo82148/go-imaging/srgb"
)

func TestAutoOrientation(t *testing.T) {
	src := fp16.NewNRGBAh(image.Rect(0, 0, 2, 2))
	one := float16.FromFloat64(1.0)
	c00 := fp16color.NRGBAh{A: one}
	c01 := fp16color.NRGBAh{R: one, A: one}
	c10 := fp16color.NRGBAh{G: one, A: one}
	c11 := fp16color.NRGBAh{B: one, A: one}
	src.SetNRGBAh(0, 0, c00)
	src.SetNRGBAh(0, 1, c01)
	src.SetNRGBAh(1, 0, c10)
	src.SetNRGBAh(1, 1, c11)

	t.Run("OrientationTopLeft", func(t *testing.T) {
		want := src
		dst := AutoOrientation(OrientationTopLeft, src)
		if diff := cmp.Diff(dst, want); diff != "" {
			t.Errorf("AutoOrientation() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("OrientationTopRight", func(t *testing.T) {
		want := fp16.NewNRGBAh(image.Rect(0, 0, 2, 2))
		want.SetNRGBAh(0, 0, c10)
		want.SetNRGBAh(0, 1, c11)
		want.SetNRGBAh(1, 0, c00)
		want.SetNRGBAh(1, 1, c01)
		dst := AutoOrientation(OrientationTopRight, src)
		if diff := cmp.Diff(dst, want); diff != "" {
			t.Errorf("AutoOrientation() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("OrientationBottomRight", func(t *testing.T) {
		want := fp16.NewNRGBAh(image.Rect(0, 0, 2, 2))
		want.SetNRGBAh(0, 0, c11)
		want.SetNRGBAh(0, 1, c10)
		want.SetNRGBAh(1, 0, c01)
		want.SetNRGBAh(1, 1, c00)
		dst := AutoOrientation(OrientationBottomRight, src)
		if diff := cmp.Diff(dst, want); diff != "" {
			t.Errorf("AutoOrientation() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("OrientationBottomLeft", func(t *testing.T) {
		want := fp16.NewNRGBAh(image.Rect(0, 0, 2, 2))
		want.SetNRGBAh(0, 0, c01)
		want.SetNRGBAh(0, 1, c00)
		want.SetNRGBAh(1, 0, c11)
		want.SetNRGBAh(1, 1, c10)
		dst := AutoOrientation(OrientationBottomLeft, src)
		if diff := cmp.Diff(dst, want); diff != "" {
			t.Errorf("AutoOrientation() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("OrientationLeftTop", func(t *testing.T) {
		want := fp16.NewNRGBAh(image.Rect(0, 0, 2, 2))
		want.SetNRGBAh(0, 0, c00)
		want.SetNRGBAh(0, 1, c10)
		want.SetNRGBAh(1, 0, c01)
		want.SetNRGBAh(1, 1, c11)
		dst := AutoOrientation(OrientationLeftTop, src)
		if diff := cmp.Diff(dst, want); diff != "" {
			t.Errorf("AutoOrientation() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("OrientationRightTop", func(t *testing.T) {
		want := fp16.NewNRGBAh(image.Rect(0, 0, 2, 2))
		want.SetNRGBAh(0, 0, c01)
		want.SetNRGBAh(0, 1, c11)
		want.SetNRGBAh(1, 0, c00)
		want.SetNRGBAh(1, 1, c10)
		dst := AutoOrientation(OrientationRightTop, src)
		if diff := cmp.Diff(dst, want); diff != "" {
			t.Errorf("AutoOrientation() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("OrientationRightBottom", func(t *testing.T) {
		want := fp16.NewNRGBAh(image.Rect(0, 0, 2, 2))
		want.SetNRGBAh(0, 0, c11)
		want.SetNRGBAh(0, 1, c01)
		want.SetNRGBAh(1, 0, c10)
		want.SetNRGBAh(1, 1, c00)
		dst := AutoOrientation(OrientationRightBottom, src)
		if diff := cmp.Diff(dst, want); diff != "" {
			t.Errorf("AutoOrientation() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("OrientationLeftBottom", func(t *testing.T) {
		want := fp16.NewNRGBAh(image.Rect(0, 0, 2, 2))
		want.SetNRGBAh(0, 0, c10)
		want.SetNRGBAh(0, 1, c00)
		want.SetNRGBAh(1, 0, c11)
		want.SetNRGBAh(1, 1, c01)
		dst := AutoOrientation(OrientationLeftBottom, src)
		if diff := cmp.Diff(dst, want); diff != "" {
			t.Errorf("AutoOrientation() mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestAutoOrientationJPEG(t *testing.T) {
	for i := 1; i <= 8; i++ {
		i := i
		t.Run(Orientation(i).String(), func(t *testing.T) {
			// decode the image
			f, err := os.Open(fmt.Sprintf("testdata/a-%d.jpg", i))
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			img, err := jpeg.DecodeWithMeta(f)
			if err != nil {
				t.Fatal(err)
			}

			// adjust the orientation
			o := img.Exif.Orientation
			src := srgb.DecodeTone(img.Image)
			dst := AutoOrientation(o, src)

			var got bytes.Buffer
			if err := pnm.Encode(&got, srgb.EncodeTone(dst)); err != nil {
				t.Fatal(err)
			}

			want, err := os.ReadFile(fmt.Sprintf("testdata/a-%d.golden.ppm", i))
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got.Bytes(), want); diff != "" {
				t.Errorf("AutoOrientation() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func BenchmarkAutoOrientation(b *testing.B) {
	for i := Orientation(1); i <= 8; i++ {
		b.Run(i.String(), func(b *testing.B) {
			src := fp16.NewNRGBAh(image.Rect(0, 0, 8192, 8192))
			AutoOrientation(i, src)
		})
	}
}
