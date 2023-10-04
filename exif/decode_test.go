package exif

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shogo82148/pointer"
)

func TestDecode(t *testing.T) {
	f, err := os.Open("testdata/senkakuwan.exif")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	got, err := Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	want := &Exif{
		Orientation: OrientationTopLeft,
		XResolution: &Rational{
			Numerator:   72,
			Denominator: 1,
		},
		YResolution: &Rational{
			Numerator:   72,
			Denominator: 1,
		},
		ResolutionUnit: ResolutionUnitInch,
		Make:           pointer.String("Apple"),
		Model:          pointer.String("iPhone 12 Pro"),
		Software:       pointer.String("16.6"),
		DateTime:       pointer.String("2023:07:13 12:41:13"),

		ExposureTime: &Rational{
			Numerator: 1, Denominator: 4464,
		},
		FNumber: &Rational{
			Numerator: 8, Denominator: 5,
		},
		ExposureProgram:   ExposureProgramNormal,
		ISOSpeedRatings:   []uint16{32},
		DateTimeOriginal:  pointer.String("2023:07:13 12:41:13"),
		DateTimeDigitized: pointer.String("2023:07:13 12:41:13"),
		ShutterSpeedValue: &SRational{
			Numerator: 98291, Denominator: 8107,
		},
		ApertureValue: &Rational{
			Numerator: 14447, Denominator: 10653,
		},
		BrightnessValue: &SRational{
			Numerator: 39410, Denominator: 3897,
		},
		ExposureBiasValue: &SRational{
			Numerator: 0, Denominator: 1,
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
