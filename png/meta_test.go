package png

import (
	"bytes"
	"fmt"
	"image"
	"math"
	"os"
	"testing"

	"github.com/shogo82148/go-imaging/icc"
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

func TestDecodeWithMeta_ICCP(t *testing.T) {
	f, err := os.Open("testdata/icc-profile.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	img, err := DecodeWithMeta(f)
	if err != nil {
		t.Fatal(err)
	}

	if img.ICCProfile == nil {
		t.Fatal("unexpected nil ICC Profile")
	}

	if got, want := img.ICCProfileName, "ICC Profile"; got != want {
		t.Errorf("unexpected ICC Profile Name: %q, want %q", got, want)
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

func TestICCProfileLatin1ToUTF8(t *testing.T) {
	tests := []struct {
		name string
		want string
		ok   bool
	}{
		{"", "", true},
		{"Valid Profile Name", "Valid Profile Name", true},
		{" Leading space is not permitted", "", false},
		{"Trailing space is not permitted ", "", false},
		{"Consecutive  spaces are not permitted", "", false},
		{"out of ascii: \xa1\xa2\xa3 \xfd\xfe\xff", "out of ascii: ¡¢£ ýþÿ", true},
		{"\x1f", "", false},
		{"\x7f", "", false},
		{"\xa0", "", false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q", tt.name), func(t *testing.T) {
			if got, ok := iccProfileLatin1ToUTF8(tt.name); got != tt.want || ok != tt.ok {
				t.Errorf("isValidICCProfileName() = %v %t, want %v %t", got, ok, tt.want, tt.ok)
			}
		})
	}
}

func TestICCProfileUTF8ToLatin1(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"", ""},
		{"Valid Profile Name", "Valid Profile Name"},
		{" Leading space is not permitted", "Leading space is not permitted"},
		{"Trailing space is not permitted ", "Trailing space is not permitted"},
		{"Consecutive  spaces are not permitted", "Consecutive spaces are not permitted"},
		{"out of ascii: ¡¢£ ýþÿ", "out of ascii: \xa1\xa2\xa3 \xfd\xfe\xff"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q", tt.name), func(t *testing.T) {
			if got := iccProfileUTF8ToLatin1(tt.name); got != tt.want {
				t.Errorf("iccProfileUTF8ToLatin1() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEncodeWithMeta_ICCProfile(t *testing.T) {
	t.Run("iPhone12Pro", func(t *testing.T) {
		f, err := os.Open("testdata/iPhone12Pro.icc")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		profile, err := icc.Decode(f)
		if err != nil {
			t.Fatal(err)
		}

		m := &ImageWithMeta{
			Image: image.NewNRGBA(image.Rect(0, 0, 100, 100)),
			// "International Color Consortium Profile" in Portuguese
			ICCProfileName: "Perfil do Consórcio Internacional de Cores",
			ICCProfile:     profile,
		}

		decoded, err := encodeDecodeWithMeta(m)
		if err != nil {
			t.Fatal(err)
		}

		if decoded.ICCProfile == nil {
			t.Fatal("unexpected nil ICC Profile")
		}
	})

	t.Run("gimp-linear", func(t *testing.T) {
		f, err := os.Open("testdata/gimp-linear.icc")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		profile, err := icc.Decode(f)
		if err != nil {
			t.Fatal(err)
		}

		m := &ImageWithMeta{
			Image: image.NewNRGBA(image.Rect(0, 0, 100, 100)),
			// "International Color Consortium Profile" in Portuguese
			ICCProfileName: "Perfil do Consórcio Internacional de Cores",
			ICCProfile:     profile,
		}

		decoded, err := encodeDecodeWithMeta(m)
		if err != nil {
			t.Fatal(err)
		}

		if decoded.ICCProfile == nil {
			t.Fatal("unexpected nil ICC Profile")
		}
	})
}
