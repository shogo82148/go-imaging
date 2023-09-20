package resize

import (
	"testing"

	"github.com/shogo82148/go-imaging/resize/internal/golden"
)

func TestNearestNeighbor(t *testing.T) {
	src := golden.InputPattern()
	dst := golden.NewDst()
	NearestNeighbor(dst, src)
	golden.Assert(t, "nearest_neighbor", dst)
}

func BenchmarkNearestNeighbor(b *testing.B) {
	src := golden.InputPattern()
	dst := golden.NewDst()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NearestNeighbor(dst, src)
	}
}
