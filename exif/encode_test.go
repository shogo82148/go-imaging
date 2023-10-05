package exif

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEncode(t *testing.T) {
	tests := []string{
		"testdata/fireworks.exif",
		"testdata/flower.exif",
		"testdata/sado-gold-mine.exif",
		"testdata/senkakuwan.exif",
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			f, err := os.Open(test)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			tiff0, err := Decode(f)
			if err != nil {
				t.Fatal(err)
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
}
