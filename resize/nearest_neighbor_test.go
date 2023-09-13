package resize

import (
	"image"
	"testing"

	"github.com/shogo82148/go-imaging/fp16"
)

func BenchmarkNearestNeighbor(b *testing.B) {
	dst := fp16.NewNRGBAh(image.Rect(0, 0, 512, 512))
	src := fp16.NewNRGBAh(image.Rect(0, 0, 256, 256))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NearestNeighbor(dst, src)
	}
}
