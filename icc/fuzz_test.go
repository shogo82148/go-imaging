package icc

import (
	"os"
	"testing"
)

func FuzzDecode(f *testing.F) {
	testdata := []string{
		"testdata/D65_XYZ.icc",
		"testdata/ILFORD_CANpro-4000_GPGFG_ProPlatin.icc",
		"testdata/iPhone12Pro.icc",
		"testdata/Probev1_ICCv2.icc",
		"testdata/Probev1_ICCv4.icc",
		"testdata/sRGB_ICC_v4_Appearance.icc",
		"testdata/sRGB_IEC61966-2-1_black_scaled.icc",
		"testdata/USWebCoatedSWOP.icc",
	}
	for _, path := range testdata {
		data, err := os.ReadFile(path)
		if err != nil {
			f.Error(err)
		}
		f.Add(data)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		_, err := Decode(data)
		if err != nil {
			return
		}
	})
}
