package srgb

import (
	"image"
	"image/color"
	"testing"
)

func TestDecode(t *testing.T) {
	input := image.NewNRGBA64(image.Rect(0, 0, 10, 1))
	input.SetRGBA64(0, 0, color.RGBA64{0x0000, 0x0000, 0x000, 0x0000})
	input.SetRGBA64(1, 0, color.RGBA64{0xffff, 0x0000, 0x000, 0xffff})
	input.SetRGBA64(2, 0, color.RGBA64{0x8000, 0x0000, 0x000, 0xffff})
	output := Decode(input)

	if got, want := output.NRGBA64At(0, 0), (color.NRGBA64{0x0000, 0x0000, 0x000, 0x0000}); got != want {
		t.Errorf("unexpected color: got %v, want %v", got, want)
	}
	if got, want := output.NRGBA64At(1, 0), (color.NRGBA64{0xffff, 0x0000, 0x000, 0xffff}); got != want {
		t.Errorf("unexpected color: got %v, want %v", got, want)
	}
	if got, want := output.NRGBA64At(2, 0), (color.NRGBA64{0x36cc, 0x0000, 0x000, 0xffff}); got != want {
		t.Errorf("unexpected color: got %v, want %v", got, want)
	}
}

func BenchmarkDecode(b *testing.B) {
	input := image.NewNRGBA64(image.Rect(0, 0, 512, 512))
	for i := 0; i < b.N; i++ {
		Decode(input)
	}
}

func TestEncode(t *testing.T) {
	input := image.NewNRGBA64(image.Rect(0, 0, 10, 1))
	input.SetNRGBA64(0, 0, color.NRGBA64{0x0000, 0x0000, 0x000, 0x0000})
	input.SetNRGBA64(1, 0, color.NRGBA64{0xffff, 0x0000, 0x000, 0xffff})
	input.SetNRGBA64(2, 0, color.NRGBA64{0x36cc, 0x0000, 0x000, 0xffff})
	output := Encode(input)

	if got, want := output.NRGBA64At(0, 0), (color.NRGBA64{0x0000, 0x0000, 0x000, 0x0000}); got != want {
		t.Errorf("unexpected color: got %v, want %v", got, want)
	}
	if got, want := output.NRGBA64At(1, 0), (color.NRGBA64{0xffff, 0x0000, 0x000, 0xffff}); got != want {
		t.Errorf("unexpected color: got %v, want %v", got, want)
	}
	if got, want := output.NRGBA64At(2, 0), (color.NRGBA64{0x8000, 0x0000, 0x000, 0xffff}); got != want {
		t.Errorf("unexpected color: got %v, want %v", got, want)
	}
}

func BenchmarkEncode(b *testing.B) {
	input := image.NewNRGBA64(image.Rect(0, 0, 512, 512))
	for i := 0; i < b.N; i++ {
		Encode(input)
	}
}
