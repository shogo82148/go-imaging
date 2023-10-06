package jpeg

import (
	"bytes"
	"os"
	"testing"
)

func TestDecodeWithMeta(t *testing.T) {
	img, err := decodeFileWithMeta("../testdata/senkakuwan.jpeg")
	if err != nil {
		t.Fatal(err)
	}
	if img.ICCProfile == nil {
		t.Fatal("ICCProfile is nil")
	}
	if img.Exif == nil {
		t.Fatal("Exif is nil")
	}

	profileID := [16]uint8{0xec, 0xfd, 0xa3, 0x8e, 0x38, 0x85, 0x47, 0xc3, 0x6d, 0xb4, 0xbd, 0x4f, 0x7a, 0xda, 0x18, 0x2f}
	if img.ICCProfile.ProfileID != profileID {
		t.Fatalf("ProfileID is %#v, want %#v", img.ICCProfile.ProfileID, profileID)
	}
}

func decodeFileWithMeta(filename string) (*ImageWithMeta, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return DecodeWithMeta(f)
}

func TestEncodeWithMeta(t *testing.T) {
	img, err := decodeFileWithMeta("../testdata/senkakuwan.jpeg")
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)
	err = EncodeWithMeta(buf, img, nil)
	if err != nil {
		t.Fatal(err)
	}

	img2, err := DecodeWithMeta(buf)
	if err != nil {
		t.Fatal(err)
	}

	if img2.ICCProfile == nil {
		t.Fatal("ICCProfile is nil")
	}
	if img2.Exif == nil {
		t.Fatal("Exif is nil")
	}

	profileID := [16]uint8{0xec, 0xfd, 0xa3, 0x8e, 0x38, 0x85, 0x47, 0xc3, 0x6d, 0xb4, 0xbd, 0x4f, 0x7a, 0xda, 0x18, 0x2f}
	if img2.ICCProfile.ProfileID != profileID {
		t.Fatalf("ProfileID is %#v, want %#v", img.ICCProfile.ProfileID, profileID)
	}
}
