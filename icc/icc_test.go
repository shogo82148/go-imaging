package icc

import (
	"log"
	"os"
	"testing"
)

func TestXxx(t *testing.T) {
	data, err := os.ReadFile("testdata/sRGB_IEC61966-2-1_black_scaled.icc")
	if err != nil {
		t.Fatal(err)
	}

	profile, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	_ = profile

	for _, t := range profile.Tags {
		log.Printf("%x, %#v", t.Tag, t.TagContent)
	}
}
