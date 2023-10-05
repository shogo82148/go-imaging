package exif

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEncode(t *testing.T) {
	f, err := os.Open("testdata/senkakuwan.exif")
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
	log.Printf("%x", buf.Bytes())
	os.WriteFile("hogehoge.exif", buf.Bytes(), 0644)

	tiff1, err := Decode(&buf)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(tiff0, tiff1); diff != "" {
		t.Errorf("Encode() mismatch (-want +got):\n%s", diff)
	}
}
