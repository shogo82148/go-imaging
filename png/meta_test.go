package png

import (
	"bytes"
	"image"
	"math"
	"os"
	"testing"
)

func encodeDecodeWithMeta(m *ImageWithMeta) (*ImageWithMeta, error) {
	var b bytes.Buffer
	err := EncodeWithMeta(&b, m)
	if err != nil {
		return nil, err
	}
	return DecodeWithMeta(&b)
}

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

func TestEncodeWithMeta_Gamma(t *testing.T) {
	m := &ImageWithMeta{
		Image: image.NewNRGBA(image.Rect(0, 0, 100, 100)),
		Gamma: 0.45455,
	}

	decoded, err := encodeDecodeWithMeta(m)
	if err != nil {
		t.Fatal(err)
	}

	if math.Abs(decoded.Gamma-m.Gamma) > 1e-5 {
		t.Errorf("unexpected gamma: %f, want %f", decoded.Gamma, m.Gamma)
	}
}

func TestEncodeWithMeta_sRGB(t *testing.T) {
	m := &ImageWithMeta{
		Image: image.NewNRGBA(image.Rect(0, 0, 100, 100)),
		SRGB: &SRGB{
			RenderingIntent: RenderingIntentPerceptual,
		},
	}

	decoded, err := encodeDecodeWithMeta(m)
	if err != nil {
		t.Fatal(err)
	}

	if decoded.SRGB == nil {
		t.Fatal("unexpected nil sRGB")
	}

	if decoded.SRGB.RenderingIntent != RenderingIntentPerceptual {
		t.Errorf("unexpected rendering intent: %v, want %v", decoded.SRGB.RenderingIntent, m.SRGB.RenderingIntent)
	}
}
