package png

import (
	"math"
	"os"
	"testing"
)

func TestDecodeWithMeta_Gamma(t *testing.T) {
	// gamma.png has the gAMA chunk but no iCCP chunk.
	// chunk: IHDR
	// chunk: gAMA
	// chunk: bKGD
	// chunk: pHYs
	// chunk: tIME
	// chunk: tEXt
	// chunk: IDAT
	// chunk: IEND
	f, err := os.Open("testdata/gamma.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	img, err := DecodeWithMeta(f)
	if err != nil {
		t.Fatal(err)
	}

	if math.Abs(img.Gamma-0.45455) > 1e-5 {
		t.Errorf("unexpected gamma: %f, want 0.45455", img.Gamma)
	}
}
