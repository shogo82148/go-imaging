package srgb

import (
	"image"
	"testing"
)

func BenchmarkDecode_RGBA(b *testing.B) {
	input := image.NewRGBA(image.Rect(0, 0, 512, 512))
	for i := 0; i < b.N; i++ {
		Decode(input)
	}
}

func BenchmarkDecode_RGBA64(b *testing.B) {
	input := image.NewRGBA64(image.Rect(0, 0, 512, 512))
	for i := 0; i < b.N; i++ {
		Decode(input)
	}
}

func BenchmarkDecode_NRGBA(b *testing.B) {
	input := image.NewNRGBA(image.Rect(0, 0, 512, 512))
	for i := 0; i < b.N; i++ {
		Decode(input)
	}
}

func BenchmarkDecode_NRGBA64(b *testing.B) {
	input := image.NewNRGBA64(image.Rect(0, 0, 512, 512))
	for i := 0; i < b.N; i++ {
		Decode(input)
	}
}
