package exif

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func FuzzDecode(f *testing.F) {
	tests := []string{
		"testdata/fireworks.exif",
		"testdata/flower.exif",
		"testdata/sado-gold-mine.exif",
		"testdata/senkakuwan.exif",
	}

	for _, test := range tests {
		data, err := os.ReadFile(test)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(data)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		tiff0, err := Decode(bytes.NewReader(data))
		if err != nil {
			return
		}

		var buf bytes.Buffer
		err = Encode(&buf, tiff0)
		if err != nil {
			t.Fatal(err)
		}

		tiff1, err := Decode(&buf)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(tiff0, tiff1); diff != "" {
			t.Errorf("Encode() mismatch (-want +got):\n%s", diff)
		}
	})
}
