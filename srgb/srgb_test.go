package srgb

import (
	"image"
	"testing"
)

func BenchmarkLinearize_RGBA(b *testing.B) {
	input := image.NewRGBA(image.Rect(0, 0, 512, 512))
	for i := 0; i < b.N; i++ {
		Linearize(input)
	}
}

func BenchmarkLinearize_RGBA64(b *testing.B) {
	input := image.NewRGBA64(image.Rect(0, 0, 512, 512))
	for i := 0; i < b.N; i++ {
		Linearize(input)
	}
}

func BenchmarkLinearize_NRGBA(b *testing.B) {
	input := image.NewNRGBA(image.Rect(0, 0, 512, 512))
	for i := 0; i < b.N; i++ {
		Linearize(input)
	}
}

func BenchmarkLinearize_NRGBA64(b *testing.B) {
	input := image.NewNRGBA64(image.Rect(0, 0, 512, 512))
	for i := 0; i < b.N; i++ {
		Linearize(input)
	}
}
