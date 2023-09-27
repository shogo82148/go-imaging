package jpeg

import (
	"os"
	"testing"
)

func TestDecodeWithMeta(t *testing.T) {
	img, err := decodeFileWithMeta("../testdata/senkakuwan.jpeg")
	if err != nil {
		t.Fatal(err)
	}
	_ = img // TODO: test me
}

func decodeFileWithMeta(filename string) (*ImageWithMeta, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return DecodeWithMeta(f)
}
