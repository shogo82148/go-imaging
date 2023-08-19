package testutils

import (
	"image/color"
	"testing"
)

func TestLoadPNM(t *testing.T) {
	img := LoadPNM(`P3
2 2
65535
0 0 0 65535 65535 65535
32768 0 0 0 32768 0
`)
	if got, want := img.NRGBA64At(0, 0), (color.NRGBA64{0x0000, 0x0000, 0x0000, 0xffff}); got != want {
		t.Errorf("unexpected color: got %v, want %v", got, want)
	}
	if got, want := img.NRGBA64At(1, 0), (color.NRGBA64{0xffff, 0xffff, 0xffff, 0xffff}); got != want {
		t.Errorf("unexpected color: got %v, want %v", got, want)
	}
	if got, want := img.NRGBA64At(0, 1), (color.NRGBA64{0x8000, 0x0000, 0x0000, 0xffff}); got != want {
		t.Errorf("unexpected color: got %v, want %v", got, want)
	}
	if got, want := img.NRGBA64At(1, 1), (color.NRGBA64{0x0000, 0x8000, 0x0000, 0xffff}); got != want {
		t.Errorf("unexpected color: got %v, want %v", got, want)
	}
}
