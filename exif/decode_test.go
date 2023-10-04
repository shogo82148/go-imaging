package exif

import (
	"os"
	"testing"
)

func TestDecode(t *testing.T) {
	f, err := os.Open("testdata/senkakuwan.exif")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	exif, err := Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", exif)
}
