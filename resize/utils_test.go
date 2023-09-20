package resize

import (
	"math"
	"testing"
)

func TestScale(t *testing.T) {
	tests := []struct {
		x, srcDx, dstDx int
		srcX            int
		dx              float64
	}{
		{0, 3, 9, -1, 2 / 3.0},
		{1, 3, 9, 0, 0},
		{2, 3, 9, 0, 1 / 3.0},
		{3, 3, 9, 0, 2 / 3.0},
		{4, 3, 9, 1, 0},
		{5, 3, 9, 1, 1 / 3.0},
		{6, 3, 9, 1, 2 / 3.0},
		{7, 3, 9, 2, 0},
		{8, 3, 9, 2, 1 / 3.0},

		{0, 3, 6, -1, 0.75},
		{1, 3, 6, 0, 0.25},
		{2, 3, 6, 0, 0.75},
		{3, 3, 6, 1, 0.25},
		{4, 3, 6, 1, 0.75},
		{5, 3, 6, 2, 0.25},
	}

	for _, tt := range tests {
		srcX, dx := scale(tt.x, tt.srcDx, tt.dstDx)
		if srcX != tt.srcX {
			t.Errorf("x: %d, srcX: want %d, got %d", tt.x, tt.srcX, srcX)
		}
		if math.Abs(dx-tt.dx) > 1e-10 {
			t.Errorf("x: %d, dx: want %f, got %f", tt.x, tt.dx, dx)
		}
	}
}

func scaleOld(x, srcDx, dstDx int) (srcX int, dx float64) {
	fx := float64(x) + 0.5
	fx = (fx * float64(srcDx)) / float64(dstDx)
	fx -= 0.5
	dx = fx - math.Floor(fx)
	srcX = int(math.Floor(fx))
	return
}

func BenchmarkScale(b *testing.B) {
	b.Run("old", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			scaleOld(i%2000, 1000, 2000)
		}
	})
	b.Run("new", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			scale(i%2000, 1000, 2000)
		}
	})
}
